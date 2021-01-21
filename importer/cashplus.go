package importer

import (
	"encoding/csv"
	"github.com/w32blaster/tax-bookkeeper/conf"
	"github.com/w32blaster/tax-bookkeeper/db"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// https://cashplus.com/
type CashPlus struct{}

const dateFormat = "02-Jan-06"

func (c CashPlus) ReadAndParseFile(path string) []db.Transaction {

	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	var transactions []db.Transaction
	firstLine := true
	for {

		record, err := r.Read()
		if err == io.EOF {
			break
		}

		if firstLine {
			firstLine = false
			continue
		}

		txType := getType(record[2])
		category := db.Unknown
		toBeAllocated := true

		// income transaction we don't need allocate, we know it is always "income"
		if txType == db.Credit {
			category = db.Income
			toBeAllocated = false
		}

		transactions = append(transactions, db.Transaction{
			Date:          getDate(record[0]),
			Card:          record[1],
			Type:          txType,
			Description:   record[3],
			Credit:        getMoneySum(record[4]),
			Debit:         getMoneySum(record[5]),
			Balance:       getMoneySum(record[6]),
			ToBeAllocated: toBeAllocated,
			Category:      category,
		})
	}

	return transactions
}

func getType(strType string) db.TransactionType {
	switch strType {
	case "Debit":
		return db.Debit
	case "Credit":
		return db.Credit
	}
	return 0
}

func getDate(strDate string) time.Time {
	parsedDate, err := time.ParseInLocation(dateFormat, strDate, conf.GMT)
	if err != nil {
		log.Fatalf("Can't parse date '%s' because of error: %s", strDate, err.Error())
	}

	return parsedDate
}

func getMoneySum(strNumber string) float64 {

	strNumber = strings.ReplaceAll(strNumber, "£", "")
	strNumber = strings.ReplaceAll(strNumber, ",", "")
	strNumber = strings.ReplaceAll(strNumber, "\"", "")
	strNumber = strings.ReplaceAll(strNumber, "(", "-")
	strNumber = strings.ReplaceAll(strNumber, ")", "")
	strNumber = strings.ReplaceAll(strNumber, " ", "")

	s, err := strconv.ParseFloat(strNumber, 64)
	if err != nil {
		log.Fatalf("Can't parse number %s because of error: %s", strNumber, err.Error())
	}

	return s
}
