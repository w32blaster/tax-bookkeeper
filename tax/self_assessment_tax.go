package tax

const (
	personalAllowance = 12500.0
)

// Tax Year is from 6 April to 5 April
func CalculateSelfAssessmentTax(income, costs float64) float64 {

	profitBeforeTaxes := income - costs

	personalTax := getPersonalTaxFrom(profitBeforeTaxes)

	return personalTax
}

// 	Band						Taxable income	 		Tax rate
//  -------------				--------------			---------
// 	Personal Allowance 			Up to £12,500	 		0%
// 	Basic rate			 		£12,501 to £50,000 		20%
// 	Higher rate 				£50,001 to £150,000 	40%
// 	Additional rate 			over £150,000 			45% (rich bastard!)
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
