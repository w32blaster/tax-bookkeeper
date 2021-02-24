package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/w32blaster/tax-bookkeeper/db"
)

const transactionsPerPage = 30

type PageableTransactions struct {
	app              *tview.Application
	table            *tview.Table
	currentPage      int
	totalPages       int
	loadTransactions FnLoadTransactions
	finderFocus      tview.Primitive // The primitive in the Finder that last had focus.
}

func BuildTxTable(a *tview.Application, totalTransactionsCnt int, fnLoadTransactions FnLoadTransactions) *PageableTransactions {
	return &PageableTransactions{
		table:            tview.NewTable().SetBorders(true),
		currentPage:      0,
		totalPages:       totalTransactionsCnt / transactionsPerPage,
		loadTransactions: fnLoadTransactions,
		app:              a,
	}
}

func (p *PageableTransactions) getPaginationText(currentPage int) string {
	return fmt.Sprintf("Page %d from %d", currentPage+1, p.totalPages+1)
}

func (p *PageableTransactions) Draw() *tview.Flex {

	paginationLabel := tview.NewTextView().SetText(p.getPaginationText(0))

	flexButtons := tview.NewFlex().SetDirection(tview.FlexColumn)
	flexButtons.SetTitleAlign(tview.AlignCenter)

	buttonBack := tview.NewButton("\u25C4").SetSelectedFunc(func() {
		go func() {
			p.app.QueueUpdateDraw(func() {
				if p.currentPage >= p.totalPages {
					return
				}

				p.currentPage = p.currentPage + 1
				tx := p.loadTransactions(30, p.currentPage)
				p.buildTransactionsListWidget(tx)
				paginationLabel.SetText(p.getPaginationText(p.currentPage))
			})
		}()
	})

	buttonForward := tview.NewButton("\u25BA").SetSelectedFunc(func() {
		go func() {
			p.app.QueueUpdateDraw(func() {
				if p.currentPage == 0 {
					return
				}

				p.currentPage = p.currentPage - 1
				tx := p.loadTransactions(30, p.currentPage)
				p.buildTransactionsListWidget(tx)
				paginationLabel.SetText(p.getPaginationText(p.currentPage))
			})
		}()
	})

	flexButtons.AddItem(buttonBack, 0, 1, true)
	flexButtons.AddItem(paginationLabel, 0, 5, false)
	flexButtons.AddItem(buttonForward, 0, 1, true)

	flexButtons.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft:
			p.app.SetFocus(buttonBack)
		case tcell.KeyRight:
			p.app.SetFocus(buttonForward)
		}
		return event
	})

	p.finderFocus = flexButtons

	// initial load of data
	tx := p.loadTransactions(transactionsPerPage, p.currentPage)
	p.buildTransactionsListWidget(tx)

	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.AddItem(flexButtons, 3, 0, true)
	flex.AddItem(p.table, 0, 1, false)

	p.app.SetFocus(flex)

	return flex
}

func (p *PageableTransactions) buildTransactionsListWidget(txs []db.Transaction) {

	p.table.Clear()

	if len(txs) == 0 {
		p.table.SetCell(0, 0,
			tview.NewTableCell("No data").
				SetTextColor(tcell.ColorRed).
				SetAlign(tview.AlignCenter))
		return
	}

	for r := 0; r < len(txs); r++ {

		// Cell 1, Date
		p.table.SetCell(r, 0,
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

		p.table.SetCell(r, 1,
			tview.NewTableCell(fmt.Sprintf("£%.02f", amount)).
				SetTextColor(color).
				SetAlign(tview.AlignLeft))

		p.table.SetCell(r, 2,
			tview.NewTableCell(txs[r].Description).
				SetTextColor(tcell.ColorWhite).
				SetAlign(tview.AlignLeft))
	}

}
