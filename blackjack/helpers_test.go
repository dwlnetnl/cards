package blackjack

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/dwlnetnl/cards/card"
	"github.com/dwlnetnl/cards/player"

	"github.com/shopspring/decimal"
)

func seededShuffler(seed int64, fn func()) {
	old := newShuffler
	newShuffler = func(d card.Deck, num uint) *card.Shuffler {
		return card.NewSeededShuffler(d, num, rand.NewSource(seed))
	}
	fn()
	newShuffler = old
}

func testPlay(t *testing.T, seed int64, rules Rules, bet, pp int64, want []event) {
	seededShuffler(seed, func() {
		f := player.NewFortune(decimal.New(50, 0))
		ui := &testUI{test: t, want: want, fort: f, bet: bet, pp: pp}
		Play(ui, rules, f)
	})
}

type testUI struct {
	test *testing.T
	want []event
	idx  int
	end  bool
	bet  int64
	pp   int64
	bal  decimal.Decimal
	fort *player.Fortune
}

func (ui *testUI) last() event {
	if ui.idx-1 < 0 {
		return nil
	}
	return ui.want[ui.idx]
}

func (ui *testUI) lastAction() Action {
	for i := ui.idx - 1; i > 0; i-- {
		if e, ok := ui.want[i].(nextAction); ok {
			return e.next
		}
	}
	ui.test.Fatal("no last action found")
	return 0
}

func (ui *testUI) nextAction() Action {
	if ui.idx == len(ui.want) {
		ui.test.Fatal("no event to select next action from")
	}
	return ui.want[ui.idx].(nextAction).next
}

func (ui *testUI) check(got event) {
	if ui.idx == len(ui.want) {
		ui.test.Fatalf("unexpected event: %#v", got)
	}
	got.test(ui.test, ui.idx+1, ui.want[ui.idx])
	ui.idx++
}

func (ui *testUI) Bet(fortune *player.Fortune) decimal.Decimal {
	if fortune != ui.fort {
		ui.test.Fatalf("got fortune %v, want: %v", fortune, ui.fort)
	}
	amount := decimal.New(ui.bet, 0)
	ui.bal = fortune.Active().Sub(amount)
	return amount
}

func (ui *testUI) Hand(dealer, player Hand) {
	ui.check(hand{dealer, player})
}

func (ui *testUI) DealerCard(card card.Card, hand Hand) {
	ui.check(dealerCard{card, hand})
}

func (ui *testUI) NextAction(actions []Action) Action {
	next := ui.nextAction()
	ui.check(nextAction{actions: actions})
	return next
}

func (ui *testUI) SplitHand(left, right Hand, amount decimal.Decimal) {
	if ui.lastAction() != Split {
		ui.test.Fatal("last action was not Split")
	}
	ui.check(splitHand{left, right})
	ui.bal = ui.bal.Sub(amount)
}

func (ui *testUI) DoubleHand(hand Hand, amount decimal.Decimal) {
	if ui.lastAction() != Double {
		ui.test.Fatal("last action was not Double")
	}
	ui.check(doubleHand{hand, amount})
	ui.bal = ui.bal.Sub(amount)
}

func (ui *testUI) Outcome(out Outcome, amount decimal.Decimal, dealer, player Hand) {
	ui.check(outcome{out, amount, dealer, player})
	if out == Won || out == Pushed || out == Blackjack || out == Surrendered {
		ui.bal = ui.bal.Add(amount)
	}
	ui.end = true
}

func (ui *testUI) NewGame() bool {
	if len(ui.want) > 0 {
		if !ui.end {
			ui.test.Error("game had no outcome")
			ui.test.Fail()
		}
		if !ui.bal.Equal(ui.fort.Active()) {
			ui.test.Errorf("got balance %v, want: %v", ui.bal, ui.fort.Active())
			ui.test.Fail()
		}
	}
	return false
}

func (ui *testUI) PerfectPairBet(fortune *player.Fortune) decimal.Decimal {
	if fortune != ui.fort {
		ui.test.Fatalf("got fortune %v, want: %v", fortune, ui.fort)
	}
	amount := decimal.New(ui.pp, 0)
	ui.bal = ui.bal.Sub(amount)
	return amount
}

