package ui

import (
	"github.com/stretchr/testify/assert"
	"github.com/w32blaster/tax-bookkeeper/conf"
	"github.com/w32blaster/tax-bookkeeper/db"
	"strconv"
	"strings"
	"testing"
	"time"
)

// NB!
// Credit - Loan return
// Debit - loan take away

func TestActiveLoanIsPaid(t *testing.T) {

	// Given:
	tx := []db.Transaction{
		{Date: dateOf("01-01-2020"), Category: db.Loan, Debit: 100.0},
		{Date: dateOf("02-01-2020"), Category: db.LoansReturn, Credit: 50.0},
		{Date: dateOf("03-01-2020"), Category: db.LoansReturn, Credit: 50.0},
	}

	// When:
	left := getActiveLoan(tx)

	// Then:
	assert.Equal(t, 0.0, left)
}

func TestActiveLoanIsNotPaid(t *testing.T) {

	// Given:
	tx := []db.Transaction{
		{Date: dateOf("01-01-2020"), Category: db.Loan, Debit: 100.0},
		{Date: dateOf("02-01-2020"), Category: db.LoansReturn, Credit: 30.0},
		{Date: dateOf("03-01-2020"), Category: db.LoansReturn, Credit: 20.0},
	}

	// When:
	left := getActiveLoan(tx)

	// Then:
	assert.Equal(t, 50.0, left)
}

func TestActiveLoanIsNotPaidTwo(t *testing.T) {

	// Given:
	tx := []db.Transaction{
		{Date: dateOf("01-01-2020"), Category: db.Loan, Debit: 100.0},
		{Date: dateOf("02-01-2020"), Category: db.LoansReturn, Credit: 100.0},

		{Date: dateOf("01-02-2020"), Category: db.Loan, Debit: 100.0},
		{Date: dateOf("02-02-2020"), Category: db.LoansReturn, Credit: 30.0},
		{Date: dateOf("03-02-2020"), Category: db.LoansReturn, Credit: 30.0},
	}

	// When:
	left := getActiveLoan(tx)

	// Then:
	assert.Equal(t, 40.0, left)
}

func TestActiveLoanNoTransactions(t *testing.T) {

	// Given:
	var tx []db.Transaction

	// When:
	left := getActiveLoan(tx)

	// Then:
	assert.Equal(t, 0.0, left)
}

// shorthand for the date creation, like "01-03-2021"
func dateOf(date string) time.Time {
	parts := strings.Split(date, "-")
	year, _ := strconv.Atoi(parts[2])
	month, _ := strconv.Atoi(parts[1])
	day, _ := strconv.Atoi(parts[0])
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, conf.GMT)
}
