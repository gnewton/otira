package otira

import (
//"strconv"
)

const INSERT = "INSERT INTO "
const VALUES = "VALUES"

// From: http://go-database-sql.org/prepared.html
// MySQL               PostgreSQL            Oracle
// =====               ==========            ======
// WHERE col = ?       WHERE col = $1        WHERE col = :col
// VALUES(?, ?, ?)     VALUES($1, $2, $3)    VALUES(:val1, :val2, :val3)

// From: https://www.sqlite.org/lang_expr.html#varparam
// ?		A question mark ...

func preparedValueFormat(dialect Dialect, counter int) string {
	return dialect.PreparedValueFormat(counter)

	// case ORACLE:
	// 	return ":val" + strconv.Itoa(counter+1)
	// case POSTGRESQL:
	// 	return "$" + strconv.Itoa(counter+1)
	// default:
	// 	// case MYSQL:
	// 	// case SQLITE3:
	// 	return "?"
	// }

}

func createTableString(dialect Dialect, t *TableMeta) {

}

type Dialect interface {
	CreateTableString(t *TableMeta) string
	PreparedValueFormat(counter int) string
	FieldType(FieldMeta) string
	Constraints(FieldMeta) string
	ForeignKeys(t *TableMeta) string
	InitPragmas() []string
}
