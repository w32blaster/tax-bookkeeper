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
		{80000, 3000, 22644},
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

func Test_getNITax(t *testing.T) {

	var tests = []struct {
		profitBeforeTaxes float64
		expectedClass2Tax float64
		expectedClass4Tax float64
	}{
		{20000, 159.00, 945.00},
		{30000, 159.00, 1845.00},
		{40000, 159.00, 2745.00},
		{50000, 159.00, 3645.00},
		{60000, 159.00, 3845.00},
		{70000, 159.00, 4045.00},
		{80000, 159.00, 4245.00},
		{90000, 159.00, 4445.00},
		{100000, 159.00, 4645.00},
		{110000, 159.00, 4845.00},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Expected NI %.0f tax from %.0f income",
			tt.expectedClass2Tax+tt.expectedClass4Tax, tt.profitBeforeTaxes),
			func(t *testing.T) {

				// When:
				class2Tax, class4Tax := getNITax(tt.profitBeforeTaxes)

				// Then:
				assert.Equal(t, tt.expectedClass2Tax, class2Tax)
				assert.Equal(t, tt.expectedClass4Tax, class4Tax)
			},
		)
	}
}

func Test_HowMuchBeforeNextThreshold(t *testing.T) {

	var tests = []struct {
		income            float64
		expectedRate      Rate
		expectedMoneyLeft float64
		expectedIsWarning bool
	}{
		// income falls within personal allowance,  how much money left before £12.500
		{0, PersonalAllowance, 12500.00, false},
		{100, PersonalAllowance, 12400.00, false},
		{5000, PersonalAllowance, 7500.00, false},
		{12000, PersonalAllowance, 500.00, true},

		// income falls within the Basic Rate, how much money left before £50.000
		{20000, BasicRate, 30000.00, false},
		{30000, BasicRate, 20000.00, false},
		{45000, BasicRate, 5000.00, true},

		// income falls within the Higher Rate, how much money left before £150.000
		{70000, HigherRate, 80000.00, false},
		{120000, HigherRate, 30000.00, true},

		// income exceeds all limits and falls into Additional Rate
		{170000, AdditionalRate, 0.00 /* infinity */, true},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Expected left %.0f when income is %.02f with current rate %d ",
			tt.expectedMoneyLeft, tt.income, tt.expectedRate),
			func(t *testing.T) {

				// When:
				rate, leftBeforeNextThreshold, isWarning := HowMuchBeforeNextThreshold(tt.income)

				// Then:
				assert.Equal(t, tt.expectedRate, rate)
				assert.Equal(t, tt.expectedMoneyLeft, leftBeforeNextThreshold)
				assert.Equal(t, tt.expectedIsWarning, isWarning)
			},
		)
	}
}

func Test_GetTaxYearDatesNowIsAfterApril(t *testing.T) {

	// Given:
	now := dateOf("01-10-2019") // after the 6th of April, we are at the beginning of tax year

	// When:
	start, end, paymentDay := GetTaxYearDates(now)

	// Then:
	assert.Equal(t, dateOf("06-04-2019"), start)
	assert.Equal(t, dateOf("05-04-2020"), end)
	assert.Equal(t, dateOf("31-01-2021"), paymentDay)
}

func Test_GetTaxYearDatesNowIsBeforeApril(t *testing.T) {

	// Given:
	now := dateOf("01-01-2019") // before the 6th of April, we are at the end of tax year

	// When:
	start, end, paymentDay := GetTaxYearDates(now)

	// Then:
	assert.Equal(t, dateOf("06-04-2018"), start)
	assert.Equal(t, dateOf("05-04-2019"), end)
	assert.Equal(t, dateOf("31-01-2020"), paymentDay)
}

func Test_GetTaxYearDatesNowIsStart(t *testing.T) {

	// Given:
	now := dateOf("06-04-2019") // exactly 6th of April, the first day of tax year

	// When:
	start, end, paymentDay := GetTaxYearDates(now)

	// Then:
	assert.Equal(t, dateOf("06-04-2019"), start)
	assert.Equal(t, dateOf("05-04-2020"), end)
	assert.Equal(t, dateOf("31-01-2021"), paymentDay)
}
