// Package blackjack implements a blackjack game engine.
package blackjack

import (
	"fmt"

	"github.com/dwlnetnl/cards/card"
	"github.com/dwlnetnl/cards/player"

	"github.com/shopspring/decimal"
)

// Action represents actions that can be taken by the player in a turn.
type Action int

// Available actions.
const (
	Hit Action = iota
	Stand
	Surrender
	Split
	Double
)

// Outcome represents the outcome of a bet.
type Outcome int

// Available outcomes.
const (
	Win Outcome = iota
	Loss
	Bust
	Push
	Surrendered
	Blackjack
	DealerBlackjack
)

// IO defines the game engine input and output interface.
type IO interface {
	Bet(fortune *player.Fortune) (amount decimal.Decimal)
	Hand(dealer, player Hand)
	NextAction(actions []Action) Action
	SplitHand(left, right Hand, amount decimal.Decimal)
	DoubleHand(hand Hand, withdrawn decimal.Decimal)
	Outcome(outcome Outcome, amount decimal.Decimal, dealer, player Hand)
	NewGame() bool

	PerfectPairBet(fortune *player.Fortune) (amount decimal.Decimal)
	PerfectPair(kind PerfectPair, amount decimal.Decimal)
}

type bet struct {
	hand    Hand
	amount  decimal.Decimal
	doubled bool
}

type game struct {
	rules    Rules
	fortune  *player.Fortune
	shuffler *card.Shuffler
	dealer   Hand
	bets     []*bet
}

var newShuffler = card.NewShuffler // for testing

// Play plays a blackjack game.
func Play(io IO, r Rules, f *player.Fortune) {
	s := newShuffler(card.NewStandardDeck(), r.NumDecks())
	g := &game{rules: r, fortune: f, shuffler: s}

	for {
		amount := io.Bet(g.fortune)
		if amount.Cmp(decimal.Zero) == 1 { // amount > 0
			g.fortune.Withdrawal(amount)
			g.setup(io, amount)

			if g.rules.PerfectPair() {
				g.perfectPair(io)
			}

			if g.dealer.IsBlackjack() {
				io.Outcome(DealerBlackjack, g.bets[0].amount, g.dealer, g.bets[0].hand)
			} else {
				g.play(io)
			}

			g.cleanup()
		}

		if !io.NewGame() {
			return
		}
	}
}

func (g *game) setup(io IO, amount decimal.Decimal) {
	g.dealer = Hand{g.shuffler.MustDraw()}
	if !g.rules.NoHoleCard() {
		g.dealer = append(g.dealer, g.shuffler.MustDraw())
	}

	g.bets = append(g.bets, &bet{
		amount: amount,
		hand:   Hand{g.shuffler.MustDraw(), g.shuffler.MustDraw()},
	})
}

func (g *game) cleanup() {
	g.shuffler.Shuffle(g.dealer...)
	for _, b := range g.bets {
		g.shuffler.Shuffle(b.hand...)
	}

	g.dealer = nil
	g.bets = nil
}

func (g *game) play(io IO) {
	// Cannot use a range expression, len(g.bets) may change during iteration.
	for i := 0; i < len(g.bets); i++ {
		b := g.bets[i]
		done := false

		for !done {
			io.Hand(g.dealer, b.hand)

			if b.hand.IsBlackjack() {
				g.blackjack(b, io)
				break
			}

			total, soft := b.hand.Points()
			if !soft && total >= 21 {
				break
			}

			switch g.nextAction(b, io) {
			case Hit:
				b.hand = append(b.hand, g.shuffler.MustDraw())
			case Stand:
				done = true
			case Surrender:
				panic("not implemented")
			case Split:
				g.fortune.Withdrawal(b.amount)
				lc := b.hand[0]
				rc := b.hand[1]
				lh := Hand{lc, g.shuffler.MustDraw()}
				rh := Hand{rc, g.shuffler.MustDraw()}
				b.hand = lh
				g.bets = append(g.bets, &bet{hand: rh, amount: b.amount})
				io.SplitHand(lh, rh, b.amount)
			case Double:
				amount := b.amount
				g.fortune.Withdrawal(amount)
				b.amount = b.amount.Add(amount)
				b.doubled = true
				b.hand = append(b.hand, g.shuffler.MustDraw())
				io.DoubleHand(b.hand, amount)
				done = true
			}
		}
	}

	if g.rules.NoHoleCard() {
		g.dealer = append(g.dealer, g.shuffler.MustDraw())
	}

	for !g.dealerFinished() {
		g.dealer = append(g.dealer, g.shuffler.MustDraw())
	}

	dealer, _ := g.dealer.Points()
	for _, b := range g.bets {
		if b.hand.IsBlackjack() {
			continue
		}

		player, _ := b.hand.Points()
		if player > 21 {
			g.bust(b, io)

		} else if player > dealer && dealer <= 21 || dealer > 21 {
			g.win(b, io)

		} else if player == dealer {
			if g.rules.DealerWinsTie() {
				g.loss(b, io)
			} else {
				g.push(b, io)
			}

		} else {
			g.loss(b, io)
		}
	}
}

