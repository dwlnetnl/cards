package blackjack

import "github.com/shopspring/decimal"

// DoubleRule represents different double rule variants.
type DoubleRule int

// When is it possible for the player to double their hand?
const (
	DoubleAny DoubleRule = iota
	DoubleOnly9_10_11
	DoubleOnly10_11
)

//go:generate stringer -type=DoubleRule

// SurrenderRule represents different surrender rule variants.
type SurrenderRule int

// Surrender rule options.
const (
	NoSurrender SurrenderRule = iota
	EarlySurrender
	LateSurrender
)

//go:generate stringer -type=SurrenderRule

// Rules represent the game rules and mechanics.
type Rules interface {
	DealerHitSoft17() bool
	NumDecks() uint
	Surrender() SurrenderRule
	CanSplit([]Hand) bool
	Double() DoubleRule
	DoubleAfterSplit() bool
	BlackjackAfterSplit() bool
	NoHoleCard() bool
	OriginalBetsOnly() bool
	BlackjackRatio() decimal.Decimal
	DealerWinsTie() bool

	PerfectPair() bool
	PerfectPairRatio() (mixed, same, perfect int)
}

// Game rules in different casino's.
var (
	HollandCasino Rules = holland{}
	TapTapBoom    Rules = tapTapBoom{}
)

// AvailableRules is a slice that hold all available game rules.
var AvailableRules = []Rules{
	HollandCasino,
	TapTapBoom,
}

type holland struct{}

func (holland) NumDecks() uint                  { return 6 }
func (holland) DealerHitSoft17() bool           { return true }
func (holland) Surrender() SurrenderRule        { return NoSurrender }
func (holland) CanSplit([]Hand) bool            { return true }
func (holland) Double() DoubleRule              { return DoubleOnly9_10_11 }
func (holland) DoubleAfterSplit() bool          { return true }
func (holland) BlackjackAfterSplit() bool       { return true }
func (holland) NoHoleCard() bool                { return true }
func (holland) OriginalBetsOnly() bool          { return false }
func (holland) BlackjackRatio() decimal.Decimal { return decimal.New(15, -1) }
func (holland) DealerWinsTie() bool             { return false }
func (holland) PerfectPair() bool               { return true }
func (holland) PerfectPairRatio() (m, s, p int) { return 6, 12, 25 }

type tapTapBoom struct{}

func (tapTapBoom) NumDecks() uint                  { return 6 } // guess
func (tapTapBoom) DealerHitSoft17() bool           { return true }
func (tapTapBoom) Surrender() SurrenderRule        { return NoSurrender }
func (tapTapBoom) CanSplit(h []Hand) bool          { return len(h) == 1 }
func (tapTapBoom) Double() DoubleRule              { return DoubleAny }
func (tapTapBoom) DoubleAfterSplit() bool          { return true }
func (tapTapBoom) BlackjackAfterSplit() bool       { return true }
func (tapTapBoom) NoHoleCard() bool                { return false }
func (tapTapBoom) OriginalBetsOnly() bool          { return false }
func (tapTapBoom) BlackjackRatio() decimal.Decimal { return decimal.New(15, -1) }
func (tapTapBoom) DealerWinsTie() bool             { return false }
func (tapTapBoom) PerfectPair() bool               { return false }
func (tapTapBoom) PerfectPairRatio() (m, s, p int) { return }
