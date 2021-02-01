package ui

import (
	"github.com/w32blaster/tax-bookkeeper/conf"
	"github.com/w32blaster/tax-bookkeeper/db"
	"github.com/w32blaster/tax-bookkeeper/tax"
	"time"
)

// accountingDateStart is only day and month, like 01-11
func CollectDataForDashboard(d *db.Database, accountingDateStart time.Time) *DashboardData {

	// TODO: handle errors!

	// get the profit for the current accounting period since accountingDateStart until now
	revenue, _ := d.GetRevenueSince(accountingDateStart)
	expenses, _ := d.GetExpensesSince(accountingDateStart)
	pension, _ := d.GetPensionSince(accountingDateStart)

	profit := revenue - expenses - pension

	// Corporate Tax
	corpTax := tax.CalculateCorporateTax(profit, accountingDateStart)

	// last 10 transactions
	lastTenTransactions, _ := d.GetAll(30)

	// Self-assessment
	movedOut, _ := d.GetMovedOut(accountingDateStart /* TODO:  ????? what is a period? */)
	selfAssessmentTax := tax.CalculateSelfAssessmentTax(movedOut, 0 /* TODO: ??? what income is??? */)
	rate, leftBeforeThreshold := tax.HowMuchBeforeNextThreshold(movedOut)

	return &DashboardData{

		// Corp Tax
		Period:                   tax.GetFinYear(accountingDateStart),
		NextPaymentDate:          time.Date(accountingDateStart.Year()+1, accountingDateStart.Month(), accountingDateStart.Day(), 0, 0, 0, 0, conf.GMT),
		CorporateTaxSoFar:        corpTax,
		EarnedAccountingPeriod:   revenue,
		ExpensesAccountingPeriod: expenses,
		PensionAccountingPeriod:  pension,

		LastTransactions: lastTenTransactions,

		// Self-assessment
		MovedOutFromCompanyTotal:   movedOut,
		SelfAssessmentTaxSoFar:     selfAssessmentTax,
		TaxRate:                    rate,
		HowMuchBeforeNextThreshold: leftBeforeThreshold,
	}
}
