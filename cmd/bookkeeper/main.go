package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/w32blaster/tax-bookkeeper/conf"
	"github.com/w32blaster/tax-bookkeeper/db"
	"github.com/w32blaster/tax-bookkeeper/importer"
	"github.com/w32blaster/tax-bookkeeper/tax"
	"github.com/w32blaster/tax-bookkeeper/ui"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var isHelp bool
var importCashPlus, accountingPeriodStartDate string
var r = regexp.MustCompile("^[0-9]{2}-[0-9]{2}$")

func main() {
	flag.Parse()
	if isHelp {
		fmt.Println("Tax Bookeeper. Helps you to analyze your taxes. Usage \n\n " +
			"-import-cashplus=/some/path - import transactions for CashPlus bank (file or directory) \n " +
			"-accounting-start=01-11 - set the accounting period date, if it doesn't match to financial year (1st of April)")
		os.Exit(0)
	}

	// TODO: validate date if set

	d := db.Init()
	defer d.Close()

	gui := ui.TerminalUI{}

	// import data and exit
	if importCashPlus != "" {
		importDataAndExit(importer.CashPlus{}, d, importCashPlus)
	}

	// if there are unallocated transactions, show the list
	if unallocatedTransactions, err := d.GetUnallocated(); err == nil && len(unallocatedTransactions) > 0 {
		gui.Start()
		gui.BeginDialogToAllocateTransactions(unallocatedTransactions, d.AllocateTransactions)
	}

	// or show the dashboards
	gui.Start()
	accPeriod, err := getNearestAccountingDate(accountingPeriodStartDate, time.Now().In(conf.GMT))
	if err != nil {
		log.Fatal(err)
	}

	dashboardData := tax.CollectDataForDashboard(d, accPeriod)
	gui.DrawDashboard(dashboardData)
}

func importDataAndExit(i importer.Importer, d *db.Database, filePath string) {

	transactions := i.ReadAndParseFiles(filePath)
	inserted, err := d.ImportTransactions(transactions)
	if err != nil {
		log.Fatal("Transactions import failed. The reason is: " + err.Error())
	}

	fmt.Printf("Successfully imported %d transactions. Exit\n", inserted)
	os.Exit(0)
}

func init() {
	flag.BoolVar(&isHelp, "h", false, "print help")
	flag.StringVar(&importCashPlus, "import-cashplus", "", "import transactions in CSV format from Cashplus")
	flag.StringVar(&accountingPeriodStartDate, "accounting-start", "", "If your Accounting Period start is different from financial year start,"+
		"you can set your date with this parameter, (example 01-11 which is 1st of November)")
}

// here we find the current account date.
// Please refer to unit tests
func getNearestAccountingDate(accountingDateStart string, now time.Time) (time.Time, error) {

	if accountingDateStart == "" {
		accountingDateStart = "01-04" // default accounting period matches the financial year, which is the 1st of April
	}

	if !r.MatchString(accountingDateStart) {
		return time.Time{}, errors.New("the incoming data is not valid! It should be like '01-11' (1st of November)")
	}

	parts := strings.Split(accountingDateStart, "-")
	day, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, err
	}
	if day < 1 || day > 31 {
		return time.Time{}, errors.New("day must be between 1 and 31")
	}

	m, err := strconv.Atoi(parts[1])
	if err != nil {
		return time.Time{}, err
	}
	if m < 1 || m > 12 {
		return time.Time{}, errors.New("month must be between 1 and 12")
	}
	month := time.Month(m)

	// today is the first day of a new accounting period year
	if day == now.Day() && month == now.Month() {
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, conf.GMT), nil
	}

	if now.Month() < month {
		return time.Date(now.Year()-1, month, day, 0, 0, 0, 0, conf.GMT), nil
	}

	return time.Date(now.Year(), month, day, 0, 0, 0, 0, conf.GMT), nil
}
