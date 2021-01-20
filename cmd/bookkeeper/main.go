package main

import (
	"flag"
	"fmt"
	"github.com/w32blaster/tax-bookkeeper/db"
	"github.com/w32blaster/tax-bookkeeper/importer"
	"github.com/w32blaster/tax-bookkeeper/ui"
	"log"
	"os"
)

var importCashPlus string

func main() {
	flag.Parse()

	d := db.Init()
	defer d.Close()

	gui := ui.TerminalUI{DB: d}

	if importCashPlus != "" {
		importDataAndExit(importer.CashPlus{}, d, importCashPlus)
	}

	if unallocatedTransactions, err := d.GetUnallocated(); err != nil {
		log.Fatal(err)
	} else if len(unallocatedTransactions) > 0 {
		gui.Start()
		gui.BeginDialogToAllocateTransactions(unallocatedTransactions)
	}
}

func importDataAndExit(i importer.Importer, d *db.Database, filePath string) {

	transactions := i.ReadAndParseFile(filePath)
	inserted, err := d.ImportTransactions(transactions)
	if err != nil {
		log.Fatal("Transactions import failed. The reason is: " + err.Error())
	}

	fmt.Printf("Successfully imported %d transactions. Exit", inserted)
	os.Exit(0)
}

func init() {
	flag.StringVar(&importCashPlus, "import-cashplus", "", "import transactions in CSV format from Cashplus")
}
