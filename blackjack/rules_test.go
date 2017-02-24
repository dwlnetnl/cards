package blackjack

import (
	"math/big"
	"testing"

	"github.com/dwlnetnl/cards/card"

	"github.com/shopspring/decimal"
)

type testRules struct {
	surrender SurrenderRule
}

func (r testRules) NumDecks() uint                  { return 6 }
func (r testRules) DealerHitSoft17() bool           { return true }
func (r testRules) Surrender() SurrenderRule        { return r.surrender }
func (r testRules) CanSplit([]Hand) bool            { return true }
func (r testRules) Double() DoubleRule              { return DoubleAny }
func (r testRules) DoubleAfterSplit() bool          { return true }
func (r testRules) NoHoleCard() bool                { return true }
func (r testRules) OriginalBetsOnly() bool          { return false }
func (r testRules) BlackjackRatio() *big.Rat        { return big.NewRat(3, 2) }
func (r testRules) DealerWinsTie() bool             { return true }
func (r testRules) PerfectPair() bool               { return false }
func (r testRules) PerfectPairRatio() (m, s, p int) { return }

func TestDoubleAfterSplit(t *testing.T) {
	rules := testRules{surrender: EarlySurrender}
	testPlay(t, 38, rules, 10, 5, []event{
		hand{
			dealer: Hand{card.Diamond(card.Nine)},
			player: Hand{card.Diamond(card.Ten), card.Diamond(card.Ten)},
		},
		nextAction{[]Action{Surrender, Continue}, Continue},
		hand{
			dealer: Hand{card.Diamond(card.Nine)},
			player: Hand{card.Diamond(card.Ten), card.Diamond(card.Ten)},
		},
		nextAction{[]Action{Hit, Stand, Split, Double}, Split},
		splitHand{
			left:  Hand{card.Diamond(card.Ten), card.Spade(card.Four)},
			right: Hand{card.Diamond(card.Ten), card.Spade(card.Ace)},
		},
		hand{
			dealer: Hand{card.Diamond(card.Nine)},
			player: Hand{card.Diamond(card.Ten), card.Spade(card.Four)},
		},
		nextAction{[]Action{Hit, Stand, Double}, Double},
		doubleHand{
			hand: Hand{
				card.Diamond(card.Ten),
				card.Spade(card.Four),
				card.Diamond(card.Nine),
			},
			withdrawn: decimal.New(10, 0),
		},
		hand{
			dealer: Hand{card.Diamond(card.Nine)},
			player: Hand{card.Diamond(card.Ten), card.Spade(card.Ace)},
		},
		nextAction{[]Action{Hit, Stand, Double}, Stand},
		outcome{
			outcome: Bust,
			amount:  decimal.New(20, 0),
			dealer:  Hand{card.Diamond(card.Nine), card.Diamond(card.Ten)},
			player: Hand{
				card.Diamond(card.Ten),
				card.Spade(card.Four),
				card.Diamond(card.Nine),
			},
		},
		outcome{
			outcome: Won,
			amount:  decimal.New(20, 0),
			dealer:  Hand{card.Diamond(card.Nine), card.Diamond(card.Ten)},
			player:  Hand{card.Diamond(card.Ten), card.Spade(card.Ace)},
		},
	})
}

func TestDealerWinsTie(t *testing.T) {
	rules := testRules{surrender: NoSurrender}
	testPlay(t, 10, rules, 10, 0, []event{
		hand{
			dealer: Hand{card.Heart(card.King)},
			player: Hand{card.Heart(card.Three), card.Club(card.Queen)},
		},
		nextAction{[]Action{Hit, Stand, Double}, Hit},
		hand{
			dealer: Hand{card.Heart(card.King)},
			player: Hand{
				card.Heart(card.Three),
				card.Club(card.Queen),
				card.Heart(card.Eight),
			},
		},
		nextAction{[]Action{Hit, Stand}, Stand},
		outcome{
			outcome: Lost,
			amount:  decimal.New(10, 0),
			dealer:  Hand{card.Heart(card.King), card.Spade(card.Ace)},
			player: Hand{
				card.Heart(card.Three),
				card.Club(card.Queen),
				card.Heart(card.Eight),
			},
		},
	})
}

// func TestDealerHitSoft17(t *testing.T) {
// 	rules := testRules{surrender: NoSurrender}
// 	testPlay(t, seed, rules, 10, 0, []event{
//
// 	})
// }

func TestEarlySurrendered(t *testing.T) {
	rules := testRules{surrender: EarlySurrender}
	testPlay(t, 5, rules, 10, 0, []event{
		hand{
			dealer: Hand{card.Spade(card.King)},
			player: Hand{card.Diamond(card.Five), card.Diamond(card.Jack)},
		},
		nextAction{[]Action{Surrender, Continue}, Surrender},
		outcome{
			outcome: Surrendered,
			amount:  decimal.New(5, 0),
			dealer:  Hand{card.Spade(card.King)},
			player:  Hand{card.Diamond(card.Five), card.Diamond(card.Jack)},
		},
	})
}

func TestLateSurrendered(t *testing.T) {
	rules := testRules{surrender: LateSurrender}
	testPlay(t, 164, rules, 10, 0, []event{
		hand{
			dealer: Hand{card.Diamond(card.Ten)},
			player: Hand{card.Heart(card.Six), card.Club(card.Jack)},
		},
		nextAction{[]Action{Surrender, Continue}, Surrender},
		outcome{
			outcome: Surrendered,
			amount:  decimal.New(5, 0),
			dealer:  Hand{card.Diamond(card.Ten)},
			player:  Hand{card.Heart(card.Six), card.Club(card.Jack)},
		},
	})
}
