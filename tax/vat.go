package tax

import (
	"github.com/w32blaster/tax-bookkeeper/conf"
	"time"
)

// Quarterly VAT return dates are due for submission 1 month and 7 days after
// the of a VAT quarter. For example, a VAT return for the quarter-end
// June 2019 would be due by 7 August 2019.
// please refer to unit tests
func GetNextReturnDate(vatRegisteredMonth time.Month, now time.Time) (time.Month, time.Time) {
	closestMonth := getClosestSubmittingMonth(vatRegisteredMonth, now.Month())
	year := now.Year()
	if closestMonth < now.Month() {
		year = year + 1
	}
	return closestMonth, time.Date(year, closestMonth+2, 7, 0, 0, 0, 0, conf.GMT)
}

// check the unit test for examples
func getClosestSubmittingMonth(vat time.Month, current time.Month) time.Month {

	if vat == current {
		return current
	}

	if vat > current {
		vat = vat - 12
	}

	for i := vat; ; i = i + 3 {
		if i >= current {
			if i > 12 {
				return i - 12
			}
			return i
		}
	}
}

// safely subtracts 3 months, considering the year beginning
func GetBeginningOfPreviousPeriod(closestMonth time.Month, year int) time.Time {
	month := closestMonth - 2
	if month < 1 {
		return time.Date(year-1, 12-(-month), 1, 0, 0, 0, 0, conf.GMT)
	}

	return time.Date(year, month, 1, 0, 0, 0, 0, conf.GMT)
}
