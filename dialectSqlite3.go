package otira

import (
	"errors"
	"log"
	"strconv"
	"strings"
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

func (d *DialectSqlite3) CreateTableString(t *TableDef) (string, error) {
	if t == nil {
		return "", errors.New("TableDef is nil")
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
		if i == 0 {
			s += SPC + PRIMARY_KEY
		}
		constraints, err := d.Constraints(fm, i == 0)
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

func (d *DialectSqlite3) ForeignKeys(t *TableDef) (string, error) {
	if t == nil {
		return "", errors.New("TableDef is nil")
	}
	var s string
	s += d.oneToManyForeignKeys(t)

	return s, nil
}

func (d *DialectSqlite3) oneToManyForeignKeys(t *TableDef) string {
	var s string
	if t.oneToMany == nil {
		return ""
	}
	log.Println(t.Name())
	for i := 0; i < len(t.oneToMany); i++ {
		one2m := t.oneToMany[i]
		log.Println(one2m)
		s += ", " + "FOREIGN KEY(" + one2m.LeftKeyField.Name() + ") REFERENCES " + one2m.RightTable.Name() + "(" + one2m.RightKeyField.Name() + ")"
	}
	return s
}

func (d *DialectSqlite3) Constraints(fm FieldDef, primaryKey bool) (string, error) {
	if fm == nil {
		return "", errors.New("FieldDef is nil")
	}
	if !primaryKey && fm.Unique() {
		return " UNIQUE", nil
	}
	return "", nil
}

func (d *DialectSqlite3) FieldType(fm FieldDef) (string, error) {
	if fm == nil {
		return "", errors.New("FieldDef is nil")
	}
	switch v := fm.(type) {
	case *FieldDefString:
		s := ""
		if fm.Fixed() {
			s += "CHAR"
		} else {
			s += "VARCHAR"
		}
		s += "(" + strconv.Itoa(v.Length()) + ")"
		return s, nil
	case *FieldDefInt, *FieldDefInt64:
		return "UNSIGNED BIG INT", nil
	case *FieldDefFloat:
		return "REAL", nil
	case *FieldDefByte:
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

func (d *DialectSqlite3) ExistsString(table, field string, id int64) (string, error) {
	if table == "" {
		return "", errors.New("table name is empty")
	}
	return SELECT + SPC + COUNT + "(" + field + ")" + SPC + FROM + SPC + table + SPC + WHERE + SPC + field + EQUALS + toString(id), nil
}

func (d *DialectSqlite3) ExistsDeepString(r *Record) (string, error) {
	if r == nil {
		return "", errors.New("Record is nil")
	}
	log.Println("TODO")
	return "", errors.New("TODO")
}

func (d *DialectSqlite3) UpdateString(rec *Record) (string, error) {
	if rec == nil {
		return "", errors.New("record is nil")
	}
	pk, err := rec.PrimaryKeyValue()
	if err != nil {
		return "", err
	}

	updateValueString, err := d.updateValuesString(rec.fields)
	if err != nil {
		return "", err
	}
	updateString := UPDATE + SPC + rec.tableDef.name + SPC + updateValueString + SPC + WHERE + SPC + rec.fields[0].Name() + EQUALS + toString(pk)
	if rec.tableDef.isJoinTable {
		updateString += SPC + AND + SPC + rec.fields[1].Name() + EQUALS + toString(rec.values[1])
	}
	return updateString, nil
}

func (d *DialectSqlite3) updateValuesString(fields []FieldDef) (string, error) {
	if fields == nil {
		return "", errors.New("Fields is nil")
	}
	s := SET + SPC

	preparedValueFormat, _ := d.PreparedValueFormat(1)

	for i := 0; i < len(fields); i++ {
		if i != 0 {
			s += COMMA
		}
		s += fields[i].Name() + EQUALS + preparedValueFormat
	}
	return s, nil
}

func (d *DialectSqlite3) IsUniqueContraintFailedError(err error) bool {
	return strings.HasPrefix(err.Error(), "UNIQUE constraint failed:")
}

func (d *DialectSqlite3) NumAllowedActiveTransactions() int {
	return 1
}

func (d *DialectSqlite3) OneTransactionPerConnection() bool {
	return true
}

func (d *DialectSqlite3) CreateIndexString(name string, tableName string, fieldNames []string) (string, error) {
	if name == "" {
		return "", errors.New("Empty index name")
	}

	if tableName == "" {
		return "", errors.New("Empty index table name")
	}

	if fieldNames == nil {
		return "", errors.New("Index field names is nil")
	}

	if len(fieldNames) == 0 {
		return "", errors.New("Index field names is zero length")
	}

	for i := 0; i < len(fieldNames); i++ {
		if fieldNames[i] == "" {
			return "", errors.New("Index field name[" + toString(i) + "] is zero length")
		}
	}

	return "", errors.New("TODO")
}
