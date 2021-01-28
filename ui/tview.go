package ui

import (
	"fmt"
	"github.com/rivo/tview"
	"github.com/w32blaster/tax-bookkeeper/db"
	"log"
	"strconv"
	"strings"
)

type TerminalUI struct {
	app   *tview.Application
	pages *tview.Pages    // The application pages.
	focus tview.Primitive // The primitive in the Finder that last had focus.
}

func (t *TerminalUI) Start() {
	t.app = tview.NewApplication()
}

func (t *TerminalUI) DrawDashboard(data DashboardData) {

	flex := tview.NewFlex().
		AddItem(tview.NewBox().SetBorder(true).SetTitle(" Last transactions "), 0, 1, false).
		AddItem(
			tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(tview.NewBox().SetBorder(true).SetTitle(" Corporate tax "), 0, 1, false).
				AddItem(tview.NewBox().SetBorder(true).SetTitle(" Self-Assessment "), 0, 3, false).
				AddItem(tview.NewBox().SetBorder(true).SetTitle(" VAT "), 0, 1, false).
				AddItem(tview.NewBox().SetBorder(true).SetTitle(" Loans "), 0, 1, false),
			0, 2, false)

	if err := t.app.SetRoot(flex, true).SetFocus(flex).Run(); err != nil {
		panic(err)
	}

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
