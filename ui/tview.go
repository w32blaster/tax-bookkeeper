package ui

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/w32blaster/tax-bookkeeper/db"
)

type TerminalUI struct {
	app   *tview.Application
	pages *tview.Pages    // The application pages.
	focus tview.Primitive // The primitive in the Finder that last had focus.
}

func (t *TerminalUI) Start() {
	t.app = tview.NewApplication()
}

func (t *TerminalUI) DrawDashboard(data *DashboardData) {

	// last 10 transactions on the left
	infoFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	infoFlex.SetBorder(true).SetTitle(" Last transactions ").SetBorderPadding(1, 1, 1, 1)
	infoFlex.AddItem(buildTransactionsListWidget(data.LastTransactions), 0, 1, false)

	// Corporate tax
	prevCorpTax := buildCorporationTaxReportWidget(&data.PreviousPeriod, false)
	currCorpTax := buildCorporationTaxReportWidget(&data.CurrentPeriod, true)

	cpFlex := buildTwoColumnsWithDescription(" Corporate tax ", prevCorpTax, currCorpTax,
		"You must pay your Corporation "+
			"Tax 9 months and 1 day after the end  of your accounting "+
			"period https://www.gov.uk/pay-corporation-tax")

	// Self-assessment
	previousSA := buildSelfAssessmentTaxReportWidget(&data.PreviousSelfAssessmentPeriod, false)
	currentSA := buildSelfAssessmentTaxReportWidget(&data.CurrentSelfAssessmentPeriod, true)

	saFlex := buildTwoColumnsWithDescription(" Self-Assessment ", previousSA, currentSA,
		"Self assessment is between 6th of April and 5th April next year and the payment day is 31st of January")

	// VAT
	previousVatTable := buildVatReportWidget(&data.PreviousVAT, false)
	currentVatTable := buildVatReportWidget(&data.CurrentVAT, true)

	vatFlex := buildTwoColumnsWithDescription(" VAT ", previousVatTable, currentVatTable,
		"Quarterly VAT return dates are due for submission 1 month and 7 days after the of a VAT quarter")

	flex := tview.NewFlex().
		AddItem(infoFlex, 0, 1, false).
		AddItem(
			tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(cpFlex, 0, 2, false).
				AddItem(saFlex, 0, 2, false).
				AddItem(vatFlex, 0, 2, false).
				AddItem(tview.NewBox().SetBorder(true).SetTitle(" Loans "), 0, 2, false),
			0, 1, false)

	if err := t.app.SetRoot(flex, true).SetFocus(flex).Run(); err != nil {
		panic(err)
	}
}

func buildTransactionsListWidget(txs []db.Transaction) *tview.Table {

	table := tview.NewTable().SetBorders(true)

	for r := 0; r < len(txs); r++ {

		// Cell 1, Date
		table.SetCell(r, 0,
			tview.NewTableCell(txs[r].Date.Format("2 Jan 06")).
				SetTextColor(tcell.ColorWhite).
				SetAlign(tview.AlignLeft))

		// Cell 2, amount
		color := tcell.ColorWhite
		amount := txs[r].Debit
		if txs[r].Type == db.Credit {
			color = tcell.ColorGreen
			amount = txs[r].Credit
		}

		table.SetCell(r, 1,
			tview.NewTableCell(fmt.Sprintf("£%.02f", amount)).
				SetTextColor(color).
				SetAlign(tview.AlignLeft))

		table.SetCell(r, 2,
			tview.NewTableCell(txs[r].Description).
				SetTextColor(tcell.ColorWhite).
				SetAlign(tview.AlignLeft))
	}
	return table
}

func buildTwoColumnsWithDescription(title string, prevTable, currentTable *tview.Table, description string) *tview.Flex {

	// two columns
	tableCorporateTaxes := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(
		tview.NewFlex().
			AddItem(prevTable, 0, 1, false).
			AddItem(currentTable, 0, 1, false),
		0, 1, false)

	// wrapping final flex
	cpFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	cpFlex.SetBorder(true).SetTitle(title).SetBorderPadding(1, 1, 1, 1)
	cpFlex.
		AddItem(tableCorporateTaxes, 0, 1, false).
		AddItem(tview.NewTextView().SetText(description), 0, 1, false)

	return cpFlex
}

