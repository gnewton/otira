package otira

import (
	"log"
)

type DialectMysql struct {
}

func (d *DialectMysql) CreateTableString(t *TableDef) string {
	log.Println("TODO")
	return "TODO"
}

func (d *DialectMysql) PreparedValueFormat(counter int) string {
	return "?"
}
