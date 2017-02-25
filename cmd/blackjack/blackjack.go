package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/dwlnetnl/cards/blackjack"
	"github.com/dwlnetnl/cards/card"
	"github.com/dwlnetnl/cards/player"

	"github.com/shopspring/decimal"
)

func main() {
	fmt.Println("Welcome to blackjack!")

	ui := &textUI{r: bufio.NewReader(os.Stdin), w: os.Stdout, err: handleError}
	f := player.NewFortune(decimal.New(50, 0))
	blackjack.Play(ui, blackjack.HollandCasino, f)
}

func handleError(err error) {
	fmt.Println("error during play:", err)
	os.Exit(1)
}

type textUI struct {
	r   *bufio.Reader
	w   io.Writer
	err func(error)
	bet decimal.Decimal
	pp  decimal.Decimal
}

func (ui *textUI) readString() string {
	s, err := ui.r.ReadString('\n')
	if err != nil {
		ui.err(err)
	}
	return strings.TrimSpace(s)
}

func (ui *textUI) readRune() rune {
	s := ui.readString()
	if s == "" {
		return '\n'
	}

	for {
		switch utf8.RuneCountInString(s) {
		default:
			ui.write("invalid input received")
			ui.writeln()
		case 0:
			ui.write("no user input received")
			ui.writeln()
		case 1:
			r, _ := utf8.DecodeRuneInString(s)
			return unicode.ToLower(r)
		}
	}
}

func (ui *textUI) write(s string) {
	_, err := fmt.Fprint(ui.w, s)
	if err != nil {
		ui.err(err)
	}
}

func (ui *textUI) writef(format string, args ...interface{}) {
	_, err := fmt.Fprintf(ui.w, format, args...)
	if err != nil {
		ui.err(err)
	}
}

func (ui *textUI) writeln() {
	_, err := fmt.Fprintln(ui.w)
	if err != nil {
		ui.err(err)
	}
}

func (ui *textUI) getRune(msg string, choices []string, accept []rune) rune {
	if len(choices) == 0 {
		panic("no choices provided")
	}
	if len(accept) == 0 {
		panic("no accept runes provided")
	}

	for {
		ui.writef("%s %s: ", msg, strings.Join(choices, "/"))
		r := ui.readRune()
		for i, r := range accept {
			accept[i] = unicode.ToLower(r)
		}

		for _, ar := range accept {
			if r == ar {
				return r
			}
		}

		buf := new(bytes.Buffer)
		buf.WriteString("invalid input received (choices: ")
		buf.WriteRune(accept[0])
		for i := 1; i < len(accept); i++ {
			buf.WriteString(", ")
			buf.WriteRune(accept[i])
		}
		buf.WriteByte(')')
		fmt.Fprintln(buf)

		_, err := buf.WriteTo(ui.w)
		if err != nil {
			ui.err(err)
		}
	}
}

func (ui *textUI) getDecimal(msg string, prev decimal.Decimal) decimal.Decimal {
	for {
		zero := prev.Equal(decimal.Zero)
		if zero {
			ui.write(msg + " ")
		} else {
			ui.writef("%s [%v] ", msg, prev)
		}

		s := ui.readString()
		if s == "" {
			if zero {
				return decimal.Zero
			}
			return prev
		}

		d, err := decimal.NewFromString(s)
		if err != nil {
			ui.writef("error: %v", err)
			ui.writeln()
			continue
		}

		return d
	}
}

func (ui *textUI) writeFortune(f *player.Fortune) {
	ui.writeln()
	ui.writef("Active: %v\tSavings: %v\tStake: %v", f.Active(), f.Savings(), f.Stake())
	ui.writeln()
}

func (ui *textUI) Bet(f *player.Fortune) decimal.Decimal {
	ui.writeFortune(f)
	amount := ui.getDecimal("How much do you want to bet?", ui.bet)
	ui.bet = amount
	return amount
}

func (ui *textUI) Hand(d, p blackjack.Hand) {
	ui.writef("d=%v p=%v", d, p)
	ui.writeln()
}

func (ui *textUI) DealerCard(c card.Card, h blackjack.Hand) {
	ui.writef("c=%v h=%v", c, h)
	ui.writeln()
}

var actionRunes = map[blackjack.Action]rune{
	blackjack.Hit:       'h',
	blackjack.Stand:     's',
	blackjack.Split:     't',
	blackjack.Double:    'd',
	blackjack.Surrender: 'r',
	blackjack.Continue:  'c',
}

var actionName = map[blackjack.Action]string{
	blackjack.Hit:       "[H]it",
	blackjack.Stand:     "[S]tand",
	blackjack.Split:     "Spli[t]",
	blackjack.Double:    "[D]ouble",
	blackjack.Surrender: "Sur[r]ender",
	blackjack.Continue:  "[C]ontinue",
}

func (ui *textUI) NextAction(actions []blackjack.Action) blackjack.Action {
	s := make([]string, len(actions))
	r := make([]rune, len(actions))
	for i, a := range actions {
		s[i] = actionName[a]
		r[i] = actionRunes[a]
	}

	ar := ui.getRune("What is your next action?", s, r)
	for a, r := range actionRunes {
		if r == ar {
			return a
		}
	}

	panic("invalid user input: " + string(ar))
}

func (ui *textUI) SplitHand(lh, rh blackjack.Hand, a decimal.Decimal) {
	ui.writef("lh=%v rh=%v a=%v", lh, rh, a)
	ui.writeln()
}

func (ui *textUI) DoubleHand(h blackjack.Hand, a decimal.Decimal) {
	ui.writef("h=%v a=%v", h, a)
	ui.writeln()
}

func (ui *textUI) Outcome(o blackjack.Outcome, a decimal.Decimal, d, p blackjack.Hand) {
	ui.writef("o=%v a=%v d=%v p=%v", o, a, d, p)
	ui.writeln()
}

func (ui *textUI) NewGame() bool {
	choices := []string{"[Y]es", "[N]o"}
	runes := []rune{'y', 'n'}
	r := ui.getRune("Do you want to play a new game?", choices, runes)
	return r == 'y'
}

func (ui *textUI) PerfectPairBet(f *player.Fortune) decimal.Decimal {
	amount := ui.getDecimal("How much do you want to bet for Perfect Pair?", ui.pp)
	ui.pp = amount
	return amount
}

func (ui *textUI) PerfectPair(pp blackjack.PerfectPair, a decimal.Decimal) {
	ui.writef("pp=%v a=%v", pp, a)
	ui.writeln()
}
