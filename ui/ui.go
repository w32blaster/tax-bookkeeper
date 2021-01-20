package ui

import "github.com/w32blaster/tax-bookkeeper/db"

// UI is a common interface for an GUI. At this moment we have only terminal UI,
// but if in the future we will need to do another UI, it would be easy possible
// to do by implementing this interface
type UI interface {
	Start()
	BeginDialogToAllocateTransactions(unallocatedTxs []db.Transaction)
}
