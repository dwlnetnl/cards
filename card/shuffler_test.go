package card

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestShuffler(t *testing.T) {
	deck := NewStandardDeck()
	ref := refbuckNew(len(deck))
	for i, c := range deck {
		refbuckShuffle(&ref, c, refBucketsFill[i])
	}
	s := NewSeededShuffler(deck, 1, rand.NewSource(1))
	testBuckets(t, s.buck, ref)

	var c Card
	var drawn []Card

	for i := 0; i < 4; i++ {
		c = s.MustDraw()
		drawn = append(drawn, c)
		testCard(t, c, refbuckDraw(&ref, refBucketsDraw[i]))
	}

	testBuckets(t, s.buck, ref)

	for i, c := range drawn {
		refbuckShuffle(&ref, c, refBucketsReFill[i])
	}
	s.Shuffle(drawn...)
	testBuckets(t, s.buck, ref)

	defer func() {
		if recover() == nil {
			t.Fatal("MustDraw did not panic")
		}
	}()

	s.cards = 0 // let MustDraw panic
	s.MustDraw()
}

func testBuckets(t *testing.T, got, want [][]Card) {
	if !reflect.DeepEqual(got, want) {
		t.Fatal("shuffler buckets differ")
	}
}

func testCard(t *testing.T, got, want Card) {
	if got != want {
		t.Fatalf("got: %v, want: %v", got, want)
	}
}

type buckets [][]Card

func refbuckNew(cards int) [][]Card {
	n := bucketsForDeck(cards, 1)
	b := make([][]Card, n)
	for i := 0; i < n; i++ {
		b[i] = make([]Card, 0, bucketSize)
	}
	return b
}

func refbuckDraw(b *[][]Card, bucket int) Card {
	s := *b
	c := s[bucket][0]
	s[bucket] = append(s[bucket][:1], s[bucket][2:]...)
	return c
}

func refbuckShuffle(b *[][]Card, c Card, index int) {
	s := *b
	s[index] = append(s[index], c)
}

var refBucketsFill = []int{
	1,
	1,
	9,
	3,
	2,
	4,
	7,
	6,
	6,
	6,
	4,
	0,
	10,
	0,
	1,
	6,
	7,
	5,
	8,
	5,
	8,
	0,
	10,
	10,
	9,
	0,
	5,
	6,
	2,
	9,
	0,
	7,
	5,
	1,
	8,
	0,
	1,
	6,
	1,
	4,
	6,
	6,
	5,
	4,
	2,
	4,
	5,
	0,
	5,
	4,
	4,
	7,
}

var refBucketsDraw = []int{
	8,
	5,
	7,
	7,
}

var refBucketsReFill = []int{
	2,
	10,
	0,
	1,
}

func BenchmarkShuffler3(b *testing.B)  { benchmarkShuffler(b, 3) }
func BenchmarkShuffler6(b *testing.B)  { benchmarkShuffler(b, 6) }
func BenchmarkShuffler8(b *testing.B)  { benchmarkShuffler(b, 8) }
func BenchmarkShuffler10(b *testing.B) { benchmarkShuffler(b, 10) }
func BenchmarkShuffler15(b *testing.B) { benchmarkShuffler(b, 15) }

func benchmarkShuffler(b *testing.B, cards int) {
	b.Run("1", func(b *testing.B) { benchmarkShuffle(b, cards, 1) })
	b.Run("4", func(b *testing.B) { benchmarkShuffle(b, cards, 4) })
}

func benchmarkShuffle(b *testing.B, cards, rounds int) {
	s := NewShuffler(NewStandardDeck(), 1)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < rounds; j++ {
			drawn := make([]Card, cards)
			for k := 0; k < cards; k++ {
				drawn[k] = s.MustDraw()
			}
			s.Shuffle(drawn...)
		}
	}
}
