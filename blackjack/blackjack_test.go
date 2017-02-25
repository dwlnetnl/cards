package blackjack

import (
	"testing"

	"github.com/dwlnetnl/cards/card"

	"github.com/shopspring/decimal"
)

func TestUnbettedGame(t *testing.T) {
	testPlay(t, 1, HollandCasino, 0, 0, []event{})
}

func TestBlackjackGame(t *testing.T) {
	testPlay(t, 1, HollandCasino, 10, 0, []event{
		outcome{
			outcome: Blackjack,
			amount:  decimal.New(25, 0),
			dealer:  Hand{card.Heart(card.Jack)},
			player:  Hand{card.Heart(card.Ace), card.Heart(card.Queen)},
		},
	})
}

func TestLostGame(t *testing.T) {
	testPlay(t, 2, HollandCasino, 10, 0, []event{
		hand{
			dealer: Hand{card.Diamond(card.Six)},
			player: Hand{card.Heart(card.Two), card.Spade(card.Two)},
		},
		nextAction{[]Action{Hit, Stand, Split}, Hit},
		hand{
			dealer: Hand{card.Diamond(card.Six)},
			player: Hand{
				card.Heart(card.Two),
				card.Spade(card.Two),
				card.Spade(card.Queen),
			},
		},
		nextAction{[]Action{Hit, Stand}, Stand},
		dealerCard{
			card: card.Diamond(card.Eight),
			hand: Hand{card.Diamond(card.Six), card.Diamond(card.Eight)},
		},
		dealerCard{
			card: card.Club(card.Five),
			hand: Hand{
				card.Diamond(card.Six),
				card.Diamond(card.Eight),
				card.Club(card.Five),
			},
		},
		outcome{
			outcome: Lost,
			amount:  decimal.New(10, 0),
			dealer: Hand{
				card.Diamond(card.Six),
				card.Diamond(card.Eight),
				card.Club(card.Five),
			},
			player: Hand{
				card.Heart(card.Two),
				card.Spade(card.Two),
				card.Spade(card.Queen),
			},
		},
	})
}

func TestBustGame(t *testing.T) {
	testPlay(t, 2, HollandCasino, 10, 0, []event{
		hand{
			dealer: Hand{card.Diamond(card.Six)},
			player: Hand{card.Heart(card.Two), card.Spade(card.Two)},
		},
		nextAction{[]Action{Hit, Stand, Split}, Hit},
		hand{
			dealer: Hand{card.Diamond(card.Six)},
			player: Hand{
				card.Heart(card.Two),
				card.Spade(card.Two),
				card.Spade(card.Queen),
			},
		},
		nextAction{[]Action{Hit, Stand}, Hit},
		hand{
			dealer: Hand{card.Diamond(card.Six)},
			player: Hand{
				card.Heart(card.Two),
				card.Spade(card.Two),
				card.Spade(card.Queen),
				card.Diamond(card.Eight),
			},
		},
		dealerCard{
			card: card.Club(card.Five),
			hand: Hand{card.Diamond(card.Six), card.Club(card.Five)},
		},
		dealerCard{
			card: card.Heart(card.Two),
			hand: Hand{
				card.Diamond(card.Six),
				card.Club(card.Five),
				card.Heart(card.Two),
			},
		},
		dealerCard{
			card: card.Spade(card.Seven),
			hand: Hand{
				card.Diamond(card.Six),
				card.Club(card.Five),
				card.Heart(card.Two),
				card.Spade(card.Seven),
			},
		},
		outcome{
			outcome: Bust,
			amount:  decimal.New(10, 0),
			dealer: Hand{
				card.Diamond(card.Six),
				card.Club(card.Five),
				card.Heart(card.Two),
				card.Spade(card.Seven),
			},
			player: Hand{
				card.Heart(card.Two),
				card.Spade(card.Two),
				card.Spade(card.Queen),
				card.Diamond(card.Eight),
			},
		},
	})
}

