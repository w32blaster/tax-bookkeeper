package ui

import (
	"github.com/w32blaster/tax-bookkeeper/db"
	"time"
)
type TaxRate int

const (
	PersonalAllowance  TaxRate = 1 + iota
	BasicRate
	HigherRate
	AdditionalRate
	)

	// struct that keeps all the data for the dashboard
type DashboardData struct {

	// Corporate tax
	Period string
	NextPaymentDate time.Time
	CorporateTaxSoFar float64
	EarnedAccountingPeriod float64
	ExpensesAccountingPeriod float64

	// VAT
	NextVATToBePaidSoFar float64

	// Self-Assessment tax
	MovedOutFromCompanyTotal float64
	// Salary, dividends?
	TaxRate TaxRate
	// warning:
	HowMuchBeforeNextThreshold float64
}

// callback function that will be fired on the Save button clicking
type FuncAllocateTransactions func (txToAllocate map[int]db.TransactionCategory) error


// UI is a common interface for an GUI. At this moment we have only terminal UI,
// but if in the future we will need to do another UI, it would be easy possible
// to do by implementing this interface
type UI interface {
	Start()
	BeginDialogToAllocateTransactions(unallocatedTxs []db.Transaction, fnAllocate FuncAllocateTransactions)
	ShowDashboard(data DashboardData)
}
