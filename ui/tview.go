package ui

import (
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

func (t *TerminalUI) BeginDialogToAllocateTransactions(unallocatedTxs []db.Transaction) {

	list := tview.NewList()
	for idx, tx := range unallocatedTxs {
		list.AddItem(tx.Date.Format("02 Jan 06"), tx.Description, rune(idx), nil)
	}

	if err := t.app.SetRoot(list, true).SetFocus(list).Run(); err != nil {
		panic(err)
	}
}
