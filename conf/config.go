package conf

import "time"

// historical Corporation Tax Rates
// https://www.gov.uk/corporation-tax-rates
var CorporationTaxRates = map[string]float64{
	"2015-2016": 0.2,
	"2016-2017": 0.2,
	"2017-2018": 0.19,
	"2019-2020": 0.19,
	"2020-2021": 0.19,
}

var GMT, _ = time.LoadLocation("GMT")

type App struct {
}