func buildCorporationTaxReportWidget(data *CorporateTax, isFuture bool) *tview.Table {

	color := "grey"
	if isFuture {
		color = "white"
	}

	cpLabel := "Corporate Tax: "
	if isFuture {
		cpLabel = "Corporate Tax so far (estimated): "
	}

	labels := [][]string{
		{"Tax for period: ", data.Period, color},
		{"Starting Date: ", data.StartingDate.Format("02 January 2006"), color},
		{"End Date: ", data.EndingDate.Format("02 January 2006"), color},
		{"Payment Date: ", data.NextPaymentDate.Format("02 January 2006"), "red"},
		{"Earned: ", "£" + floatToString(data.EarnedAccountingPeriod), color},
		{"Expenses: ", "£" + floatToString(data.ExpensesAccountingPeriod), color},
		{"Pension: ", "£" + floatToString(data.PensionAccountingPeriod), color},
		{cpLabel, "£" + floatToString(data.CorporateTaxSoFar), "green"},
	}

	cpHeader := "Previous Year Corporate tax"
	if isFuture {
		cpHeader = "Current year Corporate tax (not finished) "
	}

	table := tview.NewTable().SetBorders(false)

	var uLine tcell.Style
	uLine = uLine.Underline(true)

	// Cell 0, header
	table.SetCell(0, 0,
		tview.NewTableCell(cpHeader).
			SetStyle(uLine).
			SetTextColor(tcell.ColorWhite).
			SetAlign(tview.AlignLeft))

	for r := 0; r < len(labels); r++ {

		cellColor := tcell.GetColor(labels[r][2])

		// Cell 1, label
		table.SetCell(r+1, 0,
			tview.NewTableCell(labels[r][0]).
				SetTextColor(cellColor).
				SetAlign(tview.AlignLeft))

		// Cell 2, amount
		table.SetCell(r+1, 1,
			tview.NewTableCell(labels[r][1]).
				SetTextColor(cellColor).
				SetAlign(tview.AlignLeft))
	}

	return table
}

func buildSelfAssessmentTaxReportWidget(data *SelfAssessmentTax, isFuture bool) *tview.Table {

	color := "grey"
	if isFuture {
		color = "white"
	}

	cpLabel := "Self-Assessment tax: "
	if isFuture {
		cpLabel = "Self-Assessment tax so far: "
	}

	colorWarning := "grey"
	if data.IsWarning {
		colorWarning = "red"
	}

	labels := [][]string{
		{"Start dat: ", data.StartingDate.Format("02 January 2006"), color},
		{"End day: ", data.EndingDate.Format("02 January 2006"), color},
		{"Payment day: ", data.NextPaymentDate.Format("02 January 2006"), "red"},
		{"Moved out from company: ", "£" + floatToString(data.MovedOutFromCompanyTotal), color},
		{cpLabel, "£" + floatToString(data.SelfAssessmentTaxSoFar), "green"},
		{"Current tax rate: ", data.TaxRate.PrettyString(), color},
		{"Left before the following threshold: ", "£" + floatToString(data.HowMuchBeforeNextThreshold), colorWarning},
	}

	table := tview.NewTable().SetBorders(false)

	var uLine tcell.Style
	uLine = uLine.Underline(true)

	cpHeader := "Previous Year Self-Assessment tax"
	if isFuture {
		cpHeader = "Current year Self-Assessment (not finished) "
	}

	// Cell 0, header
	table.SetCell(0, 0,
		tview.NewTableCell(cpHeader).
			SetStyle(uLine).
			SetTextColor(tcell.ColorWhite).
			SetAlign(tview.AlignLeft))

	for r := 0; r < len(labels); r++ {

		cellColor := tcell.GetColor(labels[r][2])

		// Cell 1, label
		table.SetCell(r+1, 0,
			tview.NewTableCell(labels[r][0]).
				SetTextColor(cellColor).
				SetAlign(tview.AlignLeft))

		// Cell 2, amount
		table.SetCell(r+1, 1,
			tview.NewTableCell(labels[r][1]).
				SetTextColor(cellColor).
				SetAlign(tview.AlignLeft))
	}

	return table
}

func buildVatReportWidget(data *VAT, isFuture bool) *tview.Table {

	color := "grey"
	if isFuture {
		color = "white"
	}

	cpLabel := "VAT tax: "
	submitBy := "Submitted return by: "
	if isFuture {
		cpLabel = "VAT tax so far: "
		submitBy = "Submit your return by:"
	}

	labels := [][]string{
		{"VAT since: ", data.Since.Format("02 January 2006"), color},
		{"VAT until: ", data.Until.Format("02 January 2006"), color},
		{submitBy, data.NextMonthSubmit, color},
		{cpLabel, "£" + floatToString(data.NextVATToBePaidSoFar), color},
		{"Payment deadline: ", data.NextDateYouShouldPayFor.Format("02 January 2006"), "red"},
	}

	table := tview.NewTable().SetBorders(false)

	var uLine tcell.Style
	uLine = uLine.Underline(true)

	cpHeader := "Previous 3 months VAT tax"
	if isFuture {
		cpHeader = "Current VAT (not finished) "
	}

	// Cell 0, header
	table.SetCell(0, 0,
		tview.NewTableCell(cpHeader).
			SetStyle(uLine).
			SetTextColor(tcell.ColorWhite).
			SetAlign(tview.AlignLeft))

	for r := 0; r < len(labels); r++ {

		cellColor := tcell.GetColor(labels[r][2])

		// Cell 1, label
		table.SetCell(r+1, 0,
			tview.NewTableCell(labels[r][0]).
				SetTextColor(cellColor).
				SetAlign(tview.AlignLeft))

		// Cell 2, amount
		table.SetCell(r+1, 1,
			tview.NewTableCell(labels[r][1]).
				SetTextColor(cellColor).
				SetAlign(tview.AlignLeft))
	}

	return table
}

