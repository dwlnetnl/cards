// Code generated by "stringer -type=PerfectPair"; DO NOT EDIT

package blackjack

import "fmt"

const _PerfectPair_name = "NoPairMixedSamePerfect"

var _PerfectPair_index = [...]uint8{0, 6, 11, 15, 22}

func (i PerfectPair) String() string {
	if i < 0 || i >= PerfectPair(len(_PerfectPair_index)-1) {
		return fmt.Sprintf("PerfectPair(%d)", i)
	}
	return _PerfectPair_name[_PerfectPair_index[i]:_PerfectPair_index[i+1]]
}
