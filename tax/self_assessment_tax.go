package tax

import (
	"math"
)

const (
	personalAllowance = 12500.0
	weeksInAYear      = 52
)

type Rate int

const (
	PersonalAllowance Rate = 1 + iota
	BasicRate
	HigherRate
	AdditionalRate
)

func (r Rate) PrettyString() string {
	switch r {
	case PersonalAllowance:
		return "Personal Allowance (0%)"
	case BasicRate:
		return "Basic Rate (20%)"
	case HigherRate:
		return "Higher Rate (40%)"
	case AdditionalRate:
		return "Additional Rate (45%)"
	}
	return ""
}

// Tax Year is from 6 April to 5 April
// https://www.gov.uk/income-tax-rates
func CalculateSelfAssessmentTax(income, costs float64) float64 {

	profitBeforeTaxes := income - costs

	personalTax := getPersonalTaxFrom(profitBeforeTaxes)

	class2NITax, class4NITax := getNITax(profitBeforeTaxes)

	return personalTax + class2NITax + class4NITax
}

//    Band                    Taxable income         Tax rate
//    -------------           --------------         ---------
//    Personal Allowance      Up to £12,500          0%
//    Basic rate              £12,501 to £50,000     20%
//    Higher rate             £50,001 to £150,000    40%
//    Additional rate         over £150,000          45% (rich bastard!)
//
// please refer to unit tests for examples
//
func getPersonalTaxFrom(profitBeforeTaxes float64) float64 {

	// personalAllowance IS DIFFERENT FOR YEARS!!!
	// https://www.gov.uk/government/publications/rates-and-allowances-income-tax/income-tax-rates-and-allowances-current-and-past#tax-rates-and-bands
	if profitBeforeTaxes <= personalAllowance {
		return 0
	}

	allowance := getPersonalAllowance(profitBeforeTaxes)

	// Basic rate (£12,501 to £50,000) - 20%
	taxableProfit := profitBeforeTaxes - allowance
	if profitBeforeTaxes <= 50000 {
		return taxableProfit * 0.2
	}

	// Higher rate (£50,001 to £150,000) - 40%
	if profitBeforeTaxes <= 150000 {
		return 37500*0.2 + (taxableProfit-37500)*0.4
	}

	// Additional rate (over £150,000) - 45%
	return 37500*0.2 + (100000+12500)*0.4 + (profitBeforeTaxes-150000)*0.45
}

// Anyone earning more than £100,000 per year will have their personal
// allowance reduced by £1 for every £2 over that threshold until
// the personal allowance reaches zero.
//
// The Personal Allowance goes down by £1 for every £2 of
// income above the £100,000 limit. It can go down to zero.
// You do not get a Personal Allowance on taxable income over £125,000.
//
// https://www.gov.uk/government/publications/rates-and-allowances-income-tax/income-tax-rates-and-allowances-current-and-past#personal-allowances
func getPersonalAllowance(profitBeforeTaxes float64) float64 {
	if profitBeforeTaxes < 100000 {
		return personalAllowance
	}
	if profitBeforeTaxes > 125000 {
		return 0
	}
	return personalAllowance - (profitBeforeTaxes-100000)/2
}

// Class 	Rate for tax year 2020 to 2021
// -----    --------------------
// Class 2 	£3.05 a week
// Class 4 	9% on profits between £9,501 and £50,000
//          2% on profits over £50,000
func getNITax(profitBeforeTaxes float64) (float64, float64) {

	// THIS MUST BE CONFIGURABLE BY YEARS
	const yearlyPrimaryThreshold = 9501
	const yearlyUpperEarningsLimit = 50000
	const class2PerWeek = 3.05

	class2 := class2PerWeek * weeksInAYear

	var class4 float64
	if profitBeforeTaxes < yearlyPrimaryThreshold {
		class4 = 0.0
	} else if profitBeforeTaxes >= yearlyPrimaryThreshold && profitBeforeTaxes < yearlyUpperEarningsLimit {
		class4 = (profitBeforeTaxes - yearlyPrimaryThreshold) * 0.09
	} else {
		class4 = (yearlyUpperEarningsLimit-yearlyPrimaryThreshold)*0.09 +
			(profitBeforeTaxes-yearlyUpperEarningsLimit)*0.02
	}

	return math.Round(class2), math.Round(class4)
}

// returns current rate, how much before next threshold, and is it warning (when less than 20% left) or not
func HowMuchBeforeNextThreshold(personalIncome float64) (Rate, float64, bool) {
	const percentToWarning = 0.2
	left := 0.0
	isWarning := false
	if personalIncome < personalAllowance {
		left = personalAllowance - personalIncome
		isWarning = (left / personalIncome) <= percentToWarning
		return PersonalAllowance, left, isWarning
	}

	if personalIncome < 50000 {
		left = 50000 - personalIncome
		isWarning = (left / 50000) <= percentToWarning
		return BasicRate, left, isWarning
	}

	if personalIncome < 150000 {
		left = 150000 - personalIncome
		isWarning = (left / 150000) <= percentToWarning
		return HigherRate, left, isWarning
	}

	return AdditionalRate, 0, true
}
