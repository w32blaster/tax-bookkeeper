package db

import (
	"sort"
	"time"
)

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
	Income
)

var TransactionLabelMap = map[string]TransactionCategory{
	"Unknown": Unknown,
	"Split between company and personal accounts":                Personal,
	"Legal (accountancy, advertising)":                           Legal,
	"Travel expenses":                                            Travel,
	"Office expenses":                                            Office,
	"Equipment expenses":                                         EquipmentExpenses,
	"Premises (heat, water, electricity)":                        Premises,
	"Cost of Sales (goods purchased for resale, subcontractors)": CostOfSales,
	"Wage payments, non-director salaries":                       WagesPayment,
	"Penalties and fines":                                        Penalties,
	"Bank charges":                                               BankCharges,
	"Pension":                                                    Pension,
	"HMRC (VAT payment, Corp Tax etc)":                           HMRC,
	"Fixed assets purchase":                                      FixedAssetPurchase,
}

var LabelsTransactionType []string

func init() {
	for k, _ := range TransactionLabelMap {
		LabelsTransactionType = append(LabelsTransactionType, k)
	}
	sort.Strings(LabelsTransactionType)
}

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
