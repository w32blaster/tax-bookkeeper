package db

import "time"

type TransactionType int

const (
	Credit TransactionType = 1 + iota
	Debit
)

type TransactionCategory int

const (
	Unknown  TransactionCategory = 1 + iota
	Personal                     // split between company and personal accounts (salary and dividends)
	Legal                        // accountancy, advertising
	Travel
	Office            // rent
	EquipmentExpenses // computers, hosting
	Premises          // heat, water, electricity
	CostOfSales       // goods purchased for resale, subcontractors
	WagesPayment      // non-director sales
	Penalties
	BankCharges
	Pension
	HMRC
	FixedAssetPurchase
)

type (
	Transaction struct {
		Pk            int             `storm:"id,increment"` // primary key with auto increment
		Date          time.Time       `storm:"index"`        // midnight, GMT
		Type          TransactionType `storm:"index"`
		Card          string          // last 4 digits
		Description   string
		Credit        float64
		Debit         float64
		Balance       float64
		ToBeAllocated bool                `storm:"index"` // when category of this transaction is specified, it is "allocated"
		Category      TransactionCategory `storm:"index"`
	}
)
