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
	ui := &textUI{r: bufio.NewReader(os.Stdin), w: os.Stdout, err: handleError}
	ui.writeln("Welcome to blackjack!")

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
		if utf8.RuneCountInString(s) > 0 {
			r, _ := utf8.DecodeRuneInString(s)
			return unicode.ToLower(r)
		}

		ui.write("no user input received")
		ui.writeln()
	}
}

func (ui *textUI) write(args ...interface{}) {
	_, err := fmt.Fprint(ui.w, args...)
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

func (ui *textUI) writeln(args ...interface{}) {
	_, err := fmt.Fprintln(ui.w, args...)
	if err != nil {
		ui.err(err)
	}
}

const noDef = '\000'

func (ui *textUI) getRune(msg string, choices []string, accept []rune, def rune) rune {
	if len(choices) == 0 {
		panic("no choices provided")
	}
	if len(accept) == 0 {
		panic("no accept runes provided")
	}

	for i, r := range accept {
		accept[i] = unicode.ToLower(r)
	}

	for {
		ui.writef("%s %s: ", msg, strings.Join(choices, "/"))
		r := ui.readRune()
		if def != noDef && r == '\n' {
			r = def
		}

		for _, ar := range accept {
			if r == ar {
				return r
			}
		}

		buf := new(bytes.Buffer)
		buf.WriteString("invalid input received (choices: ")
		buf.WriteRune(toUpperIfMatch(accept[0], def))
		for i := 1; i < len(accept); i++ {
			buf.WriteString(", ")
			buf.WriteRune(toUpperIfMatch(accept[i], def))
		}
		buf.WriteByte(')')
		fmt.Fprintln(buf)

		_, err := buf.WriteTo(ui.w)
		if err != nil {
			ui.err(err)
		}
	}
}

func toUpperIfMatch(r, match rune) rune {
	if r == match {
		r = unicode.ToUpper(r)
	}
	return r
}

func (ui *textUI) getDecimal(msg string, prev decimal.Decimal, must bool) decimal.Decimal {
	for {
		zero := prev.Equal(decimal.Zero)
		if zero {
			ui.write(msg + " ")
		} else {
			ui.writef("%s [%v] ", msg, prev)
		}

		s := ui.readString()
		if s == "" {
			if must && zero {
				ui.writeln("no input received")
				continue
			}
			if zero {
				return decimal.Zero
			}
			return prev
		}

		d, err := decimal.NewFromString(s)
		if err != nil {
			ui.writeln("error:", err)
			continue
		}

		return d
	}
}

func (ui *textUI) writeFortune(f *player.Fortune) {
	ui.writeln()
	ui.writef("Active: %v\tSavings: %v\tStake: %v\n", f.Active(), f.Savings(), f.Stake())
}

func (ui *textUI) Bet(f *player.Fortune) decimal.Decimal {
	ui.writeFortune(f)

	if f.Active().Cmp(ui.bet) == -1 { // f.Active <= ui.bet
		ui.bet = decimal.Zero
	}

	amount := ui.getDecimal("How much do you want to bet?", ui.bet, true)

	if f.Active().Cmp(amount) >= 0 && ui.bet.Cmp(amount) <= 0 ||
		ui.bet.Equal(decimal.Zero) {

		// f.Active() >= amount && ui.bet <= amount || ui.bet == 0
		ui.bet = amount
	}

	return amount
}

func (ui *textUI) Hand(d, p blackjack.Hand) {
	ui.writeln()
	ui.writeln("Dealer:", d)
	ui.writeln("Player:", p)
}

func (ui *textUI) DealerCard(c card.Card, h blackjack.Hand) {
	if len(h) == 2 {
		ui.writeln()
	}
	ui.writeln("Dealer:", h)
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

	ar := ui.getRune("What is your next action?", s, r, noDef)
	for a, r := range actionRunes {
		if r == ar {
			return a
		}
	}

	panic("invalid user input: " + string(ar))
}

func (ui *textUI) SplitHand(lh, rh blackjack.Hand, a decimal.Decimal) {
	ui.writeln()
	ui.writeln("Hand splitted:")
	ui.writeln("first: ", lh)
	ui.writeln("second:", rh)
}

func (ui *textUI) DoubleHand(h blackjack.Hand, a decimal.Decimal) {
	ui.writeln()
	ui.writeln("Hand doubled:", h)
}

func (ui *textUI) Outcome(o blackjack.Outcome, a decimal.Decimal, d, p blackjack.Hand) {
	ui.writeln()
	ui.writeln("Outcome:", o)
	ui.writeln("Dealer: ", d)
	ui.writeln("Player: ", p)
	ui.writeln("Amount: ", a)
}

func (ui *textUI) NewGame(f *player.Fortune) bool {
	ui.writeFortune(f)
	choices := []string{"[Y]es", "[N]o"}
	runes := []rune{'y', 'n'}
	r := ui.getRune("Do you want to play a new game?", choices, runes, 'y')
	return r == 'y'
}

func (ui *textUI) NoActiveFortune(f *player.Fortune) bool {
	ui.writeFortune(f)
	return false
}

func (ui *textUI) NoFortune() {
	ui.writeln("No more fortune to bet with, sorry!")
}

func (ui *textUI) PerfectPairBet(f *player.Fortune) decimal.Decimal {
	const msg = "How much do you want to bet for Perfect Pair?"
	amount := ui.getDecimal(msg, ui.pp, false)
	ui.pp = amount
	return amount
}

func (ui *textUI) PerfectPair(kind blackjack.PerfectPair, a decimal.Decimal) {
	ui.writeln()
	ui.writeln("Perfect Pair")
	ui.writeln("Kind:  ", kind)
	ui.writeln("Amount:", a)
}
