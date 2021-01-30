package tax

import (
	"github.com/w32blaster/tax-bookkeeper/conf"
	"github.com/w32blaster/tax-bookkeeper/db"
	"github.com/w32blaster/tax-bookkeeper/ui"
	"time"
)

// accountingDateStart is only day and month, like 01-11
func CollectDataForDashboard(d *db.Database, accountingDateStart time.Time) *ui.DashboardData {

	// get the profit for the current accounting period since accountingDateStart until now
	revenue, _ := d.GetRevenueSince(accountingDateStart)
	expenses, _ := d.GetExpensesSince(accountingDateStart)
	pension, _ := d.GetPensionSince(accountingDateStart)

	profit := revenue - expenses - pension

	// Corporate Tax
	corpTax := CalculateCorporateTax(profit, accountingDateStart)

	// last 10 transactions
	lastTenTransactions, _ := d.GetAll(30)

	// Self-assessment
	movedOut, _ := d.GetMovedOut(accountingDateStart /* ????? what is a period? */)
	selfAssesmentTax := CalculateSelfAssessmentTax(movedOut, 0 /* ??? what costs are??? */)

	return &ui.DashboardData{

		// Corp Tax
		Period:                   GetFinYear(accountingDateStart),
		NextPaymentDate:          time.Date(accountingDateStart.Year()+1, accountingDateStart.Month(), accountingDateStart.Day(), 0, 0, 0, 0, conf.GMT),
		CorporateTaxSoFar:        corpTax,
		EarnedAccountingPeriod:   revenue,
		ExpensesAccountingPeriod: expenses,
		PensionAccountingPeriod:  pension,

		LastTransactions: lastTenTransactions,

		// Self-assessment
		MovedOutFromCompanyTotal: movedOut,
		SelfAssessmentTaxSoFar:   selfAssesmentTax,
	}
}
