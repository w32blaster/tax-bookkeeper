package main

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/w32blaster/tax-bookkeeper/conf"
)

func Test_calculateCorporateTax(t *testing.T) {
	var tests = []struct {
		name                  string
		profit                float64
		accountingPeriodStart time.Time
		expectedCorpTax       float64
	}{
		// the accounting period matches the financial year (1st of April), so
		// tax is only 19% of the profit 40.000 x 19% = 7.600
		{"accounting date matches financial year 19%",
			40000.00, dateOf("01-04-2019"), 7600.00},
		{"accounting date matches financial year 20%",
			60000.00, dateOf("01-04-2016"), 12000.00},

		// the accounting period doesn't match, both both periods have the same rate
		{"accounting date matches financial year 19%",
			40000.00, dateOf("30-12-2019"), 7600.00},
		{"accounting date matches financial year 20%",
			60000.00, dateOf("01-09-2015"), 12000.00},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// When:
			corpTax := CalculateCorporateTax(tt.profit, tt.accountingPeriodStart)

			// Then:
			assert.Equal(t, tt.expectedCorpTax, corpTax)
		})
	}
}

func Test_getFinYear(t *testing.T) {
	var tests = []struct {
		accountingPeriodStart time.Time
		expectedPeriod        string
	}{

		{dateOf("01-04-2020"), "2020-2021"},
		{dateOf("20-11-2019"), "2019-2020"},
		{dateOf("10-01-2018"), "2018-2019"},
		{dateOf("05-07-2017"), "2017-2018"},
	}

	for _, tt := range tests {
		t.Run(tt.expectedPeriod, func(t *testing.T) {

			// When:
			period := getFinYear(tt.accountingPeriodStart)

			// Then:
			assert.Equal(t, tt.expectedPeriod, period)
		})
	}
}

// accounting period (1st of January 2017) starts before the financial year (1st of April),
// so it is divided by two periods:
//   1) 2016-2017 (which is 1 January 2017 to 31 March 2017)
//   2) 2017-2018 (which is 1 April 2017 to 31 December 2017)
func Test_getTwoPeriodsBefore(t *testing.T) {

	// When:
	prevPeriod, nextPeriod := getTwoPeriods(dateOf("01-01-2017"))

	// Then:
	assert.Equal(t, "2016-2017", prevPeriod)
	assert.Equal(t, "2017-2018", nextPeriod)
}

// accounting period (1st of August 2018) starts after the financial year (1st of April) begins,
// so it is divided by two periods:
//   1) 2018-2019 (which is 1 August 2018 to 1 April 2019)
//   2) 2019-2020 (which is 1 April 2019 to 31 March 2020)
func Test_getTwoPeriodsAfter(t *testing.T) {

	// When:
	prevPeriod, nextPeriod := getTwoPeriods(dateOf("01-08-2018"))

	// Then:
	assert.Equal(t, "2018-2019", prevPeriod)
	assert.Equal(t, "2019-2020", nextPeriod)
}

func Test_getDaysForPeriods(t *testing.T) {
	var tests = []struct {
		accountingPeriodStart time.Time
		daysPrev              int
		daysNext              int
	}{

		// after 1st April
		{dateOf("20-11-2019"), 132, 234},

		// before 1st April
		{dateOf("01-01-2017"), 90, 275},
	}

	for _, tt := range tests {
		t.Run(tt.accountingPeriodStart.Format("02 Jan 06 "), func(t *testing.T) {

			// When:
			daysPrev, daysNext := getDaysForPeriods(tt.accountingPeriodStart)

			// Then:
			assert.Equal(t, tt.daysPrev, daysPrev)
			assert.Equal(t, tt.daysNext, daysNext)
		})
	}
}

// shorthand for the date creation, like "01-03-2021"
func dateOf(date string) time.Time {
	parts := strings.Split(date, "-")
	year, _ := strconv.Atoi(parts[2])
	month, _ := strconv.Atoi(parts[1])
	day, _ := strconv.Atoi(parts[0])
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, conf.GMT)
}
