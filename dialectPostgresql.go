package otira

import (
	"errors"
	"log"
	"strconv"
)

type DialectPostgresql struct {
}

func (d *DialectPostgresql) CreateTableString(t *TableMeta) (string, error) {
	if t == nil {
		return "", errors.New("TableMeta is nil")
	}

	log.Println("TODO")
	return "", errors.New("TODO")
}

func (d *DialectPostgresql) PreparedValueFormat(counter int) (string, error) {
	if counter < 0 {
		return "", errors.New("counter is <1")
	}
	return "$" + strconv.Itoa(counter+1), nil
}

func (d *DialectPostgresql) FieldType(fm FieldMeta) (string, error) {
	if fm == nil {
		return "", errors.New("FieldMeta is nil")
	}
	return "", errors.New("TODO")
}

func (d *DialectPostgresql) Constraints(fm FieldMeta) (string, error) {
	if fm == nil {
		return "", errors.New("FieldMeta is nil")
	}
	return "", errors.New("TODO")
}

func (d *DialectPostgresql) ForeignKeys(t *TableMeta) (string, error) {
	if t == nil {
		return "", errors.New("TableMeta is nil")
	}
	return "", errors.New("TODO")

}

//TODO
func (d *DialectPostgresql) Pragmas() []string {
	var pragmas []string

	return pragmas
}

//TODO
func (d *DialectPostgresql) DropTableIfExistsString(tableName string) (string, error) {
	return "", errors.New("TODO")
}

func (d *DialectPostgresql) ExistsString(table, field string, id uint64) (string, error) {
	return "", errors.New("TODO")
}

func (d *DialectPostgresql) ExistsDeepString(*Record) (string, error) {
	return "", errors.New("TODO")
}
