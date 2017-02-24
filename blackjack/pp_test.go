package blackjack

import (
	"testing"

	"github.com/dwlnetnl/cards/card"
	"github.com/shopspring/decimal"
)

func TestMixedPerfectPairGame(t *testing.T) {
	testPlay(t, 2, HollandCasino, 10, 5, []event{
		perfectPair{Mixed, decimal.New(30, 0)},
		hand{
			dealer: Hand{card.Diamond(card.Six)},
			player: Hand{card.Heart(card.Two), card.Spade(card.Two)},
		},
		nextAction{[]Action{Hit, Stand, Split}, Stand},
		outcome{
			outcome: Won,
			amount:  decimal.New(20, 0),
			dealer: Hand{
				card.Diamond(card.Six),
				card.Spade(card.Queen),
				card.Diamond(card.Eight),
			},
			player: Hand{card.Heart(card.Two), card.Spade(card.Two)},
		},
	})
}

func TestSamePerfectPairGame(t *testing.T) {
	testPlay(t, 107, HollandCasino, 10, 5, []event{
		perfectPair{Same, decimal.New(60, 0)},
		hand{
			dealer: Hand{card.Heart(card.Ace)},
			player: Hand{card.Club(card.Ace), card.Spade(card.Ace)},
		},
		nextAction{[]Action{Hit, Stand, Split}, Hit},
		hand{
			dealer: Hand{card.Heart(card.Ace)},
			player: Hand{
				card.Club(card.Ace),
				card.Spade(card.Ace),
				card.Club(card.Seven),
			},
		},
		nextAction{[]Action{Hit, Stand}, Stand},
		outcome{
			outcome: Pushed,
			amount:  decimal.New(10, 0),
			dealer:  Hand{card.Heart(card.Ace), card.Spade(card.Eight)},
			player: Hand{
				card.Club(card.Ace),
				card.Spade(card.Ace),
				card.Club(card.Seven),
			},
		},
	})
}

func TestPerfectPerfectPairGame(t *testing.T) {
	testPlay(t, 38, HollandCasino, 10, 5, []event{
		perfectPair{Perfect, decimal.New(125, 0)},
		hand{
			dealer: Hand{card.Diamond(card.Nine)},
			player: Hand{card.Diamond(card.Ten), card.Diamond(card.Ten)},
		},
		nextAction{[]Action{Hit, Stand, Split}, Stand},
		outcome{
			outcome: Won,
			amount:  decimal.New(20, 0),
			dealer: Hand{
				card.Diamond(card.Nine),
				card.Spade(card.Four),
				card.Spade(card.Ace),
				card.Diamond(card.Nine),
			},
			player: Hand{card.Diamond(card.Ten), card.Diamond(card.Ten)},
		},
	})
}
