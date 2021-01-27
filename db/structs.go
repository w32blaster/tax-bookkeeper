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
	LoansReturn
	Loan // when you borrow some money from your company and have to return it back
)

var TransactionDebitLabelMap = map[string]TransactionCategory{
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
	"Load":                                                       Loan,
}

var TransactionCreditLabelMap = map[string]TransactionCategory{
	"income":      Income,
	"loan Return": LoansReturn,
}

type TransactionUi struct {
	labels                  []string
	categoryToLabelPosition map[TransactionCategory]int
}

func (d TransactionUi) GetLabels() []string {
	return d.labels
}

func (d TransactionUi) GetPositionFor(category TransactionCategory) int {
	return d.categoryToLabelPosition[category]
}

var DebitTransactionUI TransactionUi
var CreditTransactionUI TransactionUi

func init() {

	// Debit labels
	DebitTransactionUI = createUiElementFrom(TransactionDebitLabelMap)

	// Credit labels
	CreditTransactionUI = createUiElementFrom(TransactionCreditLabelMap)
}

// creates one element to me used for UI, to draw a Dropdown List
func createUiElementFrom(categories map[string]TransactionCategory) TransactionUi {

	// make array of labels to be displayed in dropdown list
	var labels = make([]string, len(categories))
	i := 0
	for k := range categories {
		labels[i] = k
		i++
	}
	sort.Strings(labels)

	// build map category <=> position in an array
	var positions = make(map[TransactionCategory]int, len(categories))
	for k, v := range categories {
		positions[v] = getPosition(k, labels)
	}

	return TransactionUi{
		labels:                  labels,
		categoryToLabelPosition: positions,
	}
}

func getPosition(label string, arr []string) int {
	for i, l := range arr {
		if l == label {
			return i
		}
	}
	return 0
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
