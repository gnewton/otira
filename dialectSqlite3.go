package otira

import (
	"log"
	"strconv"
)

type DialectSqlite3 struct {
}

func (d *DialectSqlite3) CreateTableString(t *TableMeta) string {
	s := "CREATE TABLE " + t.name + " ("

	for i, fm := range t.Fields() {
		if i != 0 {
			s += ", "
		}
		s += fm.Name() + " " + d.FieldType(fm) + d.Constraints(fm)
	}

	s += ")"
	return s
}

func (d *DialectSqlite3) Constraints(fm FieldMeta) string {
	//if fm.PrimaryKey() {
	//return " PRIMARY KEY"
	//}
	return ""
}

func (d *DialectSqlite3) FieldType(fm FieldMeta) string {
	switch v := fm.(type) {
	case *FieldMetaString:
		s := ""
		if fm.Fixed() {
			s += "CHAR"
		} else {
			s += "VARCHAR"
		}
		s += "(" + strconv.Itoa(fm.Length()) + ")"
		return s
	case *FieldMetaInt:
		return "INTEGER"
	case *FieldMetaFloat:
		return "REAL"
	case *FieldMetaByte:
		return "BLOB"
	default:
		log.Println("Unknown type", v)
		return "ERROR-UNKNOWN"
	}

}

func (d *DialectSqlite3) PreparedValueFormat(counter int) string {
	return "?"
}
