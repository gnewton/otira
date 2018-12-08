package otira

import (
	"database/sql"
	"errors"
	"log"
)

type Persister struct {
	dialect            Dialect
	db                 *sql.DB
	tx                 *sql.Tx
	preparedStatements map[string]*sql.Stmt
	preparedStrings    map[string]string

	relationPKCacheMap map[Relation]map[string]int64
}

// needs to also have
//func NewPersister(ctx context.Context, db *sql.DB, dialect Dialect, size int) (*Persister, error) {
func NewPersister(db *sql.DB, dialect Dialect) (*Persister, error) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var err error
	//log.SetFlags(log.LstdFlags | log.Lshortfile)
	if db == nil {
		return nil, errors.New("DB cannot be nil")
	}

	if dialect == nil {
		return nil, errors.New("Dialect cannot be nil")
	}

	pers := new(Persister)
	pers.db = db
	pers.dialect = dialect
	pers.initPragmas()

	//pers.ctx = ctx

	err = pers.beginTx()
	if err != nil {
		return nil, err
	}

	pers.preparedStatements = make(map[string]*sql.Stmt, 0)
	pers.preparedStrings = make(map[string]string, 0)
	pers.relationPKCacheMap = make(map[Relation]map[string]int64)

	return pers, nil
}

func (pers *Persister) initPragmas() error {
	pragmas := pers.dialect.Pragmas()
	for i := 0; i < len(pragmas); i++ {
		_, err := pers.db.Exec(pragmas[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (pers *Persister) CreateTable(tm *TableMeta) error {
	if pers.db == nil {
		return errors.New("db cannot be nil")
	}
	if pers.tx == nil {
		return errors.New("Tx cannot be nil")
	}
	if tm == nil {
		return errors.New("TableMeta cannot be nil")
	}

	if pers.dialect == nil {
		return errors.New("Dialect cannot be nil")
	}
	createTableString, err := tm.createTableString(pers.dialect)
	if err != nil {
		return err
	}

	log.Println("createTableString=" + createTableString)
	// Create the table in the db
	_, err = pers.tx.Exec(createTableString)

	err = pers.Commit()
	if err != nil {
		return err
	}
	err = pers.BeginTx()
	return err

}

func (pers *Persister) beginTx() error {
	var err error
	pers.tx, err = pers.db.Begin()
	return err
}

func (pers *Persister) commit() error {
	var err error
	err = pers.tx.Commit()
	return err

}

func (pers *Persister) BeginTx() error {
	var err error
	pers.tx, err = pers.db.Begin()
	return err
}

func (pers *Persister) Commit() error {
	return pers.tx.Commit()
}

func (pers *Persister) Done() error {
	// commit last transation
	// close db
	return pers.Commit()
}

func (pers *Persister) prepareRelationRecords(record *Record) ([]*Record, error) {
	var relationRecords []*Record
	//relationRecords := make([]*Record, 0)
	for i := 0; i < len(record.relationRecords); i++ {
		rr := record.relationRecords[i]

		switch v := rr.relation.(type) {
		case *OneToMany:
			relationRecords = append(relationRecords, pers.prepareOneToManyRecord(record, rr.record, v)...)
		case *ManyToMany:
			relationRecords = append(relationRecords, pers.prepareManyToManyRecord(record, rr.record, v)...)
			log.Println(v.String())
		}
	}
	return relationRecords, nil
}

func (pers *Persister) prepareOneToManyRecord(record *Record, relationRecord *Record, relation *OneToMany) []*Record {
	var relationPKCache map[string]int64
	var ok bool
	oneToManyRecords := make([]*Record, 0)
	if relationPKCache, ok = pers.relationPKCacheMap[relation]; !ok {
		relationPKCache = make(map[string]int64)
		pers.relationPKCacheMap[relation] = relationPKCache
	}

	//k, err := makeKey(relationRecord, relation)
	_, _ = findRelationPK(relationRecord, relation)
	return oneToManyRecords
}

func (pers *Persister) prepareManyToManyRecord(record *Record, relationRecord *Record, relation *ManyToMany) []*Record {
	manyToManyRecords := make([]*Record, 0)
	return manyToManyRecords
}

func (pers *Persister) preparedString(record *Record) (string, error) {
	var ok bool
	var preparedString string
	var err error
	if preparedString, ok = pers.preparedStrings[record.tableMeta.name]; !ok {
		preparedString, err = record.tableMeta.CreatePreparedStatementInsertAllFields(pers.dialect)
		//log.Println("Prepared String=" + preparedString)
		if err != nil {
			log.Println("error=")
			log.Println(err)
			return "", err
		}
		pers.preparedStrings[record.tableMeta.name] = preparedString
	} else {
		log.Println("cached")
	}
	return preparedString, nil
}

func (pers *Persister) preparedStatement(record *Record) (*sql.Stmt, error) {
	var ok bool
	var stmt *sql.Stmt
	if stmt, ok = pers.preparedStatements[record.tableMeta.name]; !ok {
		//preparedString, err := record.tableMeta.CreatePreparedStatementInsertAllFields(pers.dialect)
		preparedString, err := pers.preparedString(record)
		if err != nil {
			return nil, err
		}
		if pers.tx == nil {
			log.Println("*******************  pers.tx is nil")
			return nil, errors.New("pers.tx is nil")
		}
		log.Println("preparedString=[" + preparedString + "]")
		stmt, err = pers.tx.Prepare(preparedString)
		if err != nil {
			return nil, err
		}
		pers.preparedStatements[record.tableMeta.name] = stmt
		pers.preparedStrings[record.tableMeta.name] = preparedString
	}
	return stmt, nil
}

// Saves single record
func (pers *Persister) save(record *Record) error {
	stmt, err := pers.preparedStatement(record)
	log.Println("SAVE()")
	//log.Println(record.Values())
	log.Println(record.String())
	if err != nil {
		return err
	}
	result, err := stmt.Exec(record.Values()...)
	if err != nil {
		return err
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return err
	}
	log.Println(lastInsertId)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	log.Println(rowsAffected)
	return err
}

func (pers *Persister) CreatePreparedStatementInsertAllFields(record *Record) (string, error) {
	a, b := pers.CreatePreparedStatementInsertSomeFields(record.tableMeta.name, record.tableMeta.fields...)
	return a, b
}

func (pers *Persister) CreatePreparedStatementInsertFromRecord(record *Record) (string, error) {
	if record == nil {
		return "", errors.New("Record cannot be nil")
	}
	fields := make([]FieldMeta, 0)
	for i, _ := range record.values {
		fields = append(fields, record.tableMeta.fields[i])
	}
	return pers.CreatePreparedStatementInsertSomeFields(record.tableMeta.name, fields...)
}

func (pers *Persister) CreatePreparedStatementInsertSomeFields(tablename string, fields ...FieldMeta) (string, error) {
	st := "INSERT INTO " + tablename + " ("
	values := "("
	for i, _ := range fields {
		if i != 0 {
			st += ", "
			values += ", "
		}
		st += fields[i].Name()
		values += preparedValueFormat(pers.dialect, i)
	}

	st += ")"
	values += ")"

	st = st + " VALUES " + values
	return st, nil
}
