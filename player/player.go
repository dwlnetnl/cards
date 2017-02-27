// Package player provides card game player related data types.
package player

import "github.com/shopspring/decimal"

// Fortune represents a player assets.
type Fortune struct {
	stake  decimal.Decimal
	active decimal.Decimal
	saving decimal.Decimal
}

// NewFortune makes a fortune by taking a stake.
func NewFortune(stake decimal.Decimal) *Fortune {
	return &Fortune{stake: stake, active: stake}
}

// Stake returns the amount the fortune is started with.
func (f Fortune) Stake() decimal.Decimal { return f.stake }

// Active returns the amount of active assets.
func (f Fortune) Active() decimal.Decimal { return f.active }

// Savings returns the amount of saved assets.
func (f Fortune) Savings() decimal.Decimal { return f.saving }

// Total returns the total amount of assets.
func (f Fortune) Total() decimal.Decimal { return f.active.Add(f.saving) }

// Has returns true if an amount can be taken from fortune f.
func (f Fortune) Has(amount decimal.Decimal) bool {
	return f.active.Cmp(amount) >= 0 // f.active >= amount
}

// Skim transfers an amount from active to savings.
func (f *Fortune) Skim(amount decimal.Decimal) {
	f.active.Sub(amount)
	f.saving.Add(amount)
}

// Pour transfers an amount from savings to active.
func (f *Fortune) Pour(amount decimal.Decimal) {
	f.saving.Sub(amount)
	f.active.Add(amount)
}

// Withdrawal reduces the active fortune by an amount.
func (f *Fortune) Withdrawal(amount decimal.Decimal) {
	f.active = f.active.Sub(amount)
}

// Deposit increases the active fortune by an amount.
func (f *Fortune) Deposit(amount decimal.Decimal) {
	f.active = f.active.Add(amount)
}
