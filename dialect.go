package otira

import (
	//"strconv"
	"errors"
)

type Dialect interface {
	Constraints(FieldDef, bool) (string, error)
	CreateIndexString(string, string, []string) (string, error)
	CreateTableString(t *TableDef) (string, error)
	DropTableIfExistsString(tableName string) (string, error)
	ExistsDeepString(*Record) (string, error)
	ExistsString(table, field string, id int64) (string, error)
	FieldType(FieldDef) (string, error)
	ForeignKeys(t *TableDef) (string, error)

	IsUniqueContraintFailedError(error) bool
	Pragmas() []string
	PreparedValueFormat(counter int) (string, error)
	UpdateString(*Record) (string, error)

	NumAllowedActiveTransactions() int
	OneTransactionPerConnection() bool
}

const CREATE_TABLE = "CREATE TABLE"
const SPC = " "
const AND = "AND"
const PRIMARY_KEY = "PRIMARY KEY"
const INSERT = "INSERT INTO "
const VALUES = "VALUES"
const SELECT = "SELECT"
const WHERE = "WHERE"
const FROM = "FROM"
const COUNT = "COUNT"
const EQUALS = "="
const UPDATE = "UPDATE"
const SET = "SET"
const COMMA = ","

// From: http://go-database-sql.org/prepared.html
// MySQL               PostgreSQL            Oracle
// =====               ==========            ======
// WHERE col = ?       WHERE col = $1        WHERE col = :col
// VALUES(?, ?, ?)     VALUES($1, $2, $3)    VALUES(:val1, :val2, :val3)

// From: https://www.sqlite.org/lang_expr.html#varparam
// ?		A question mark ...

func preparedValueFormat(dialect Dialect, counter int) (string, error) {
	if dialect == nil {
		return "", errors.New("Dialect is nil")
	}

	format, err := dialect.PreparedValueFormat(counter)
	if err != nil {
		return "", err
	}

	return format, nil

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

func createTableString(dialect Dialect, t *TableDef) {

}
