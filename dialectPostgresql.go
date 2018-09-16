package otira

import (
	"log"
	"strconv"
)

type DialectPostgresql struct {
}

func (d *DialectPostgresql) CreateTableString(t *TableMeta) string {
	log.Println("TODO")
	return "TODO"
}

func (d *DialectPostgresql) PreparedValueFormat(counter int) string {
	return "$" + strconv.Itoa(counter+1)
}
