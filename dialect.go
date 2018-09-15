package otira

import (
	"strconv"
)

// From: http://go-database-sql.org/prepared.html
// MySQL               PostgreSQL            Oracle
// =====               ==========            ======
// WHERE col = ?       WHERE col = $1        WHERE col = :col
// VALUES(?, ?, ?)     VALUES($1, $2, $3)    VALUES(:val1, :val2, :val3)

// From: https://www.sqlite.org/lang_expr.html#varparam
// ?		A question mark ...

func preparedValueFormat(dialect string, counter int) string {
	switch dialect {
	case "oracle":
		return ":val" + strconv.Itoa(counter+1)
	case "postgresql":
		return "$" + strconv.Itoa(counter+1)
	default:
		// case "mysql":
		// case "sqlite3":
		return "?"
	}

}
