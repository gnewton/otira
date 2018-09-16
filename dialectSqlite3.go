package otira

import ()

type DialectSqlite3 struct {
}

func (d *DialectSqlite3) CreateTableString(t *TableMeta) string {
	s := "CREATE TABLE " + t.name

	return s
}

func (d *DialectSqlite3) PreparedValueFormat(counter int) string {
	return "?"
}