// here we attempt to guess and prefill category dropdown list by some words in description,
// add here as many "common" words so it could be easily to pre-fill dropdown list
func getInitialOptionByDescription(tx db.Transaction) int {
	descr := strings.ToLower(tx.Description)
	if tx.Type == db.Credit {
		if strings.Contains(descr, "loan") {
			return db.CreditTransactionUI.GetPositionFor(db.LoansReturn)
		}
		return db.CreditTransactionUI.GetPositionFor(db.Income)

	} else {

		if strings.Contains(descr, "dividend") {
			return db.DebitTransactionUI.GetPositionFor(db.Personal)
		} else if strings.Contains(descr, "salary") {
			return db.DebitTransactionUI.GetPositionFor(db.Personal)
		} else if strings.Contains(descr, "amznmktplace" /* amazon payment */) {
			return db.DebitTransactionUI.GetPositionFor(db.EquipmentExpenses)
		} else if strings.Contains(descr, "amazon" /* amazon payment */) {
			return db.DebitTransactionUI.GetPositionFor(db.EquipmentExpenses)
		} else if strings.Contains(descr, "energy") {
			return db.DebitTransactionUI.GetPositionFor(db.Premises)
		} else if strings.Contains(descr, "water") {
			return db.DebitTransactionUI.GetPositionFor(db.Premises)
		} else if strings.Contains(descr, "forx") {
			return db.DebitTransactionUI.GetPositionFor(db.BankCharges)
		} else if strings.Contains(descr, "fee") {
			return db.DebitTransactionUI.GetPositionFor(db.BankCharges)
		} else if strings.Contains(descr, "loan") {
			return db.DebitTransactionUI.GetPositionFor(db.Loan)
		} else if strings.Contains(descr, "hmrc") {
			return db.DebitTransactionUI.GetPositionFor(db.HMRC)
		} else if strings.Contains(descr, "pension") {
			return db.DebitTransactionUI.GetPositionFor(db.Pension)
		}

		return db.DebitTransactionUI.GetPositionFor(db.EquipmentExpenses)
	}
}

func (t *TerminalUI) BeginDialogToAllocateTransactions(unallocatedTxs []db.Transaction, fnAllocate FuncAllocateTransactions) {

	form := tview.NewForm()

	// populate dropdown list
	mapSelectedOptions := make(map[int]db.TransactionCategory)
	for idx, tx := range unallocatedTxs {
		if tx.Type == db.Credit {
			rowText := fmt.Sprintf("%d) %.2f (%s) - %s", idx, tx.Credit, tx.Date.Format("02 Jan 06"), tx.Description)
			form.AddDropDown(rowText, db.CreditTransactionUI.GetLabels(), getInitialOptionByDescription(tx),
				func(option string, optionIndex int) {
					mapSelectedOptions[tx.Pk] = db.TransactionCreditLabelMap[option]
				},
			)
		} else {
			rowText := fmt.Sprintf("%d) %.2f (%s) - %s", idx, tx.Debit, tx.Date.Format("02 Jan 06"), tx.Description)
			form.AddDropDown(rowText, db.DebitTransactionUI.GetLabels(), getInitialOptionByDescription(tx), func(option string, optionIndex int) {
				mapSelectedOptions[tx.Pk] = db.TransactionDebitLabelMap[option]
			})
		}
	}

	// button "Save" with callback
	form.AddButton(" Save ", func() {
		log.Println(mapSelectedOptions)
		if err := fnAllocate(mapSelectedOptions); err != nil {
			log.Println(err)
		} else {
			modal := tview.NewModal().
				SetText("All data were updated successfully").
				AddButtons([]string{" Ok "}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					if buttonLabel == " Ok " {
						t.app.Stop()
					}
				})
			t.app.SetRoot(modal, true)
		}
	})

	form.SetBorder(true).SetTitle("    Please allocate all " + strconv.Itoa(len(unallocatedTxs)) + " transactions    ").SetTitleAlign(tview.AlignLeft)
	if err := t.app.SetRoot(form, true).SetFocus(form).Run(); err != nil {
		panic(err)
	}
}

func floatToString(inputNum float64) string {
	return strconv.FormatFloat(inputNum, 'f', 2, 64)
}