func (g *game) nextAction(b *bet, io IO) Action {
	actions := g.availableActions(b)
	action := io.NextAction(actions)

	if !validAction(action, actions...) {
		panic(fmt.Sprintf("action %v is invalid, allowed: %v", action, actions))
	}

	return action
}

func (g *game) availableActions(b *bet) []Action {
	actions := []Action{Hit, Stand}

	if g.canSplit(b) {
		actions = append(actions, Split)
	}

	if !b.doubled && g.canDouble(b) {
		actions = append(actions, Double)
	}

	if g.rules.Surrender() != NoSurrender && !b.doubled && len(b.hand) == 2 {
		actions = append(actions, Surrender)
	}

	return actions
}

func validAction(action Action, actions ...Action) bool {
	for _, a := range actions {
		if a == action {
			return true
		}
	}

	return false
}

func (g *game) dealerFinished() bool {
	total, soft := g.dealer.Points()

	if total == 17 && soft {
		return !g.rules.DealerHitSoft17()
	}

	return total >= 17
}

func (g *game) blackjack(b *bet, io IO) {
	ratio := g.rules.BlackjackRatio()
	num := decimal.New(ratio.Num().Int64(), 0)
	denom := decimal.New(ratio.Denom().Int64(), 0)

	// e.g. ratio is 3/2 and wager is 10
	// amount = 10 * 3 / 2 + 10
	//        = 30 / 2 + 10
	//        = 15 + 10
	//        = 25
	amount := b.amount.Mul(num).Div(denom).Add(b.amount)

	g.fortune.Deposit(amount)
	io.Outcome(Blackjack, amount, g.dealer, b.hand)
}

func (g *game) win(b *bet, io IO) {
	amount := b.amount.Mul(decimal.New(2, 0))
	g.fortune.Deposit(amount)
	io.Outcome(Win, amount, g.dealer, b.hand)
}

func (g *game) push(b *bet, io IO) {
	g.fortune.Deposit(b.amount)
	io.Outcome(Push, b.amount, g.dealer, b.hand)
}

func (g *game) loss(b *bet, io IO) {
	io.Outcome(Loss, b.amount, g.dealer, b.hand)
}

func (g *game) bust(b *bet, io IO) {
	io.Outcome(Bust, b.amount, g.dealer, b.hand)
}

func (g *game) canSplit(b *bet) bool {
	if len(b.hand) != 2 {
		return false
	}

	if b.hand[0].Rank != b.hand[1].Rank {
		return false
	}

	if !g.fortune.Has(b.amount) {
		return false
	}

	hands := make([]Hand, len(g.bets))
	for i, b := range g.bets {
		hands[i] = b.hand
	}

	return g.rules.CanSplit(hands)
}

func (g *game) canDouble(b *bet) bool {
	if len(b.hand) > 2 {
		return false
	}

	if !g.fortune.Has(b.amount) {
		return false
	}

	if dr := g.rules.Double(); dr != DoubleAny {
		pts, _ := b.hand.Points()
		if !(pts == 11 || pts == 10 || pts == 9 && dr == DoubleOnly9_10_11) {
			return false
		}
	}

	if len(g.bets) > 1 {
		return g.rules.DoubleAfterSplit()
	}

	return true
}
