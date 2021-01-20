package main

import (
	"flag"
	"fmt"
	"github.com/w32blaster/tax-bookkeeper/ui"
	"log"
	"os"

	"github.com/w32blaster/tax-bookkeeper/db"
	"github.com/w32blaster/tax-bookkeeper/importer"
)

var importCashPlus string

func main() {
	flag.Parse()

	d := db.Init()
	defer d.Close()

	gui := ui.TerminalUI{}

	if importCashPlus != "" {
		importDataAndExit(importer.CashPlus{}, d, importCashPlus, &gui)
	}
}

func importDataAndExit(i importer.Importer, d *db.Database, filePath string, gui ui.UI) {

	transactions := i.ReadAndParseFile(filePath)
	inserted, err := d.ImportTransactions(transactions)
	if err != nil {
		log.Fatal("Transactions import failed. The reason is: " + err.Error())
	}

	gui.Start()
	gui.BeginDialogToAllocateTransactions(transactions)

	fmt.Printf("Successfully imported %d transactions. Exit", inserted)

	d.Close()
	os.Exit(0)
}

func init() {
	flag.StringVar(&importCashPlus, "import-cashplus", "", "import transactions in CSV format from Cashplus")
}
