package otira

import (
	"log"
	"strconv"
)

type DialectOracle struct {
}

func (d *DialectOracle) CreateTableString(t *TableDef) string {
	log.Println("TODO")
	return "TODO"
}

func (d *DialectOracle) PreparedValueFormat(counter int) string {
	return ":val" + strconv.Itoa(counter+1)
}
