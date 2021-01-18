package importer

import "github.com/w32blaster/tax-bookkeeper/db"

type Importer interface {
	ReadAndParseFile(path string) []db.Transaction
}
