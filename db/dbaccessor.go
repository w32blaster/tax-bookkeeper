package db

import (
	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/codec/msgpack"
	"github.com/asdine/storm/v3/q"
	"go.etcd.io/bbolt"
	"log"
	"math"
	"time"
)

type Database struct {
	db *storm.DB
}

func Init(dbPathFile string) *Database {

	// Open Storm DB
	boltdb, err := storm.Open(dbPathFile, storm.Codec(msgpack.Codec), storm.BoltOptions(0600, &bbolt.Options{Timeout: 5 * time.Second}))
	if err != nil {
		panic(err)
	}

	boltdb.Init(&Transaction{})

	return &Database{
		db: boltdb,
	}
}

func (d Database) Close() {
	d.db.Close()
}

func (d Database) GetAll(skipTo int) ([]Transaction, error) {
	var transactions []Transaction
	var err error
	if skipTo == 0 {
		err = d.db.All(&transactions)
	} else {
		err = d.db.All(&transactions, storm.Limit(skipTo))
	}
	return transactions, err
}

func (d Database) AllocateTransactions(cats map[int]TransactionCategory) error {
	tx, err := d.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for pk, cat := range cats {
		if err := tx.UpdateField(&Transaction{Pk: pk}, "ToBeAllocated", false); err != nil {
			return err
		}
		if err := tx.UpdateField(&Transaction{Pk: pk}, "Category", cat); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (d Database) GetUnallocated() ([]Transaction, error) {
	var transactions []Transaction
	err := d.db.Find("ToBeAllocated", true, &transactions)
	return transactions, err
}

func (d Database) ImportTransactions(transactions []Transaction) (int, error) {

	tx, err := d.db.Begin(true)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	for _, v := range transactions {

		if len(v.Description) == 0 && v.Balance == 0.0 {
			continue
		}

		if err := tx.Save(&v); err != nil {
			log.Fatalf("Can't save transaction, because %s", err.Error())
		}
	}

	return len(transactions), tx.Commit()
}

func (d Database) GetRevenueSince(accountingDateStart time.Time, accountingDateEnd time.Time) (float64, error) {

	var transactions []Transaction
	query := d.db.Select(
		q.And(
			q.Gt("Date", accountingDateStart),
			q.Lt("Date", accountingDateEnd),
			q.Eq("Type", Credit),
			q.Eq("ToBeAllocated", false),
			q.Eq("Category", Income),
		),
	)
	if err := query.Find(&transactions); err != nil {
		if err == storm.ErrNotFound {
			return 0, nil
		}
		return 0, err
	}

	var revenue float64
	for _, idx := range transactions {
		revenue = revenue + idx.Credit
	}
	return revenue, nil
}

func (d Database) GetExpensesSince(accountingDateStart time.Time, accountingDateEnd time.Time) (float64, error) {
	return _calculateExpensesByType(d.db, accountingDateStart, accountingDateEnd, Legal, Travel, Office, EquipmentExpenses, Premises, FixedAssetPurchase)
}

func (d Database) GetPensionSince(accountingDateStart time.Time, accountingDateEnd time.Time) (float64, error) {
	return _calculateExpensesByType(d.db, accountingDateStart, accountingDateEnd, Pension)
}

func (d Database) GetMovedOut(since time.Time, until time.Time) (float64, error) {
	return _calculateExpensesByType(d.db, since, until, Personal)
}

func _calculateExpensesByType(db *storm.DB, since time.Time, until time.Time, categories ...TransactionCategory) (float64, error) {

	// prepare the query
	var catMatcher q.Matcher
	if len(categories) == 1 {
		catMatcher = q.Eq("Category", categories[0])
	} else {
		var orMatcher = make([]q.Matcher, len(categories))
		for i, cat := range categories {
			orMatcher[i] = q.Eq("Category", cat)
		}
		catMatcher = q.Or(orMatcher...)
	}

	query := db.Select(
		q.And(
			q.Gt("Date", since),
			q.Lt("Date", until),
			q.Eq("Type", Debit),
			q.Eq("ToBeAllocated", false),
			catMatcher,
		),
	)

	var transactions []Transaction
	if err := query.Find(&transactions); err != nil {
		if err == storm.ErrNotFound {
			return 0, nil
		}
		return 0, err
	}

	var total float64
	for _, idx := range transactions {
		total = total + idx.Debit
	}
	return math.Abs(total), nil

}
