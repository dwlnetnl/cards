// Code generated by "stringer -type=Outcome"; DO NOT EDIT

package blackjack

import "fmt"

const _Outcome_name = "WonLostBustPushedSurrenderedBlackjackDealerBlackjack"

var _Outcome_index = [...]uint8{0, 3, 7, 11, 17, 28, 37, 52}

func (i Outcome) String() string {
	if i < 0 || i >= Outcome(len(_Outcome_index)-1) {
		return fmt.Sprintf("Outcome(%d)", i)
	}
	return _Outcome_name[_Outcome_index[i]:_Outcome_index[i+1]]
}