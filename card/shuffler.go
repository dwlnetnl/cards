package card

import (
	"math/rand"
	"time"
)

// A Shuffler holds one or more shuffled decks of playing cards.
type Shuffler struct {
	rand  *rand.Rand
	buck  [][]Card
	nbuck int
	cards int
}

const bucketSize = 8

func bucketsForDeck(size, num int) int {
	n := size
	n *= num
	n *= 5
	n /= 3
	n /= bucketSize
	return n + 1
}

// NewShuffler returns a shuffler that shuffles a number of particular decks.
func NewShuffler(d Deck, num uint) *Shuffler {
	return NewSeededShuffler(d, num, rand.NewSource(time.Now().UnixNano()))
}

// NewSeededShuffler returns a shuffler that shuffles a number of particular
// decks. The shuffler is seeded by a random source src.
func NewSeededShuffler(d Deck, num uint, src rand.Source) *Shuffler {
	n := bucketsForDeck(len(d), int(num))
	s := &Shuffler{
		rand:  rand.New(src),
		buck:  make([][]Card, n),
		nbuck: n,
	}

	for i := 0; i < n; i++ {
		s.buck[i] = make([]Card, 0, bucketSize)
	}

	for i := 0; i < int(num); i++ {
		s.Shuffle(d...)
	}

	return s
}

// Shuffle shuffles zero or more cards back into the deck(s).
func (s *Shuffler) Shuffle(cards ...Card) {
	for _, c := range cards {
		s.shuffle(c)
	}
}

func (s *Shuffler) shuffle(c Card) {
	for {
		i := s.rand.Intn(s.nbuck)
		if len(s.buck[i]) < bucketSize {
			s.buck[i] = append(s.buck[i], c)
			s.cards++
			return
		}
	}
}

// Draw draws randomly a card from the deck(s).
// It returns false if no card is left.
func (s *Shuffler) Draw() (c Card, ok bool) {
	if s.cards > 0 {
		c = s.draw()
		ok = true
	}
	return
}

// MustDraw is like Draw but panics if there is no  card left.
func (s *Shuffler) MustDraw() Card {
	c, ok := s.Draw()
	if !ok {
		panic("card: no cards left in shuffler")
	}
	return c
}

func (s *Shuffler) draw() Card {
	for {
		i := s.rand.Intn(s.nbuck)
		n := len(s.buck[i])
		if n > 0 {
			c := s.buck[i][0]
			if n == 1 {
				s.buck[i] = s.buck[i][1:]
			} else {
				s.buck[i] = append(s.buck[i][:1], s.buck[i][2:]...)
			}
			s.cards--
			return c
		}
	}
}
