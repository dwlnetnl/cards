package blackjack

import (
	"testing"

	"github.com/dwlnetnl/cards/card"
)

func makeHand(ranks ...card.Rank) Hand {
	h := make(Hand, len(ranks))
	for i, r := range ranks {
		h[i] = card.Card{Suit: card.Naked, Rank: r}
	}

	return h
}

func TestHandIsBlackjack(t *testing.T) {
	cases := []struct {
		in   Hand
		want bool
	}{
		{makeHand(card.Ace, card.Ten), true},
		{makeHand(card.Ten, card.Ace), true},
		{makeHand(card.Ace, card.Six, card.Four), false},
	}

	for _, c := range cases {
		t.Run(c.in.Names(","), func(t *testing.T) {
			got := c.in.IsBlackjack()
			if got != c.want {
				t.Errorf("got: %v, want: %v", got, c.want)
			}
		})
	}
}

func TestHandPoints(t *testing.T) {
	cases := []struct {
		in    Hand
		total int
		soft  bool
	}{
		{makeHand(card.Ace, card.Ten), 21, true},
		{makeHand(card.Ace, card.Six, card.Four), 21, true},
		{makeHand(card.Ace, card.Ace, card.Ten, card.Nine), 21, false},
		{makeHand(card.Ace, card.Ace, card.Ace), 13, false},
		{makeHand(card.Ace, card.Six), 17, true},
		{makeHand(card.Ace, card.Six, card.Five), 12, false},
		{makeHand(card.Ace, card.Seven, card.Nine), 17, false},
	}

	for _, c := range cases {
		t.Run(c.in.Names(","), func(t *testing.T) {
			total, soft := c.in.Points()

			if total != c.total {
				t.Errorf("[%s] = %d, want: %d", c.in.Names(", "), total, c.total)
			}

			if soft != c.soft {
				t.Errorf("soft = %v, want: %v", soft, c.soft)
			}
		})
	}
}

func TestHandPerfectPair(t *testing.T) {
	cases := []struct {
		in   Hand
		want PerfectPair
	}{
		{Hand{card.Spade(card.Queen), card.Diamond(card.Queen)}, Mixed},
		{Hand{card.Spade(card.Queen), card.Club(card.Queen)}, Same},
		{Hand{card.Spade(card.Queen), card.Spade(card.Queen)}, Perfect},
	}

	for _, c := range cases {
		t.Run(c.in.Names(","), func(t *testing.T) {
			got := c.in.perfectPair()
			if got != c.want {
				t.Errorf("got: %v, want: %v", got, c.want)
			}
		})
	}

	t.Run("Invalid", func(t *testing.T) {
		cases := []Hand{
			Hand{card.Spade(card.Queen)},
			Hand{card.Spade(card.Queen), card.Spade(card.King)},
			Hand{card.Spade(card.Queen), card.Spade(card.King), card.Heart(card.Ace)},
		}

		for _, h := range cases {
			t.Run(h.Names(","), func(t *testing.T) {
				got := h.perfectPair()
				const want = NoPair

				if got != want {
					t.Errorf("got: %v, want: %v", got, want)
				}
			})
		}
	})
}
