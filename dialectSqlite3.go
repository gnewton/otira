package otira

import (
	"errors"
	"log"
	"strconv"
)

var sqlite3DefaultPragmas = []string{
	"PRAGMA foreign_keys = ON;",
	"PRAGMA schema.cache_size =-20000;",
	"pragma synchronous = OFF;",
	"pragma journal_mode = OFF;",
	"pragma count_changes = OFF;",
	"pragma temp_store = MEMORY;",
}

type DialectSqlite3 struct {
	pragmas []string
}

func NewDialectSqlite3(pragmas []string, overwriteDefaultPragmas bool) Dialect {
	dialect := new(DialectSqlite3)
	if pragmas != nil {
		if overwriteDefaultPragmas {
			dialect.pragmas = pragmas
		} else {
			dialect.pragmas = append(pragmas, sqlite3DefaultPragmas...)
		}
	}
	return new(DialectSqlite3)
}

func (d *DialectSqlite3) DropTableIfExists(tm *TableMeta) string {
	return "DROP TABLE IF EXISTS " + tm.name
}

func (d *DialectSqlite3) CreateTableString(t *TableMeta) (string, error) {
	if t == nil {
		return "", errors.New("Tablemeta is nil")
	}
	s := CREATE_TABLE + SPC + t.name + " ("

	for i, fm := range t.Fields() {
		if i != 0 {
			s += ", "
		}
		s += fm.Name() + " " + d.FieldType(fm)
		if fm == t.primaryKey {
			s += SPC + PRIMARY_KEY
		}
		s += d.Constraints(fm)
	}
	s += d.ForeignKeys(t)
	s += ")"
	return s, nil
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
	if fm.Unique() {
		return " UNIQUE"
	}
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
	case *FieldMetaInt, *FieldMetaUint64:
		return "UNSIGNED BIG INT"
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

func (d *DialectSqlite3) Pragmas() []string {
	return sqlite3DefaultPragmas
}
