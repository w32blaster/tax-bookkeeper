package ui

import (
	"time"

	"github.com/w32blaster/tax-bookkeeper/db"
	"github.com/w32blaster/tax-bookkeeper/tax"
)

// struct that keeps all the data for the dashboard
type (
	CorporateTax struct {
		Period                   string
		StartingDate             time.Time
		EndingDate               time.Time
		NextPaymentDate          time.Time
		CorporateTaxSoFar        float64
		EarnedAccountingPeriod   float64
		ExpensesAccountingPeriod float64
		PensionAccountingPeriod  float64
	}

	// TODO: Salary, dividends?
	SelfAssessmentTax struct {
		StartingDate             time.Time
		EndingDate               time.Time
		NextPaymentDate          time.Time
		MovedOutFromCompanyTotal float64
		SelfAssessmentTaxSoFar   float64
		TaxRate                  tax.Rate
		// warning:
		HowMuchBeforeNextThreshold float64

		IsWarning bool
	}

	VAT struct {
		Since                   time.Time
		NextVATToBePaidSoFar    float64
		NextDateYouShouldPayFor time.Time
		NextMonthSubmit         string
	}

	DashboardData struct {

		// Last 10 transactions
		LastTransactions []db.Transaction

		PreviousPeriod CorporateTax
		CurrentPeriod  CorporateTax

		PreviousSelfAssessmentPeriod SelfAssessmentTax
		CurrentSelfAssessmentPeriod  SelfAssessmentTax

		VAT
	}
)

// callback function that will be fired on the Save button clicking
type FuncAllocateTransactions func(txToAllocate map[int]db.TransactionCategory) error

// UI is a common interface for an GUI. At this moment we have only terminal UI,
// but if in the future we will need to do another UI, it would be easy possible
// to do by implementing this interface
type UI interface {
	Start()
	BeginDialogToAllocateTransactions(unallocatedTxs []db.Transaction, fnAllocate FuncAllocateTransactions)
	ShowDashboard(data DashboardData)
}
