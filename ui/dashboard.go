package ui

import (
	"math"
	"time"

	"github.com/w32blaster/tax-bookkeeper/conf"
	"github.com/w32blaster/tax-bookkeeper/db"
	"github.com/w32blaster/tax-bookkeeper/tax"
)

// CollectDataForDashboard accountingDateStart is only day and month, like 01-11
func CollectDataForDashboard(d *db.Database, accountingDateStart time.Time, vatMonth time.Month) (*DashboardData, error) {

	now := time.Now().In(conf.GMT)

	// get the profit for the current accounting period since accountingDateStart until now
	currentCorporateTax, err := collectSummaryCorporateTax(d, accountingDateStart, now)
	if err != nil {
		return nil, err
	}

	previousCorporateTax, err := collectSummaryCorporateTax(d,
		accountingDateStart.AddDate(-1, 0, 0),
		accountingDateStart)
	if err != nil {
		return nil, err
	}

	// last 10 transactions
	lastTenTransactions, _ := d.GetAll(30)

	// Self-assessment tax
	currentSelfAssessmentTax, err := collectSummarySelfAssessmentTax(d, now)
	if err != nil {
		return nil, err
	}

	previousSelfAssessmentTax, err := collectSummarySelfAssessmentTax(d, now.AddDate(-1, 0, 0))
	if err != nil {
		return nil, err
	}

	currentVAT, err := collectSummaryVAT(d, vatMonth, now)
	previousVAT, err := collectSummaryVAT(d, vatMonth, now.AddDate(0, -3, 0))

	if err != nil {
		return nil, err
	}

	return &DashboardData{
		PreviousPeriod:   previousCorporateTax,
		CurrentPeriod:    currentCorporateTax,
		LastTransactions: lastTenTransactions,

		PreviousSelfAssessmentPeriod: previousSelfAssessmentTax,
		CurrentSelfAssessmentPeriod:  currentSelfAssessmentTax,

		PreviousVAT: previousVAT,
		CurrentVAT:  currentVAT,
	}, nil
}

func collectSummarySelfAssessmentTax(d *db.Database, now time.Time) (SelfAssessmentTax, error) {

	startDate, endDate, paymentDate := tax.GetTaxYearDates(now)

	movedOut, _ := d.GetMovedOut(startDate, endDate)
	selfAssessmentTax := tax.CalculateSelfAssessmentTax(movedOut, 0)
	rate, leftBeforeThreshold, isWarning := tax.HowMuchBeforeNextThreshold(math.Abs(movedOut))

	return SelfAssessmentTax{
		StartingDate:               startDate,
		EndingDate:                 endDate,
		NextPaymentDate:            paymentDate,
		MovedOutFromCompanyTotal:   movedOut,
		SelfAssessmentTaxSoFar:     selfAssessmentTax,
		TaxRate:                    rate,
		HowMuchBeforeNextThreshold: leftBeforeThreshold,
		IsWarning:                  isWarning,
	}, nil
}

func collectSummaryCorporateTax(d *db.Database, accountingDateStart time.Time, accountingDateEnd time.Time) (CorporateTax, error) {

	var revenue, expenses, pension float64
	var err error

	if revenue, _ = d.GetRevenueSince(accountingDateStart, accountingDateEnd); err != nil {
		return CorporateTax{}, err
	}
	if expenses, err = d.GetExpensesSince(accountingDateStart, accountingDateEnd); err != nil {
		return CorporateTax{}, err
	}
	if pension, err = d.GetPensionSince(accountingDateStart, accountingDateEnd); err != nil {
		return CorporateTax{}, err
	}

	profit := revenue - expenses - pension

	// Corporate Tax
	corpTax := tax.CalculateCorporateTax(profit, accountingDateStart)

	// You must pay your Corporation Tax 9 months and 1 day after the end
	// of your accounting period
	// https://www.gov.uk/pay-corporation-tax
	paymentDate := accountingDateStart.AddDate(1, 9, 1)

	return CorporateTax{
		Period:                   tax.GetFinYear(accountingDateStart),
		StartingDate:             accountingDateStart,
		EndingDate:               accountingDateStart.AddDate(1, 0, -1),
		NextPaymentDate:          paymentDate,
		CorporateTaxSoFar:        corpTax,
		EarnedAccountingPeriod:   revenue,
		ExpensesAccountingPeriod: expenses,
		PensionAccountingPeriod:  pension,
	}, nil
}

func collectSummaryVAT(d *db.Database, vatMonth time.Month, until time.Time) (VAT, error) {

	submitMonth, payDeadline := tax.GetNextReturnDate(vatMonth, until)

	beginningVATPeriod := tax.GetBeginningOfPreviousPeriod(submitMonth, until.Year())
	vatExpenses, err := d.GetExpensesSince(beginningVATPeriod, until)
	if err != nil {
		return VAT{}, err
	}

	return VAT{
		Since:                   beginningVATPeriod,
		Until:                   beginningVATPeriod.AddDate(0, 3, -1),
		NextVATToBePaidSoFar:    vatExpenses * 0.2,
		NextDateYouShouldPayFor: payDeadline,
		NextMonthSubmit:         submitMonth.String(),
	}, nil
}
