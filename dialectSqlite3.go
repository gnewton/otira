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

func (d *DialectSqlite3) DropTableIfExistsString(tableName string) (string, error) {
	if tableName == "" {
		return "", errors.New("Tablename is empty")
	}
	return "DROP TABLE IF EXISTS " + tableName, nil
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
		fieldType, err := d.FieldType(fm)
		if err != nil {
			return "", err
		}
		s += fm.Name() + " " + fieldType
		if fm == t.primaryKey {
			s += SPC + PRIMARY_KEY
		}
		constraints, err := d.Constraints(fm)
		if err != nil {
			return "", err
		}
		s += constraints
	}
	foreignKeys, err := d.ForeignKeys(t)
	if err != nil {
		return "", err
	}
	s += foreignKeys
	s += ")"
	return s, nil
}

func (d *DialectSqlite3) ForeignKeys(t *TableMeta) (string, error) {
	if t == nil {
		return "", errors.New("TableMeta is nil")
	}
	var s string
	s += d.oneToManyForeignKeys(t)

	return s, nil
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

func (d *DialectSqlite3) Constraints(fm FieldMeta) (string, error) {
	if fm == nil {
		return "", errors.New("FieldMeta is nil")
	}
	if fm.Unique() {
		return " UNIQUE", nil
	}
	return "", nil
}

func (d *DialectSqlite3) FieldType(fm FieldMeta) (string, error) {
	if fm == nil {
		return "", errors.New("FieldMeta is nil")
	}
	switch v := fm.(type) {
	case *FieldMetaString:
		s := ""
		if fm.Fixed() {
			s += "CHAR"
		} else {
			s += "VARCHAR"
		}
		s += "(" + strconv.Itoa(v.Length()) + ")"
		return s, nil
	case *FieldMetaInt, *FieldMetaUint64:
		return "UNSIGNED BIG INT", nil
	case *FieldMetaFloat:
		return "REAL", nil
	case *FieldMetaByte:
		return "BLOB", nil
	default:
		return "", errors.New("Unknown type")
	}

}

func (d *DialectSqlite3) PreparedValueFormat(counter int) (string, error) {
	return "?", nil
}

func (d *DialectSqlite3) Pragmas() []string {
	return sqlite3DefaultPragmas
}

func (d *DialectSqlite3) ExistsString(table string, id uint64) (bool, error) {
	if table == "" {
		return false, errors.New("table name is empty")
	}
	return false, errors.New("TODO")
}

func (d *DialectSqlite3) ExistsDeepString(r *Record) (bool, error) {
	if r == nil {
		return false, errors.New("Record is nil")
	}
	return false, errors.New("TODO")
}
