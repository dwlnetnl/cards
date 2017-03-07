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
	Split
	Double
	Surrender
	Continue
)

//go:generate stringer -type=Action

// Outcome represents the outcome of a bet.
type Outcome int

// Available outcomes.
const (
	Won Outcome = iota
	Lost
	Bust
	Pushed
	Surrendered
	Blackjack
	DealerBlackjack
)

//go:generate stringer -type=Outcome

// UI defines the game engine interface for user interaction.
type UI interface {
	Bet(fortune *player.Fortune) (amount decimal.Decimal)
	Hand(dealer, player Hand)
	DealerCard(card card.Card, hand Hand)
	NextAction(actions []Action) Action
	SplitHand(left, right Hand, amount decimal.Decimal)
	DoubleHand(hand Hand, withdrawn decimal.Decimal)
	Outcome(outcome Outcome, amount decimal.Decimal, dealer, player Hand)
	NewGame(fortune *player.Fortune) bool
	NoActiveFortune(fortune *player.Fortune) bool
	NoFortune()

	PerfectPairBet(fortune *player.Fortune) (amount decimal.Decimal)
	PerfectPair(kind PerfectPair, amount decimal.Decimal)
}

type bet struct {
	hand    Hand
	amount  decimal.Decimal
	doubled bool
}

type game struct {
	ui       UI
	rules    Rules
	fortune  *player.Fortune
	shuffler *card.Shuffler
	dealer   Hand
	bets     []*bet
}

var newShuffler = card.NewShuffler // for testing

// Play starts a blackjack game.
func Play(ui UI, r Rules, f *player.Fortune) {
	s := newShuffler(card.NewStandardDeck(), r.NumDecks())
	g := &game{ui: ui, rules: r, fortune: f, shuffler: s}

	for {
		amount := g.bet()
		if amount.Equal(decimal.Zero) {
			return
		}

		g.setup(amount)

		if g.rules.PerfectPair() {
			g.perfectPair()
		}

		b := g.bets[0]
		if g.dealer.IsBlackjack() {
			ui.Outcome(DealerBlackjack, b.amount, g.dealer, b.hand)
		} else {
			if b.hand.IsBlackjack() {
				g.blackjack()
			} else {
				g.play()
			}
		}

		g.cleanup()

		if !ui.NewGame(f) {
			return
		}
	}
}

func (g *game) bet() decimal.Decimal {
	if g.fortune.Total().Equal(decimal.Zero) {
		g.ui.NoFortune()
		return decimal.Zero
	}

	amount := g.ui.Bet(g.fortune)
	if amount.Cmp(decimal.Zero) <= 0 { // amount <= 0
		if !g.ui.NewGame(g.fortune) {
			return decimal.Zero
		}
		return g.bet()
	}

	if !g.fortune.Has(amount) {
		return g.bet()
	}

	g.fortune.Withdrawal(amount)
	return amount
}

