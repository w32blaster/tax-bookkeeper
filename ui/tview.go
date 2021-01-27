package ui

import (
	"fmt"
	"github.com/rivo/tview"
	"github.com/w32blaster/tax-bookkeeper/db"
	"log"
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

func (t *TerminalUI) ShowDashboard() {

}

// here we attempt to guess and prefill category dropdown list by some words in description
func getInitialOptionByDescription(tx db.Transaction) int {
	descr := strings.ToLower(tx.Description)
	if tx.Type == db.Credit {
		if strings.Contains(descr, "loan") {
			db.CreditTransactionUI.GetPositionFor(db.LoansReturn)
		}
		return db.CreditTransactionUI.GetPositionFor(db.Income)

	} else {

		if strings.Contains(descr, "dividend") {
			db.DebitTransactionUI.GetPositionFor(db.Personal)
		} else if strings.Contains(descr, "salary") {
			db.DebitTransactionUI.GetPositionFor(db.Personal)
		} else if strings.Contains(descr, "AMZNMktplace" /* amazon payment */) {
			db.DebitTransactionUI.GetPositionFor(db.EquipmentExpenses)
		} else if strings.Contains(descr, "energy") {
			db.DebitTransactionUI.GetPositionFor(db.Premises)
		} else if strings.Contains(descr, "water") {
			db.DebitTransactionUI.GetPositionFor(db.Premises)
		} else if strings.Contains(descr, "forx") {
			db.DebitTransactionUI.GetPositionFor(db.BankCharges)
		} else if strings.Contains(descr, "loan") {
			db.DebitTransactionUI.GetPositionFor(db.Loan)
		}

		return db.DebitTransactionUI.GetPositionFor(db.BankCharges)
	}
}

func (t *TerminalUI) BeginDialogToAllocateTransactions(unallocatedTxs []db.Transaction, fnAllocate FuncAllocateTransactions) {

	form := tview.NewForm()

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

	form.SetBorder(true).SetTitle("    Please allocate all these transactions    ").SetTitleAlign(tview.AlignLeft)
	if err := t.app.SetRoot(form, true).SetFocus(form).Run(); err != nil {
		panic(err)
	}
}
