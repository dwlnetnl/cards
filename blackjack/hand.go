package blackjack

import (
	"fmt"
	"strings"

	"github.com/dwlnetnl/cards/card"
)

// Hand represents a players hand in blackjack.
type Hand []card.Card

// Names returns a string of all card names joined together separated by sep.
func (h Hand) Names(sep string) string {
	names := make([]string, len(h))
	for i, c := range h {
		names[i] = c.String()
	}

	return strings.Join(names, sep)
}

func (h Hand) String() string {
	if h == nil {
		return "<nil>"
	}
	points, _ := h.Points()
	return fmt.Sprintf("%s (%d)", h.Names(", "), points)
}

var rankPoints = [...]int{
	2, 3, 4, 5, 6, 7, 8,
	9, 10, 10, 10, 10, 11,
}

// IsBlackjack returns true if hand h is a blackjack.
func (h Hand) IsBlackjack() bool {
	if len(h) != 2 {
		return false
	}
	return rankPoints[h[0].Rank]+rankPoints[h[1].Rank] == 21
}

// Points returns the total number of points in hand h.
// When hand h contains an ace that is counted as 11, soft is true.
func (h Hand) Points() (total int, soft bool) {
	aces := 0

	for _, c := range h {
		total += rankPoints[c.Rank]
		if c.Rank == card.Ace {
			aces++
		}
	}

	if total <= 21 {
		soft = true
	} else {
		if aces > 0 {
			for i := 0; i < aces; i++ {
				total = total - 10
				if total < 21 {
					break
				}
			}
		}
	}

	return total, soft
}