func (g *game) setup(amount decimal.Decimal) {
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

func (g *game) play() {
	if g.canEarlySurrender() {
		b := g.bets[0]
		g.ui.Hand(g.dealer, b.hand)
		switch a := g.nextAction(b, []Action{Surrender, Continue}); a {
		case Continue:
		case Surrender:
			g.surrender(b)
			return
		default:
			panic(fmt.Sprintf("unexpected action: %v", a))
		}
	}

	// Can't use a range expression because it evaluates the length
	// of g.bets only once, g.bets may change during iteration.
	for i := 0; i < len(g.bets); i++ {
		b := g.bets[i]
		done := false

		for !done {
			g.ui.Hand(g.dealer, b.hand)

			total, soft := b.hand.Points()
			if !soft && total >= 21 {
				break
			}

			var a Action
			if total == 21 {
				a = Stand
			} else {
				a = g.nextAction(b, g.availableActions(b))
			}

			switch a {
			case Hit:
				b.hand = append(b.hand, g.shuffler.MustDraw())
			case Stand:
				done = true
			case Split:
				g.fortune.Withdrawal(b.amount)
				lc := b.hand[0]
				rc := b.hand[1]
				lh := Hand{lc, g.shuffler.MustDraw()}
				rh := Hand{rc, g.shuffler.MustDraw()}
				b.hand = lh
				g.bets = append(g.bets, &bet{hand: rh, amount: b.amount})
				g.ui.SplitHand(lh, rh, b.amount)
			case Double:
				amount := b.amount
				g.fortune.Withdrawal(amount)
				b.amount = b.amount.Add(amount)
				b.doubled = true
				b.hand = append(b.hand, g.shuffler.MustDraw())
				g.ui.DoubleHand(b.hand, amount)
				done = true
			case Surrender: // late surrender
				g.surrender(b)
				return
			case Continue:
			default:
				panic(fmt.Sprintf("unexpected action: %v", a))
			}
		}
	}

	if g.rules.NoHoleCard() && len(g.dealer) == 1 {
		c := g.shuffler.MustDraw()
		g.dealer = append(g.dealer, c)
		g.ui.DealerCard(c, g.dealer)
	}

	for !g.dealerFinished() {
		c := g.shuffler.MustDraw()
		g.dealer = append(g.dealer, c)
		g.ui.DealerCard(c, g.dealer)
	}

	dealer, _ := g.dealer.Points()
	for _, b := range g.bets {
		player, _ := b.hand.Points()
		if player > 21 {
			g.bust(b)

		} else if player > dealer && dealer <= 21 || dealer > 21 {
			g.win(b)

		} else if player == dealer {
			if g.rules.DealerWinsTie() {
				g.loss(b)
			} else {
				g.push(b)
			}

		} else {
			g.loss(b)
		}
	}
}

func (g *game) availableActions(b *bet) []Action {
	if g.canLateSurrender(b) {
		return []Action{Surrender, Continue}
	}

	actions := []Action{Hit, Stand}
	if g.canSplit(b) {
		actions = append(actions, Split)
	}

	if !b.doubled && g.canDouble(b) {
		actions = append(actions, Double)
	}

	return actions
}

func (g *game) nextAction(b *bet, actions []Action) Action {
	action := g.ui.NextAction(actions)
	if !validAction(action, actions...) {
		panic(fmt.Sprintf("action %v is invalid, allowed: %v", action, actions))
	}
	return action
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

func (g *game) blackjack() {
	b := g.bets[0]
	amount := b.amount.Mul(g.rules.BlackjackRatio()).Add(b.amount)
	g.fortune.Deposit(amount)
	g.ui.Outcome(Blackjack, amount, g.dealer, b.hand)
}

func (g *game) win(b *bet) {
	amount := b.amount.Mul(decimal.New(2, 0))
	g.fortune.Deposit(amount)
	g.ui.Outcome(Won, amount, g.dealer, b.hand)
}

func (g *game) push(b *bet) {
	g.fortune.Deposit(b.amount)
	g.ui.Outcome(Pushed, b.amount, g.dealer, b.hand)
}

func (g *game) loss(b *bet) {
	amount := decimal.Zero.Sub(b.amount)
	g.ui.Outcome(Lost, amount, g.dealer, b.hand)
}

func (g *game) bust(b *bet) {
	amount := decimal.Zero.Sub(b.amount)
	g.ui.Outcome(Bust, amount, g.dealer, b.hand)
}

func (g *game) surrender(b *bet) {
	amount := b.amount.Div(decimal.New(2, 0))
	g.fortune.Deposit(amount)
	g.ui.Outcome(Surrendered, amount, g.dealer, b.hand)
}

func (g *game) canSplit(b *bet) bool {
	if len(b.hand) != 2 || b.hand[0].Rank != b.hand[1].Rank {
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
	if len(b.hand) > 2 || !g.fortune.Has(b.amount) {
		return false
	}

	if dr := g.rules.Double(); dr != DoubleAny {
		pts, _ := b.hand.Points()

		if dr == DoubleOnly9_10_11 && (pts < 9 || pts > 11) {
			return false
		}

		if dr == DoubleOnly10_11 && (pts < 10 || pts > 11) {
			return false
		}
	}

	if len(g.bets) > 1 {
		return g.rules.DoubleAfterSplit()
	}

	return true
}

func (g *game) canEarlySurrender() bool {
	if len(g.bets) != 1 {
		panic("only first bet can be surrendered early")
	}
	return g.rules.Surrender() == EarlySurrender && len(g.dealer) == 1
}

func (g *game) canLateSurrender(b *bet) bool {
	return g.rules.Surrender() == LateSurrender && len(g.dealer) == 1 &&
		len(g.bets) == 1 && len(g.bets[0].hand) == 2 && !b.doubled
}
