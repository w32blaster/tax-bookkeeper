package tax

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetNextReturnDate(t *testing.T) {
	var tests = []struct {
		VATMonth                  time.Month
		now                       time.Time
		expectedSubmittingMonth   time.Month
		expectedVARReturnDeadline time.Time
	}{
		// if you were registered in November, then VAT Return dates are (when you submit return):
		// May, August, November and February
		// so that deadlines for payment are
		// 7 July,  7 Oct,  7 January and 7th of April
		{time.November, dateOf("15-12-2018"), time.February, dateOf("07-04-2019")},
		{time.November, dateOf("15-01-2019"), time.February, dateOf("07-04-2019")},
		{time.November, dateOf("15-02-2019"), time.February, dateOf("07-04-2019")},

		{time.November, dateOf("15-03-2019"), time.May, dateOf("07-07-2019")},
		{time.November, dateOf("15-04-2019"), time.May, dateOf("07-07-2019")},
		{time.November, dateOf("15-05-2019"), time.May, dateOf("07-07-2019")},

		{time.November, dateOf("15-06-2019"), time.August, dateOf("07-10-2019")},
		{time.November, dateOf("15-07-2019"), time.August, dateOf("07-10-2019")},
		{time.November, dateOf("15-08-2019"), time.August, dateOf("07-10-2019")},

		{time.November, dateOf("15-09-2019"), time.November, dateOf("07-01-2020")},
		{time.November, dateOf("15-10-2019"), time.November, dateOf("07-01-2020")},
		{time.November, dateOf("15-11-2019"), time.November, dateOf("07-01-2020")},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Expected submitting month %s and the VAT return deadline %s for the %s when VAT was registered in %s",
			tt.expectedSubmittingMonth.String(), tt.expectedVARReturnDeadline.Format("02 Jan 06 "), tt.now.Format("02 Jan 06 "), tt.VATMonth.String()),
			func(t *testing.T) {

				// When:
				submittingMonth, returnDeadline := GetNextReturnDate(tt.VATMonth, tt.now)

				// Then:
				assert.Equal(t, tt.expectedSubmittingMonth, submittingMonth)
				assert.Equal(t, tt.expectedVARReturnDeadline, returnDeadline)
			},
		)
	}
}

func TestFindClosestSubmittingMonth(t *testing.T) {
	var tests = []struct {
		VATMonth time.Month
		now      time.Month
		expected time.Month
	}{
		// if you were registered in November, then VAT Return dates are (when you submit return):
		// May, August, November and February
		{time.November, time.December, time.February},
		{time.November, time.January, time.February},
		{time.November, time.February, time.February},

		{time.November, time.March, time.May},
		{time.November, time.April, time.May},
		{time.November, time.May, time.May},

		{time.November, time.June, time.August},
		{time.November, time.July, time.August},
		{time.November, time.August, time.August},

		{time.November, time.September, time.November},
		{time.November, time.October, time.November},
		{time.November, time.November, time.November},

		// if VAT was registered in January, then VAT Return dates are:
		// Jan, Apr, Jul and Oct
		{time.January, time.February, time.April},
		{time.January, time.March, time.April},
		{time.January, time.April, time.April},

		{time.January, time.May, time.July},
		{time.January, time.June, time.July},
		{time.January, time.July, time.July},

		{time.January, time.August, time.October},
		{time.January, time.September, time.October},
		{time.January, time.October, time.October},

		{time.January, time.November, time.January},
		{time.January, time.December, time.January},
		{time.January, time.January, time.January},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Expected submitting month %s and the VAT submittion deadline %s for the %s",
			tt.expected.String(), tt.VATMonth.String(), tt.now.String()),
			func(t *testing.T) {

				// When:
				closestSubmittingMonth := getClosestSubmittingMonth(tt.VATMonth, tt.now)

				// Then:
				assert.Equal(t, tt.expected, closestSubmittingMonth)
			},
		)
	}
}

func TestGetBeginningOfPreviousPeriod(t *testing.T) {
	var tests = []struct {
		month    time.Month
		year     int
		expected time.Time
	}{
		// returns previous year!
		{time.January, 2020, dateOf("01-11-2019")},
		{time.February, 2020, dateOf("01-12-2019")},

		// the same year
		{time.March, 2020, dateOf("01-01-2020")},
		{time.April, 2020, dateOf("01-02-2020")},
		{time.May, 2020, dateOf("01-03-2020")},
		{time.June, 2020, dateOf("01-04-2020")},
		{time.July, 2020, dateOf("01-05-2020")},
		{time.August, 2020, dateOf("01-06-2020")},
		{time.September, 2020, dateOf("01-07-2020")},
		{time.October, 2020, dateOf("01-08-2020")},
		{time.November, 2020, dateOf("01-09-2020")},
		{time.December, 2020, dateOf("01-10-2020")},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Expected getting %s when the end of period is  %s",
			tt.expected.String(), tt.month.String()),
			func(t *testing.T) {

				// When:
				beginning := GetBeginningOfPreviousPeriod(tt.month, tt.year)

				// Then:
				assert.Equal(t, tt.expected.Format("02 Jan 06 "), beginning.Format("02 Jan 06 "))
			},
		)
	}
}
