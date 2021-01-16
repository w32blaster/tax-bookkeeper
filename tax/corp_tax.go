package tax

import (
	"github.com/w32blaster/tax-bookkeeper/conf"
	"math"
	"strconv"
	"time"
)

const (
	// 1st of April
	financialYearStartDay   = 1
	financialYearStartMonth = time.April
)

// CalculateCorporateTax calculate corporate tax for a given financial year
//     yearProfit - company profit in £ (profit = revenue - expenses - pension - salary)
//     accountingPeriodEnd - end of accounting period specific for your company. It may be different from
//                           financial year (1 April - 31 March).
//                           in GMT timezone
//                           See more: https://www.gov.uk/corporation-tax-accounting-period
func CalculateCorporateTax(yearProfit float64, accountingPeriodStartDate time.Time) float64 {

	// simply multiply profit by rate
	if isMatchingFinYear(accountingPeriodStartDate) {
		finYear := getFinYear(accountingPeriodStartDate)
		rate := conf.CorporationTaxRates[finYear]
		return yearProfit * rate
	}

	// if both periods has the same rate, then calculate as in previous step
	prevPeriod, nextPeriod := getTwoPeriods(accountingPeriodStartDate)
	ratePrev := conf.CorporationTaxRates[prevPeriod]
	rateNext := conf.CorporationTaxRates[nextPeriod]
	if ratePrev == rateNext {
		return yearProfit * ratePrev
	}

	// otherwise, necessary tax will be calculated proportionally against
	// the government's tax year period date
	daysOne, daysTwo := getDaysForPeriods(accountingPeriodStartDate)
	return calculateTwoPeriodsDifferentRate(daysOne, ratePrev, daysTwo, rateNext, yearProfit)
}

// split accounting period by two slices. Depending on if the start date before of after 1st of April,
// it can return different periods. For examples please refer to unit test.
//
//   Before 1st of April:
//
//                 ┌------------------------┐
//                 |   accounting period    |
//   +---------|---+-----+-------------|----+---+----->
//            2019       ↑            2020      ↑
//                   1st of April             1st of April
//
//
//   After 1st of April:
//                          ┌----------------------┐
//                          |   accounting period  |
//   +---------|---------+--+----------|--------+--+---->
//            2019       ↑            2020      ↑
//                    1st of April            1st of April
//
func getTwoPeriods(accPeriodStartDate time.Time) (string, string) {

	year := accPeriodStartDate.Year()
	financialYearStartInThisYear := time.Date(year, financialYearStartMonth, financialYearStartDay, 0, 0, 0, 0, conf.GMT)

	var prevPeriod, nextPeriod string
	if accPeriodStartDate.After(financialYearStartInThisYear) {
		prevPeriod = strconv.Itoa(year) + "-" + strconv.Itoa(year+1)
		nextPeriod = strconv.Itoa(year+1) + "-" + strconv.Itoa(year+2)
	} else {
		prevPeriod = strconv.Itoa(year-1) + "-" + strconv.Itoa(year)
		nextPeriod = strconv.Itoa(year) + "-" + strconv.Itoa(year+1)
	}

	return prevPeriod, nextPeriod
}

// check whether accounting period for current company matches financial year start (1st of April)
func isMatchingFinYear(accPeriodStartDate time.Time) bool {
	return accPeriodStartDate.Day() == financialYearStartDay &&
		accPeriodStartDate.Month() == financialYearStartMonth
}

// returns year period for the giving accounting period
func getFinYear(accPeriodStartDate time.Time) string {
	year := accPeriodStartDate.Year()
	return strconv.Itoa(year) + "-" + strconv.Itoa(year+1)
}

// if accounting period doesn't match financial year, we need to find how much
// in days takes each period from different tax years.
//
// For example, if your accounting period is 1 January 2017 to 31 December 2017 we have:
//
//    1) financial year starting 1 April 2016 for 90 days (1 January 2017 to 31 March 2017)
//    2) financial year starting 1 April 2017 for 275 days (1 April 2017 to 31 December 2017)
//
func getDaysForPeriods(accPeriodStartDate time.Time) (int, int) {
	year := accPeriodStartDate.Year()
	financialYearStartInThisYear := time.Date(year, financialYearStartMonth, financialYearStartDay, 0, 0, 0, 0, conf.GMT)

	var daysPrev, daysNext int
	if accPeriodStartDate.After(financialYearStartInThisYear) {
		finYearEnd := financialYearStartInThisYear.AddDate(1, 0, 0)
		daysPrev = int(finYearEnd.Sub(accPeriodStartDate).Hours() / 24)

		accPeriodEnd := accPeriodStartDate.AddDate(1, 0, 0)
		daysNext = int(accPeriodEnd.Sub(financialYearStartInThisYear.AddDate(1, 0, 0)).Hours() / 24)

	} else {
		daysPrev = int(financialYearStartInThisYear.Sub(accPeriodStartDate).Hours() / 24)
		daysNext = int(accPeriodStartDate.AddDate(1, 0, 0).Sub(financialYearStartInThisYear).Hours() / 24)
	}

	return daysPrev, daysNext
}

// if our accounting period doesn't match financial year, then it is divided by 1st of April by
// two periods. And if these periods have different Corporate Tax Rate, we should calculate it
// proportionally against the government's tax year period date.
// Please refer to unit test for examples
func calculateTwoPeriodsDifferentRate(daysOne int, rateOne float64, daysTwo int, rateTwo float64, profit float64) float64 {
	daysInYear := float64(daysOne + daysTwo) // should be 365 or 366

	tax := (profit*(float64(daysOne)/daysInYear))*rateOne + /* period before 1st of April */
		(profit*(float64(daysTwo)/daysInYear))*rateTwo /* period after 1st of April */

	return math.Round(tax*100) / 100 // round for two decimal places
}