func (ui *testUI) PerfectPair(kind PerfectPair, amount decimal.Decimal) {
	ui.check(perfectPair{kind, amount})
	ui.bal = ui.bal.Add(amount)
}

type event interface {
	test(t *testing.T, num int, other event)
}

type hand struct {
	dealer Hand
	player Hand
}

func (want hand) test(t *testing.T, num int, other event) {
	got := other.(hand)
	if !reflect.DeepEqual(got.dealer, want.dealer) {
		t.Errorf("#%d: got dealer hand %v, want: %v", num, got.dealer, want.dealer)
		t.Fail()
	}
	if !reflect.DeepEqual(got.player, want.player) {
		t.Errorf("#%d: got player hand %v, want: %v", num, got.player, want.player)
		t.Fail()
	}
}

type dealerCard struct {
	card card.Card
	hand Hand
}

func (want dealerCard) test(t *testing.T, num int, other event) {
	got := other.(dealerCard)
	if got.card != want.card {
		t.Errorf("#%d: got card %v, want: %v", num, got.card, want.card)
		t.Fail()
	}
	if !reflect.DeepEqual(got.hand, want.hand) {
		t.Errorf("#%d: got hand %v, want: %v", num, got.hand, want.hand)
		t.Fail()
	}
}

type nextAction struct {
	actions []Action
	next    Action
}

func (want nextAction) test(t *testing.T, num int, other event) {
	got := other.(nextAction)
	if !reflect.DeepEqual(got.actions, want.actions) {
		t.Fatalf("#%d: got actions %v, want: %v", num, got.actions, want.actions)
	}
}

type splitHand struct {
	left  Hand
	right Hand
}

func (want splitHand) test(t *testing.T, num int, other event) {
	got := other.(splitHand)
	if !reflect.DeepEqual(got.left, want.left) {
		t.Errorf("#%d: got left hand %v, want: %v", num, got.left, want.left)
		t.Fail()
	}
	if !reflect.DeepEqual(got.right, want.right) {
		t.Errorf("#%d: got right hand %v, want: %v", num, got.right, want.right)
		t.Fail()
	}
}

type doubleHand struct {
	hand      Hand
	withdrawn decimal.Decimal
}

func (want doubleHand) test(t *testing.T, num int, other event) {
	got := other.(doubleHand)
	if !reflect.DeepEqual(got.hand, want.hand) {
		t.Errorf("#%d: got hand %v, want: %v", num, got.hand, want.hand)
		t.Fail()
	}
	if !got.withdrawn.Equal(want.withdrawn) {
		t.Errorf("#%d: got withdrawn %v, want: %v", num, got.withdrawn, want.withdrawn)
		t.Fail()
	}
}

type outcome struct {
	outcome Outcome
	amount  decimal.Decimal
	dealer  Hand
	player  Hand
}

func (want outcome) test(t *testing.T, num int, other event) {
	got := other.(outcome)
	if got.outcome != want.outcome {
		t.Errorf("#%d: got outcome %v, want: %v", num, got.outcome, want.outcome)
		t.Fail()
	}
	if !got.amount.Equal(want.amount) {
		t.Errorf("#%d: got amount %v, want: %v", num, got.amount, want.amount)
		t.Fail()
	}
	if !reflect.DeepEqual(got.dealer, want.dealer) {
		t.Errorf("#%d: got dealer hand %v, want: %v", num, got.dealer, want.dealer)
		t.Fail()
	}
	if !reflect.DeepEqual(got.player, want.player) {
		t.Errorf("#%d: got player hand %v, want: %v", num, got.player, want.player)
		t.Fail()
	}
}

type perfectPair struct {
	kind   PerfectPair
	amount decimal.Decimal
}

func (want perfectPair) test(t *testing.T, num int, other event) {
	got := other.(perfectPair)
	if got.kind != want.kind {
		t.Errorf("#%d: got kind %v, want: %v", num, got.kind, want.kind)
		t.Fail()
	}
	if !got.amount.Equal(want.amount) {
		t.Errorf("#%d: got amount %v, want: %v", num, got.amount, want.amount)
		t.Fail()
	}
}
