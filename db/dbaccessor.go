package db

import (
	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/codec/msgpack"
	"github.com/asdine/storm/v3/q"
	"go.etcd.io/bbolt"
	"log"
	"time"
)

type Database struct {
	db *storm.DB
}

func Init() *Database {

	// Open Storm DB
	boltdb, err := storm.Open("./tax-bookkeeper.db", storm.Codec(msgpack.Codec), storm.BoltOptions(0600, &bbolt.Options{Timeout: 5 * time.Second}))
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

func (d Database) GetAll() ([]Transaction, error) {
	var transactions []Transaction
	err := d.db.All(&transactions)
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
		if err := tx.Save(&v); err != nil {
			log.Fatalf("Can't save transaction, because %s", err.Error())
		}
	}

	return len(transactions), tx.Commit()
}

func (d Database) GetRevenueSince(accountingDateStart time.Time) (float64, error) {

	var transactions []Transaction
	query := d.db.Select(
		q.And(
			q.Gt("Date", accountingDateStart),
			q.Eq("Type", Credit),
			q.Eq("ToBeAllocated", true),
			q.Eq("Category", Income),
		),
	)
	if err := query.Find(&transactions); err != nil {
		return 0, err
	}

	var revenue float64
	for _, idx := range transactions {
		revenue = revenue + idx.Credit
	}
	return revenue, nil
}

func (d Database) GetExpensesSince(accountingDateStart time.Time) (float64, error) {
	var transactions []Transaction
	query := d.db.Select(
		q.And(
			q.Gt("Date", accountingDateStart),
			q.Eq("Type", Debit),
			q.Eq("ToBeAllocated", true),
			q.Or(
				q.Eq("Category", Legal),
				q.Eq("Category", Travel),
				q.Eq("Category", Office),
				q.Eq("Category", EquipmentExpenses),
				q.Eq("Category", Premises),
				q.Eq("Category", FixedAssetPurchase),
			),
		),
	)
	if err := query.Find(&transactions); err != nil {
		return 0, err
	}

	var expenses float64
	for _, idx := range transactions {
		expenses = expenses + idx.Credit
	}
	return expenses, nil

}

func (d Database) GetPensionSince(accountingDateStart time.Time) (float64, error) {

	var transactions []Transaction
	query := d.db.Select(
		q.And(
			q.Gt("Date", accountingDateStart),
			q.Eq("Type", Debit),
			q.Eq("ToBeAllocated", true),
			q.Eq("Category", Pension),
		),
	)
	if err := query.Find(&transactions); err != nil {
		return 0, err
	}

	var pension float64
	for _, idx := range transactions {
		pension = pension + idx.Credit
	}
	return pension, nil
}
