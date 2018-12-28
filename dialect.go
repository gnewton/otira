package otira

import (
	//"strconv"
	"errors"
)

const CREATE_TABLE = "CREATE TABLE"
const SPC = " "
const PRIMARY_KEY = "PRIMARY_KEY"
const INSERT = "INSERT INTO "
const VALUES = "VALUES"

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

func createTableString(dialect Dialect, t *TableMeta) {

}

type Dialect interface {
	Constraints(FieldMeta) (string, error)
	CreateTableString(t *TableMeta) (string, error)
	DropTableIfExistsString(tableName string) (string, error)
	ExistsString(table string, id uint64) (bool, error)
	ExistsDeepString(*Record) (bool, error)
	FieldType(FieldMeta) (string, error)
	ForeignKeys(t *TableMeta) (string, error)
	Pragmas() []string
	PreparedValueFormat(counter int) (string, error)
}
