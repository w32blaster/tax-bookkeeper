package main

import (
    "github.com/stretchr/testify/assert"
    "github.com/w32blaster/tax-bookkeeper/conf"
    "strconv"
    "strings"
    "testing"
    "time"
)

func Test_getNearestAccountingDate(t *testing.T) {

    var tests = []struct {
        accountingPeriodStart string // like "01-11"
        now                   time.Time
        expectedAccPeriod     time.Time
        isErrorExpected       bool
    }{
        // invalid accounting period data
        {"0-0", time.Time{}, time.Time{}, true},
        {"incorrect_data", time.Time{}, time.Time{}, true},
        {"1-12", time.Time{}, time.Time{}, true},
        {"01.12", time.Time{}, time.Time{}, true},
        {"01/04", time.Time{}, time.Time{}, true},
        {"33-04", time.Time{}, time.Time{}, true},
        {"00-04", time.Time{}, time.Time{}, true},
        {"01-13", time.Time{}, time.Time{}, true},
        {"01-00", time.Time{}, time.Time{}, true},

        // matches: today is the first day of the next accounting period
        {"01-04", dateOf("01-04-2019"), dateOf("01-04-2019"), false},
        {"01-10", dateOf("01-10-2020"), dateOf("01-10-2020"), false},

        // current date (1st of April) is before the accounting period end in this year (1st of Oct),
        // so it was started in previous year (1st of October 2019)
        {"01-10", dateOf("01-04-2020"), dateOf("01-10-2019"), false},
        {"01-04", dateOf("23-02-2019"), dateOf("01-04-2018"), false},

        // current date (23rd of July 2019) is after accounting period start in this year (1st of April), so
        // it will end up in the next year, and start in this year (1st of April 2019)
        {"01-04", dateOf("23-07-2019"), dateOf("01-04-2019"), false},

        // today is the last day of current accounting period
        {"01-04", dateOf("31-03-2019"), dateOf("01-04-2018"), false},

    }

    for _, tt := range tests {
        t.Run(tt.accountingPeriodStart, func(t *testing.T) {

            // When:
            accPeriod, err := getNearestAccountingDate(tt.accountingPeriodStart, tt.now)

            // Then:
            assert.Equal(t, tt.expectedAccPeriod, accPeriod)
            assert.Equal(t, tt.isErrorExpected, err != nil)
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
