package ui

import (
	"fmt"
	"github.com/rivo/tview"
	"github.com/w32blaster/tax-bookkeeper/db"
	"log"
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

func (t *TerminalUI) BeginDialogToAllocateTransactions(unallocatedTxs []db.Transaction, fnAllocate FuncAllocateTransactions) {

	form := tview.NewForm()

	mapSelectedOptions := make(map[int]db.TransactionCategory)
	for idx, tx := range unallocatedTxs {
		if tx.Type == db.Credit {
			rowText := fmt.Sprintf("%d) %.2f (%s) - %s", idx, tx.Credit, tx.Date.Format("02 Jan 06"), tx.Description)
			form.AddDropDown(rowText, db.LabelsTransactionTypeCredit, 0, func(option string, optionIndex int) {
				mapSelectedOptions[tx.Pk] = db.TransactionCreditLabelMap[option]
			})
		} else {
			rowText := fmt.Sprintf("%d) %.2f (%s) - %s", idx, tx.Debit, tx.Date.Format("02 Jan 06"), tx.Description)
			form.AddDropDown(rowText, db.LabelsTransactionTypeDebit, 0, func(option string, optionIndex int) {
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
