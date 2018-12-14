package otira

import (
	"database/sql"
	"errors"
	"log"
	"sync"
)

type Persister struct {
	TransactionSize    int
	transactionCounter int
	dialect            Dialect
	db                 *sql.DB
	tx                 *sql.Tx
	preparedStatements map[string]*sql.Stmt
	preparedStrings    map[string]string

	relationPKCacheMap map[Relation]map[string]int64

	saveMutex          sync.Mutex
	createMutex        sync.Mutex
	doneCreatingTables bool
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

	if err != nil {
		return nil, err
	}

	pers.initPragmas()
	pers.TransactionSize = 500

	//err = pers.beginTx()
	pers.preparedStatements = make(map[string]*sql.Stmt, 0)
	pers.preparedStrings = make(map[string]string, 0)
	pers.relationPKCacheMap = make(map[Relation]map[string]int64)
	pers.doneCreatingTables = false
	return pers, nil
}

func (pers *Persister) initPragmas() error {
	if pers.db == nil {
		return errors.New("db is nil")
	}
	pragmas := pers.dialect.Pragmas()
	for i := 0; i < len(pragmas); i++ {
		_, err := exec(pers.db, pragmas[i])
		log.Println("PRAGMAS: " + pragmas[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (pers *Persister) CreateTables(tms ...*TableMeta) error {
	if pers.doneCreatingTables {
		return errors.New("CreateTables can only be called once")
	}

	pers.createMutex.Lock()
	defer pers.createMutex.Unlock()

	if pers.db == nil {
		return errors.New("db cannot be nil")
	}
	if tms == nil {
		return errors.New("TableMeta cannot be nil")
	}

	if pers.dialect == nil {
		return errors.New("Dialect cannot be nil")
	}

	for i := 0; i < len(tms); i++ {
		tm := tms[i]
		createTableString, err := tm.createTableString(pers.dialect)
		if err != nil {
			return err
		}

		// Delete table
		sql := pers.dialect.DropTableIfExists(tm)
		log.Println(sql)
		_, err = exec(pers.db, sql)
		if err != nil {
			return err
		}
		log.Println("createTableString=" + createTableString)
		// Create the table in the db
		_, err = exec(pers.db, createTableString)
		if err != nil {
			return err
		}
	}
	pers.doneCreatingTables = true
	err := pers.BeginTx()
	return err

}

func exec(db *sql.DB, sql string) (sql.Result, error) {
	return db.Exec(sql)
}

func execStatement(stmt *sql.Stmt, values []interface{}) (sql.Result, error) {
	return stmt.Exec(values...)
}

func (pers *Persister) BeginTx() error {
	if pers.tx != nil {
		return nil
	}
	log.Println("=============================START TX")
	pers.saveMutex.Lock()
	defer pers.saveMutex.Unlock()

	var err error
	pers.tx, err = pers.db.Begin()
	pers.transactionCounter = 0
	return err
}

func (pers *Persister) commit() error {
	log.Println("=============================END TX")
	return pers.tx.Commit()
}

func (pers *Persister) Done() error {
	// commit last transation
	// close db
	if pers.tx != nil {
		return pers.commit()
	} else {
		return nil
	}

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

// saves record and all related records
func (pers *Persister) Save(record *Record) error {
	err := pers.saveRelations(record)
	if err != nil {
		return err
	}
	err = pers.save(record)
	return err
}

func (pers *Persister) saveRelations(record *Record) error {
	log.Println("saveRelations")
	err := pers.saveOneToMany(record)
	//saveManyToMany(record)
	return err

}

func (pers *Persister) saveOneToMany(record *Record) error {
	log.Println("saveOneToMany")
	for i := 0; i < len(record.relationRecords); i++ {
		relation := record.relationRecords[i].relation
		relRecord := record.relationRecords[i].record
		log.Println(relation)
		if one2m, ok := relation.(*OneToMany); ok {
			log.Println(i)
			log.Println("*****************")
			//relRecordValueIndex, ok := relRecord.fieldsMap[record.tableMeta.PrimaryKey().Name()]
			relRecordValueIndex, ok := relRecord.fieldsMap[one2m.rightKeyField.Name()]
			if !ok {
				return errors.New("Cannot find relation record primary key")
			} else {
				record.SetByName(one2m.leftKeyField.Name(), relRecord.values[relRecordValueIndex])
			}
			pers.save(relRecord)
		}
	}
	return nil
}

// Saves single record
func (pers *Persister) save(record *Record) error {
	pers.saveMutex.Lock()
	defer pers.saveMutex.Unlock()

	if pers.transactionCounter > pers.TransactionSize {
		err := pers.commit()
		if err != nil {
			return err
		}
		err = pers.BeginTx()
		if err != nil {
			return err
		}
		pers.preparedStatements = make(map[string]*sql.Stmt, 0)
	}

	stmt, err := pers.preparedStatement(record)
	log.Println("SAVE()")
	//log.Println(record.Values())
	log.Println(record.String())
	if err != nil {
		return err
	}
	result, err := execStatement(stmt, record.Values())
	if err != nil {
		log.Println(err)
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
