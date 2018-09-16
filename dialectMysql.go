package otira

import (
	"log"
)

type DialectMysql struct {
}

func (d *DialectMysql) CreateTableString(t *TableMeta) string {
	log.Println("TODO")
	return "TODO"
}

func (d *DialectMysql) PreparedValueFormat(counter int) string {
	return "?"
}
