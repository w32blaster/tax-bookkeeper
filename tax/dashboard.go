package tax

import (
    "github.com/w32blaster/tax-bookkeeper/db"
    "github.com/w32blaster/tax-bookkeeper/ui"
    "log"
    "time"
)

// accountingDateStart is only day and month, like 01-11
func CollectDataForDashboard(d db.Database, accountingDateStart time.Time) ui.DashboardData {

    // get the profit for the current accounting period since accountingDateStart until now
    revenue, err := d.GetRevenueSince(accountingDateStart)
    if err != nil {
        log.Println(err)
    }

    expenses := d.GetExpensesSince(accountingDateStart)
    pension := d.GetPensionSince(accountingDateStart)

    // Corporate Tax
    corpTax := CalculateCorporateTax(0.0, accountingDateStart)

    return ui.DashboardData{
        CorporateTaxSoFar: corpTax,
    }
}
