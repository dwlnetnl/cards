package blackjack

import "github.com/shopspring/decimal"

// PerfectPair represents an outcome of a perfect pair side bet.
type PerfectPair int

// Outcomes of a perfect pair side bet.
const (
	NoPair PerfectPair = iota
	Mixed
	Same
	Perfect
)

//go:generate stringer -type=PerfectPair

func (h Hand) perfectPair() PerfectPair {
	if len(h) != 2 || h[0].Rank != h[1].Rank {
		return NoPair
	}
	if h[0].Suit == h[1].Suit {
		return Perfect
	}
	if h[0].Color() == h[1].Color() {
		return Same
	}
	return Mixed
}

func (g *game) perfectPair() {
	amount := g.ui.PerfectPairBet(g.fortune)
	if amount.Equal(decimal.Zero) {
		return
	}

	g.fortune.Withdrawal(amount)

	b := g.bets[0]
	pp := b.hand.perfectPair()
	if pp != NoPair {
		var factor int
		mixed, same, perfect := g.rules.PerfectPairRatio()

		switch pp {
		case Mixed:
			factor = mixed
		case Same:
			factor = same
		case Perfect:
			factor = perfect
		}

		if factor > 0 {
			factor := decimal.New(int64(factor), 0)
			amount = amount.Mul(factor)
			g.fortune.Deposit(amount)
			g.ui.PerfectPair(pp, amount)
		}
	}
}
