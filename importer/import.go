package importer

import "github.com/w32blaster/tax-bookkeeper/db"

type Importer interface {
	ReadAndParseFiles(path string) []db.Transaction
}
