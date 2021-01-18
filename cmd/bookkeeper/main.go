package main

import (
	"flag"
	"fmt"
	"github.com/w32blaster/tax-bookkeeper/db"
	"github.com/w32blaster/tax-bookkeeper/importer"
)

func main() {
	fmt.Println("dfd")

	importCashPlusPtr := flag.String("import-cashplus", "", "import transactions in CSV format from Cashplus")
	flag.Parse()

	d := db.Init()
	defer d.Close()

	var i importer.Importer
	if importCashPlusPtr != nil {
		i = importer.CashPlus{}
		transactions := i.ReadAndParseFile(*importCashPlusPtr)
		d.ImportTransactions(transactions)

		fmt.Println("Imported!")
	}
}
