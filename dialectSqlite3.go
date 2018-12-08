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
	s += d.ForeignKeys(t)
	s += ")"
	return s
}

func (d *DialectSqlite3) ForeignKeys(t *TableMeta) string {
	var s string
	s += d.oneToManyForeignKeys(t)

	return s
}

func (d *DialectSqlite3) oneToManyForeignKeys(t *TableMeta) string {
	var s string
	if t.oneToMany == nil {
		return ""
	}
	log.Println(t.GetName())
	for i := 0; i < len(t.oneToMany); i++ {
		one2m := t.oneToMany[i]
		log.Println(one2m)
		s += ", " + "FOREIGN KEY(" + one2m.leftKeyField.Name() + ") REFERENCES " + one2m.rightTable.GetName() + "(" + one2m.rightKeyField.Name() + ")"
	}
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

func (d *DialectSqlite3) InitPragmas() []string {
	var pragmas []string

	pragmas = append(pragmas, "PRAGMA foreign_keys = ON;")
	pragmas = append(pragmas, "PRAGMA schema.cache_size =-20000;")

	return pragmas
}
