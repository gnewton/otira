package otira

import (
	"database/sql"
	"errors"
	"log"
)

func createTable(db *sql.DB, dialect Dialect, tm *TableDef) error {
	if tm == nil {
		return errors.New("createTable: Table is nil")
	}

	if db == nil {
		return errors.New("createTable: DB is nil")
	}

	createTableString, err := tm.createTableString(dialect)
	log.Println("CREATE::::::::: " + createTableString)
	if err != nil {
		log.Println(err)
		return err
	}

	// Delete table
	sql, err := dialect.DropTableIfExistsString(tm.name)
	if err != nil {
		return err
	}
	_, err = exec(db, sql)
	if err != nil {
		return err
	}
	log.Println("createTableString=" + createTableString)

	// Create the table in the db
	_, err = exec(db, createTableString)
	if err != nil {
		return err
	}
	tm.created = true
	return err
}
