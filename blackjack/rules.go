package blackjack

import "math/big"

// DoubleRule represents different double rule variants.
type DoubleRule int

// When is it possible for the player to double their hand?
const (
	DoubleAny DoubleRule = iota
	DoubleOnly9_10_11
	DoubleOnly10_11
)

// SurrenderRule represents different surrender rule variants.
type SurrenderRule int

// Surrender rule options.
const (
	NoSurrender SurrenderRule = iota
	EarlySurrender
	LateSurrender
)

// Rules represent the game rules and mechanics.
type Rules interface {
	DealerHitSoft17() bool
	NumDecks() uint
	Surrender() SurrenderRule
	CanSplit([]Hand) bool
	Double() DoubleRule
	DoubleAfterSplit() bool
	NoHoleCard() bool
	OriginalBetsOnly() bool
	BlackjackRatio() *big.Rat
	DealerWinsTie() bool

	PerfectPair() bool
	PerfectPairRatio() (mixed, same, perfect int)
}

// Game rules in different casino's.
var (
	HollandCasino Rules = holland{}
	London        Rules = london{}
)

type holland struct{}

func (holland) NumDecks() uint                  { return 6 }
func (holland) DealerHitSoft17() bool           { return true }
func (holland) Surrender() SurrenderRule        { return NoSurrender }
func (holland) CanSplit([]Hand) bool            { return true }
func (holland) Double() DoubleRule              { return DoubleOnly9_10_11 }
func (holland) DoubleAfterSplit() bool          { return true }
func (holland) NoHoleCard() bool                { return true }
func (holland) OriginalBetsOnly() bool          { return false }
func (holland) BlackjackRatio() *big.Rat        { return big.NewRat(3, 2) }
func (holland) DealerWinsTie() bool             { return false }
func (holland) PerfectPair() bool               { return true }
func (holland) PerfectPairRatio() (m, s, p int) { return 6, 12, 25 }

type london struct{}

func (london) NumDecks() uint                  { return 6 }
func (london) DealerHitSoft17() bool           { return true }
func (london) Surrender() SurrenderRule        { return NoSurrender }
func (london) CanSplit(h []Hand) bool          { return len(h) == 1 }
func (london) Double() DoubleRule              { return DoubleAny }
func (london) DoubleAfterSplit() bool          { return true }
func (london) NoHoleCard() bool                { return false }
func (london) OriginalBetsOnly() bool          { return false }
func (london) BlackjackRatio() *big.Rat        { return big.NewRat(3, 2) }
func (london) DealerWinsTie() bool             { return false }
func (london) PerfectPair() bool               { return false }
func (london) PerfectPairRatio() (m, s, p int) { return }
