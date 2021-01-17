package tax

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_CalculateSelfAssessmentTax(t *testing.T) {
	var tests = []struct {
		income      float64
		costs       float64
		expectedTax float64
	}{
		{90000, 0, 28104},
		{60000, 1000, 15084},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Expected %.0f tax from %.0f income", tt.expectedTax, tt.income),
			func(t *testing.T) {

				// When:
				selfAssessmentTax := CalculateSelfAssessmentTax(tt.income, tt.costs)

				// Then:
				assert.Equal(t, tt.expectedTax, selfAssessmentTax)
			},
		)
	}
}

// for testing I used these calculators:
// https://www.uktaxcalculators.co.uk/tax-calculators/personal-tax-calculators/self-employed-tax-calculator/#self-employed-income
// https://www.employedandselfemployed.co.uk/self-employed-tax-calculator
func Test_getPersonalTaxFrom(t *testing.T) {
	var tests = []struct {
		profitBeforeTaxes float64
		expectedTax       float64
	}{
		// income fits to personal allowance, no taxes are paid
		{10000, 0},

		// income fits to basic rate tax payer (under £50.000 per year)
		{15000, 500},
		{20000, 1500},
		{40000, 5500},

		// income fits to higher rate, but still below additional rate (under $150.000)
		{70000, 15500},
		{90000, 23500},
		{100000, 27500},

		// additional rate
		{120000, 39500},
		{170000, 61500},
		{200000, 75000},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Expected £%.0f tax from £%.0f income", tt.expectedTax, tt.profitBeforeTaxes),
			func(t *testing.T) {

				// When:
				tax := getPersonalTaxFrom(tt.profitBeforeTaxes)

				// Then:
				assert.Equal(t, tt.expectedTax, tax)
			},
		)
	}
}

func Test_getPersonalAllowanceForRich(t *testing.T) {

	var tests = []struct {
		profitBeforeTaxes float64
		expectedAllowance float64
	}{
		{50000, personalAllowance},

		{100000, personalAllowance},
		{100002, personalAllowance - 1},
		{100004, personalAllowance - 2},
		{100006, personalAllowance - 3},
		{100008, personalAllowance - 4},
		{100010, personalAllowance - 5},

		{100500, personalAllowance - 250},
		{101000, personalAllowance - 500},
		{110000, personalAllowance - 5000},

		{120000, 2500},
		{124996, 2},
		{124998, 1},

		// You do not get a Personal Allowance on taxable income over £125,000.
		{125000, 0},
		{130000, 0},
		{150000, 0},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Expected £%.0f allowance from £%.0f income", tt.expectedAllowance, tt.profitBeforeTaxes),
			func(t *testing.T) {

				// When:
				allowance := getPersonalAllowance(tt.profitBeforeTaxes)

				// Then:
				assert.Equal(t, tt.expectedAllowance, allowance)
			},
		)
	}
}