func TestDoubleOnFirstHandOnly(t *testing.T) {
	testPlay(t, 33, HollandCasino, 10, 0, []event{
		hand{
			dealer: Hand{card.Spade(card.Ten)},
			player: Hand{card.Club(card.Six), card.Club(card.Two)},
		},
		nextAction{[]Action{Hit, Stand}, Hit},
		hand{
			dealer: Hand{card.Spade(card.Ten)},
			player: Hand{
				card.Club(card.Six),
				card.Club(card.Two),
				card.Heart(card.Two),
			},
		},
		nextAction{[]Action{Hit, Stand}, Stand},
		dealerCard{
			card: card.Spade(card.Nine),
			hand: Hand{card.Spade(card.Ten), card.Spade(card.Nine)},
		},
		outcome{
			outcome: Lost,
			amount:  decimal.New(10, 0),
			dealer:  Hand{card.Spade(card.Ten), card.Spade(card.Nine)},
			player: Hand{
				card.Club(card.Six),
				card.Club(card.Two),
				card.Heart(card.Two),
			},
		},
	})
}

func TestPushedGame(t *testing.T) {
	testPlay(t, 10, HollandCasino, 10, 0, []event{
		hand{
			dealer: Hand{card.Heart(card.King)},
			player: Hand{card.Heart(card.Three), card.Club(card.Queen)},
		},
		nextAction{[]Action{Hit, Stand}, Hit},
		hand{
			dealer: Hand{card.Heart(card.King)},
			player: Hand{
				card.Heart(card.Three),
				card.Club(card.Queen),
				card.Heart(card.Eight),
			},
		},
		nextAction{[]Action{Hit, Stand}, Stand},
		dealerCard{
			card: card.Spade(card.Ace),
			hand: Hand{card.Heart(card.King), card.Spade(card.Ace)},
		},
		outcome{
			outcome: Pushed,
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

func TestSplittedGame(t *testing.T) {
	testPlay(t, 14, HollandCasino, 10, 0, []event{
		hand{
			dealer: Hand{card.Spade(card.Six)},
			player: Hand{card.Heart(card.Eight), card.Spade(card.Eight)},
		},
		nextAction{[]Action{Hit, Stand, Split}, Split},
		splitHand{
			left:  Hand{card.Heart(card.Eight), card.Diamond(card.Ace)},
			right: Hand{card.Spade(card.Eight), card.Club(card.Eight)},
		},
		hand{
			dealer: Hand{card.Spade(card.Six)},
			player: Hand{card.Heart(card.Eight), card.Diamond(card.Ace)},
		},
		nextAction{[]Action{Hit, Stand}, Stand},
		hand{
			dealer: Hand{card.Spade(card.Six)},
			player: Hand{card.Spade(card.Eight), card.Club(card.Eight)},
		},
		nextAction{[]Action{Hit, Stand, Split}, Split},
		splitHand{
			left:  Hand{card.Spade(card.Eight), card.Club(card.King)},
			right: Hand{card.Club(card.Eight), card.Diamond(card.Eight)},
		},
		hand{
			dealer: Hand{card.Spade(card.Six)},
			player: Hand{card.Spade(card.Eight), card.Club(card.King)},
		},
		nextAction{[]Action{Hit, Stand}, Stand},
		hand{
			dealer: Hand{card.Spade(card.Six)},
			player: Hand{card.Club(card.Eight), card.Diamond(card.Eight)},
		},
		nextAction{[]Action{Hit, Stand, Split}, Split},
		splitHand{
			left:  Hand{card.Club(card.Eight), card.Club(card.Seven)},
			right: Hand{card.Diamond(card.Eight), card.Heart(card.Eight)},
		},
		hand{
			dealer: Hand{card.Spade(card.Six)},
			player: Hand{card.Club(card.Eight), card.Club(card.Seven)},
		},
		nextAction{[]Action{Hit, Stand}, Stand},
		hand{
			dealer: Hand{card.Spade(card.Six)},
			player: Hand{card.Diamond(card.Eight), card.Heart(card.Eight)},
		},
		nextAction{[]Action{Hit, Stand}, Stand},
		dealerCard{
			card: card.Club(card.Ten),
			hand: Hand{card.Spade(card.Six), card.Club(card.Ten)},
		},
		dealerCard{
			card: card.Heart(card.Three),
			hand: Hand{
				card.Spade(card.Six),
				card.Club(card.Ten),
				card.Heart(card.Three),
			},
		},
		outcome{
			outcome: Pushed,
			amount:  decimal.New(10, 0),
			dealer: Hand{
				card.Spade(card.Six),
				card.Club(card.Ten),
				card.Heart(card.Three),
			},
			player: Hand{card.Heart(card.Eight), card.Diamond(card.Ace)},
		},
		outcome{
			outcome: Lost,
			amount:  decimal.New(10, 0),
			dealer: Hand{
				card.Spade(card.Six),
				card.Club(card.Ten),
				card.Heart(card.Three),
			},
			player: Hand{card.Spade(card.Eight), card.Club(card.King)},
		},
		outcome{
			outcome: Lost,
			amount:  decimal.New(10, 0),
			dealer: Hand{
				card.Spade(card.Six),
				card.Club(card.Ten),
				card.Heart(card.Three),
			},
			player: Hand{card.Club(card.Eight), card.Club(card.Seven)},
		},
		outcome{
			outcome: Lost,
			amount:  decimal.New(10, 0),
			dealer: Hand{
				card.Spade(card.Six),
				card.Club(card.Ten),
				card.Heart(card.Three),
			},
			player: Hand{card.Diamond(card.Eight), card.Heart(card.Eight)},
		},
	})
}

func TestDoubledGame(t *testing.T) {
	testPlay(t, 27, HollandCasino, 10, 0, []event{
		hand{
			dealer: Hand{card.Club(card.Ten)},
			player: Hand{card.Club(card.Two), card.Diamond(card.Seven)},
		},
		nextAction{[]Action{Hit, Stand, Double}, Double},
		doubleHand{
			hand: Hand{
				card.Club(card.Two),
				card.Diamond(card.Seven),
				card.Club(card.Queen),
			},
			withdrawn: decimal.New(10, 0),
		},
		dealerCard{
			card: card.Club(card.Five),
			hand: Hand{card.Club(card.Ten), card.Club(card.Five)},
		},
		dealerCard{
			card: card.Spade(card.Nine),
			hand: Hand{
				card.Club(card.Ten),
				card.Club(card.Five),
				card.Spade(card.Nine),
			},
		},
		outcome{
			outcome: Won,
			amount:  decimal.New(40, 0),
			dealer: Hand{
				card.Club(card.Ten),
				card.Club(card.Five),
				card.Spade(card.Nine),
			},
			player: Hand{
				card.Club(card.Two),
				card.Diamond(card.Seven),
				card.Club(card.Queen),
			},
		},
	})
}

func TestDealerLostGame(t *testing.T) {
	testPlay(t, 16, HollandCasino, 10, 0, []event{
		hand{
			dealer: Hand{card.Spade(card.Two)},
			player: Hand{card.Diamond(card.Two), card.Heart(card.King)},
		},
		nextAction{[]Action{Hit, Stand}, Stand},
		dealerCard{
			card: card.Diamond(card.Eight),
			hand: Hand{card.Spade(card.Two), card.Diamond(card.Eight)},
		},
		dealerCard{
			card: card.Spade(card.Queen),
			hand: Hand{
				card.Spade(card.Two),
				card.Diamond(card.Eight),
				card.Spade(card.Queen),
			},
		},
		outcome{
			outcome: Lost,
			amount:  decimal.New(10, 0),
			dealer: Hand{
				card.Spade(card.Two),
				card.Diamond(card.Eight),
				card.Spade(card.Queen),
			},
			player: Hand{card.Diamond(card.Two), card.Heart(card.King)},
		},
	})
}

func TestDealerBlackjackGame(t *testing.T) {
	testPlay(t, 15, TapTapBoom, 10, 0, []event{
		outcome{
			outcome: DealerBlackjack,
			amount:  decimal.New(10, 0),
			dealer:  Hand{card.Club(card.Ace), card.Club(card.Jack)},
			player:  Hand{card.Heart(card.Eight), card.Spade(card.Seven)},
		},
	})
}
