// Package card defines a playing card data type.
package card

// Suit represents a playing card suit.
type Suit int

// Playing card suits.
const (
	Naked Suit = iota
	Spades
	Hearts
	Diamonds
	Clubs
)

var suitNames = [...]string{"Naked", "Spades", "Hearts", "Diamonds", "Clubs"}
var suitRunes = [...]rune{' ', '\u2664', '\u2661', '\u2662', '\u2667'}

// Symbol returns the unicode rune for suit s.
func (s Suit) Symbol() rune   { return suitRunes[s] }
func (s Suit) String() string { return suitNames[s] }

// Rank represents a playing card rank.
type Rank int

// Playing card ranks.
const (
	Two Rank = iota
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace

	JokerRed
	JokerBlack
	JokerWhite
)

var rankNames = [...]string{
	"Two", "Three", "Four", "Five", "Six", "Seven", "Eight",
	"Nine", "Ten", "Jack", "Queen", "King", "Ace",
	"Red Joker", "Black Joker", "White Joker",
}

var rankSymbols = [...]string{
	"2", "3", "4", "5", "6", "7", "8",
	"9", "10", "J", "Q", "K", "A",
	"*R", "*B", "*W",
}

// Symbol returns a short identifier for rank r.
func (r Rank) Symbol() string { return rankSymbols[r] }
func (r Rank) String() string { return rankNames[r] }

// Card represents a playing card from a particular deck.
type Card struct {
	Suit Suit
	Rank Rank
}

// Color represents a playing card color.
type Color int

// Colors of playing cards.
const (
	Colorless Color = iota
	Red
	Black
)

// Color returns the color of playing card c.
func (c Card) Color() Color {
	switch c.Suit {
	case Hearts, Diamonds:
		return Red
	case Spades, Clubs:
		return Black
	default:
		switch c.Rank {
		case JokerRed:
			return Red
		case JokerBlack:
			return Black
		default:
			return Colorless
		}
	}
}

func (c Card) String() string {
	if c.Suit == Naked {
		return c.Rank.Symbol()
	}
	return string(c.Suit.Symbol()) + " " + c.Rank.Symbol()
}

func nonJokerRank(r Rank) Rank {
	if r == JokerRed || r == JokerBlack || r == JokerWhite {
		panic("card: joker in normal rank is not allowed")
	}

	return r
}

// Spade returns a spades card with rank r.
func Spade(r Rank) Card { return Card{Spades, nonJokerRank(r)} }

// Heart returns a hearts card with rank r.
func Heart(r Rank) Card { return Card{Hearts, nonJokerRank(r)} }

// Diamond returns a diamonds card with rank r.
func Diamond(r Rank) Card { return Card{Diamonds, nonJokerRank(r)} }

// Club returns a clubs card with rank r.
func Club(r Rank) Card { return Card{Clubs, nonJokerRank(r)} }

// Joker returns an unspecified kind of joker.
func Joker() Card { return RedJoker() }

// RedJoker returns a red joker.
func RedJoker() Card { return Card{Naked, JokerRed} }

// BlackJoker returns a black joker.
func BlackJoker() Card { return Card{Naked, JokerBlack} }

// WhiteJoker returns a white joker.
func WhiteJoker() Card { return Card{Naked, JokerWhite} }

// Deck represents a deck of playing cards.
type Deck []Card

// NewStandardDeck returns a new standard 52-cards deck of playing cards.
func NewStandardDeck() Deck {
	const suits = 4
	const ranks = 13

	d := make(Deck, suits*ranks)
	for s := 0; s < suits; s++ {
		for r := 0; r < ranks; r++ {
			d[s*ranks+r] = Card{Suit(s + 1), Rank(r)}
		}
	}

	return d
}
