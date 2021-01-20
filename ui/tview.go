package ui

import (
	"fmt"
	"github.com/rivo/tview"
	"github.com/w32blaster/tax-bookkeeper/db"
	"log"
)

type TerminalUI struct {
	DB    *db.Database
	app   *tview.Application
	pages *tview.Pages    // The application pages.
	focus tview.Primitive // The primitive in the Finder that last had focus.
}

func (t *TerminalUI) Start() {
	t.app = tview.NewApplication()
}

func (t *TerminalUI) BeginDialogToAllocateTransactions(unallocatedTxs []db.Transaction) {

	var labels []string
	for k, _ := range db.TransactionLabelMap {
		labels = append(labels, k)
	}

	form := tview.NewForm()

	var currentPrice float64
	mapSelectedOptions := make(map[int]db.TransactionCategory)
	for idx, tx := range unallocatedTxs {
		if tx.Debit != 0 {
			currentPrice = tx.Debit
		} else {
			currentPrice = tx.Credit
		}
		rowText := fmt.Sprintf("%d) %.2f (%s) - \n %s", idx, currentPrice, tx.Date.Format("02 Jan 06"), tx.Description)
		form.AddDropDown(rowText, labels, 0, func(option string, optionIndex int) {
			mapSelectedOptions[tx.Pk] = db.TransactionLabelMap[option]
		})
	}

	form.AddButton("Save", func() {
		log.Println(mapSelectedOptions)
		if err := t.DB.AllocateTransactions(mapSelectedOptions); err != nil {
			log.Println(err)
		} else {
			modal := tview.NewModal().
				SetText("All data updated").
				AddButtons([]string{"Ok"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					if buttonLabel == "Ok" {
						t.DB.Close()
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
