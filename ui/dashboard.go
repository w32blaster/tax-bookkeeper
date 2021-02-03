package ui

import (
	"github.com/w32blaster/tax-bookkeeper/conf"
	"github.com/w32blaster/tax-bookkeeper/db"
	"github.com/w32blaster/tax-bookkeeper/tax"
	"math"
	"time"
)

// accountingDateStart is only day and month, like 01-11
func CollectDataForDashboard(d *db.Database, accountingDateStart time.Time, vatMonth time.Month) (*DashboardData, error) {

	// get the profit for the current accounting period since accountingDateStart until now
	corporateTax, err := collectSummaryCorporateTax(d, accountingDateStart)
	if err != nil {
		return nil, err
	}

	// last 10 transactions
	lastTenTransactions, _ := d.GetAll(30)

	// Self-assessment tax
	selfAssessmentTax, err := collectSummarySelfAssessmentTax(d, accountingDateStart)
	if err != nil {
		return nil, err
	}

	vat, err := collectSummaryVAT(d, vatMonth)
	if err != nil {
		return nil, err
	}

	return &DashboardData{
		CorporateTax:      corporateTax,
		LastTransactions:  lastTenTransactions,
		SelfAssessmentTax: selfAssessmentTax,
		VAT:               vat,
	}, nil
}

func collectSummarySelfAssessmentTax(d *db.Database, accountingDateStart time.Time) (SelfAssessmentTax, error) {

	movedOut, _ := d.GetMovedOut(accountingDateStart /* TODO:  ????? what is a period? */)
	selfAssessmentTax := tax.CalculateSelfAssessmentTax(movedOut, 0 /* TODO: ??? what income is??? */)
	rate, leftBeforeThreshold, isWarning := tax.HowMuchBeforeNextThreshold(math.Abs(movedOut))

	return SelfAssessmentTax{
		MovedOutFromCompanyTotal:   movedOut,
		SelfAssessmentTaxSoFar:     selfAssessmentTax,
		TaxRate:                    rate,
		HowMuchBeforeNextThreshold: leftBeforeThreshold,
		IsWarning:                  isWarning,
	}, nil
}

func collectSummaryCorporateTax(d *db.Database, accountingDateStart time.Time) (CorporateTax, error) {

	var revenue, expenses, pension float64
	var err error

	if revenue, _ = d.GetRevenueSince(accountingDateStart); err != nil {
		return CorporateTax{}, err
	}
	if expenses, err = d.GetExpensesSince(accountingDateStart); err != nil {
		return CorporateTax{}, err
	}
	if pension, err = d.GetPensionSince(accountingDateStart); err != nil {
		return CorporateTax{}, err
	}

	profit := revenue - expenses - pension

	// Corporate Tax
	corpTax := tax.CalculateCorporateTax(profit, accountingDateStart)

	return CorporateTax{
		Period:                   tax.GetFinYear(accountingDateStart),
		NextPaymentDate:          time.Date(accountingDateStart.Year()+1, accountingDateStart.Month(), accountingDateStart.Day(), 0, 0, 0, 0, conf.GMT),
		CorporateTaxSoFar:        corpTax,
		EarnedAccountingPeriod:   revenue,
		ExpensesAccountingPeriod: expenses,
		PensionAccountingPeriod:  pension,
	}, nil
}

func collectSummaryVAT(d *db.Database, vatMonth time.Month) (VAT, error) {
	now := time.Now().In(conf.GMT)
	submitMonth, payDeadline := tax.GetNextReturnDate(vatMonth, now)

	beginningVATPeriod := tax.GetBeginningOfPreviousPeriod(submitMonth, now.Year())
	vatExpenses, err := d.GetExpensesSince(beginningVATPeriod)
	if err != nil {
		return VAT{}, err
	}

	return VAT{
		Since:                   beginningVATPeriod,
		NextVATToBePaidSoFar:    vatExpenses * 0.2,
		NextDateYouShouldPayFor: payDeadline,
		NextMonthSubmit:         submitMonth.String(),
	}, nil
}
