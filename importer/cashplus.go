package importer

import (
	"encoding/csv"
	"fmt"
	"github.com/w32blaster/tax-bookkeeper/conf"
	"github.com/w32blaster/tax-bookkeeper/db"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// https://cashplus.com/
type CashPlus struct{}

const dateFormat = "02 January 2006"

func (c CashPlus) ReadAndParseFiles(path string) []db.Transaction {

	importPath, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Fatal("File " + path + " does not exist. Exit")
	}

	if importPath.IsDir() {

		// in a loop import all these files
		filePaths, err := listFilesInDir(path)
		if err != nil {
			log.Fatal("Can't list files inside directory, because: " + err.Error())
		}

		fmt.Printf("Found %d files\n\n", len(filePaths))
		var transactions = make([]db.Transaction, len(filePaths))
		for _, filePath := range filePaths {
			transactions = append(transactions, readAndImportSingleFile(filePath)...)
		}
		return transactions

	} else {
		// this is a file, just import one file
		return readAndImportSingleFile(path)
	}
}

func listFilesInDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}


func readAndImportSingleFile(filePath string) []db.Transaction {

	fmt.Println("  - Parsing file " + filePath)
	f, err := os.Open(filePath)
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

		transactions = append(transactions, db.Transaction{
			Date:          getDate(record[0]),
			Card:          record[1],
			Type:          txType,
			Description:   record[3],
			Credit:        getMoneySum(record[4]),
			Debit:         getMoneySum(record[5]),
			Balance:       getMoneySum(record[6]),
			ToBeAllocated: true,
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
