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
	query := d.db.Select(
		q.And(
			q.Eq("ToBeAllocated", true),
			q.Eq("Type", Debit),
		),
	)
	err := query.Find(&transactions)
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
