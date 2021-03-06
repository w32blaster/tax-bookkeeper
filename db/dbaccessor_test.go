package db

import (
	"github.com/stretchr/testify/assert"
	"github.com/w32blaster/tax-bookkeeper/conf"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

var now = time.Now().In(conf.GMT)

func TestCalculatePension(t *testing.T) {

	// create real DB
	const dbFile = "/tmp/tax-bookkeeper-pensions.db"
	db := Init(dbFile)
	defer func() {
		db.Close()
		os.Remove(dbFile)
	}()

	// Dates:
	tooLate := dateOf("01-10-2019")
	recently := dateOf("20-12-2019")

	// Populate with data:
	inserted, err := db.ImportTransactions([]Transaction{

		// these will be ignored, because they were too far away from current date
		_debitTransaction(Pension, 10.0, "Pension, ignored, too far away", tooLate),
		_debitTransaction(Pension, 30.0, "Pension, ignored, to far away", tooLate),

		// ignored because not pension
		_debitTransaction(Legal, 100.0, "should be ignored", tooLate),
		_debitTransaction(Premises, 100.0, "should be ignored", recently),

		// counting, because resent pensions
		_debitTransaction(Pension, 50.0, "Pension", recently),
		_debitTransaction(Pension, 60.0, "Pension", recently),
	})
	assert.Nil(t, err)
	assert.Equal(t, 6, inserted)

	// When:
	total, err := db.GetPensionSince(dateOf("01-12-2019"), now)

	// Then:
	assert.Nil(t, err)
	assert.Equal(t, 110.0, total)
}

func TestCalculateExpenses(t *testing.T) {

	// create real DB
	const dbFile = "/tmp/tax-bookkeeper-expenses.db"
	db := Init(dbFile)
	defer func() {
		db.Close()
		os.Remove(dbFile)
	}()

	// Dates:
	tooLate := dateOf("01-10-2019")
	recently := dateOf("20-12-2019")

	// Populate with data:
	inserted, err := db.ImportTransactions([]Transaction{

		// these will be ignored, because they were too far away from current date
		_debitTransaction(Legal, 10.0, "ignored, too far away", tooLate),
		_debitTransaction(Travel, 30.0, "ignored, to far away", tooLate),
		_debitTransaction(Office, 30.0, "ignored, to far away", tooLate),

		// ignored because not expenses
		_debitTransaction(Pension, 100.0, "should be ignored", recently),
		_debitTransaction(Personal, 200.0, "should be ignored", recently),

		// counting, because resent expenses
		_debitTransaction(Legal, 50.0, "Ok", recently),
		_debitTransaction(Travel, 60.0, "Ok", recently),
		_debitTransaction(Office, 30.0, "Ok", recently),
		_debitTransaction(EquipmentExpenses, 70.0, "Ok", recently),
		_debitTransaction(Premises, 50.0, "Ok", recently),
	})
	assert.Nil(t, err)
	assert.Equal(t, 10, inserted)

	// When:
	total, err := db.GetExpensesSince(dateOf("01-12-2019"), now)

	// Then:
	assert.Nil(t, err)
	assert.Equal(t, 260.0, total)
}

func TestCalculateExpensesRecently(t *testing.T) {

	// create real DB
	const dbFile = "/tmp/tax-bookkeeper-expenses2.db"
	db := Init(dbFile)
	defer func() {
		db.Close()
		os.Remove(dbFile)
	}()

	// Dates:
	tooLate := dateOf("01-10-2019")
	middle := dateOf("20-12-2019")
	tooEarly := dateOf("10-02-2020")

	// Populate with data:
	inserted, err := db.ImportTransactions([]Transaction{

		// these will be ignored, because they were too far away from current date
		_debitTransaction(Legal, 10.0, "ignored, too far away", tooLate),
		_debitTransaction(Travel, 30.0, "ignored, to far away", tooLate),
		_debitTransaction(Office, 30.0, "ignored, to far away", tooLate),

		// ignored because not expenses
		_debitTransaction(Legal, 100.0, "ok, between tooEarly and tooLate", middle),
		_debitTransaction(Travel, 200.0, "ok, between tooEarly and tooLate", middle),
		_debitTransaction(Pension, 400.0, "proper time, but its pension, ignored", middle),

		// counting, because resent expenses
		_debitTransaction(Legal, 50.0, "ignored, too recently", tooEarly),
		_debitTransaction(Travel, 60.0, "ignored, too recently", tooEarly),
		_debitTransaction(Pension, 30.0, "ignored, its pension", tooEarly),
		_debitTransaction(EquipmentExpenses, 70.0, "ignored, too recently", tooEarly),
		_debitTransaction(Premises, 50.0, "ignored, too recently", tooEarly),
	})
	assert.Nil(t, err)
	assert.Equal(t, 11, inserted)

	// When:
	total, err := db.GetExpensesSince(dateOf("01-11-2019"), dateOf("01-01-2020")) // between tooEarly and tooLate

	// Then:
	assert.Nil(t, err)
	assert.Equal(t, 100.0+200.0, total)
}

func TestCalculateExpensesNegativeNumbers(t *testing.T) {

	// create real DB
	const dbFile = "/tmp/tax-bookkeeper-expenses-neg.db"
	db := Init(dbFile)
	defer func() {
		db.Close()
		os.Remove(dbFile)
	}()

	// Dates:
	recently := dateOf("20-12-2019")

	// Populate with data:
	inserted, err := db.ImportTransactions([]Transaction{

		// counting, because resent expenses
		_debitTransaction(Legal, -50.0, "Ok", recently),
		_debitTransaction(Travel, -60.0, "Ok", recently),
		_debitTransaction(Office, -30.0, "Ok", recently),
		_debitTransaction(EquipmentExpenses, -70.0, "Ok", recently),
		_debitTransaction(Premises, -50.0, "Ok", recently),
	})
	assert.Nil(t, err)
	assert.Equal(t, 5, inserted)

	// When:
	total, err := db.GetExpensesSince(dateOf("01-12-2019"), now)

	// Then:
	assert.Nil(t, err)
	assert.Equal(t, 260.0, total) // still positive number
}

func _debitTransaction(cat TransactionCategory, debit float64, description string, txDate time.Time) Transaction {
	return Transaction{
		Date:          txDate,
		Card:          "0000",
		Type:          Debit,
		Description:   description,
		Credit:        0.0,
		Balance:       0.0,
		Debit:         debit,
		ToBeAllocated: false,
		Category:      cat,
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
