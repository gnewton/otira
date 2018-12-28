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
	pers.TransactionSize = 50000

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
		if tm.created {
			continue
		}

		err := pers.createTable(tm)
		if err != nil {
			return err
		}

	}
	pers.doneCreatingTables = true
	err := pers.BeginTx()
	return err

}

func (pers *Persister) createTable(tm *TableMeta) error {
	log.Println("createTable: " + tm.name)
	createTableString, err := tm.createTableString(pers.dialect)
	if err != nil {
		log.Println(err)
		return err
	}

	// Delete table
	sql, err := pers.dialect.DropTableIfExistsString(tm.name)
	if err != nil {
		return err
	}
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
	tm.created = true
	err = pers.createRelationTables(tm)
	return err
}

func (pers *Persister) createRelationTables(tm *TableMeta) error {
	pers.createOneToManyTables(tm)
	err := pers.createManyToManyTables(tm)
	return err
}

func (pers *Persister) createOneToManyTables(tm *TableMeta) error {
	return nil
}

func (pers *Persister) createManyToManyTables(tm *TableMeta) error {
	log.Println("Create M2M")
	log.Println(len(tm.manyToMany))
	for i := 0; i < len(tm.manyToMany); i++ {
		m2m := tm.manyToMany[i]
		log.Println(m2m.leftTable.name, m2m.rightTable.name)
		log.Println(m2m.joinTable)
		err := pers.createTable(m2m.joinTable)
		if err != nil {
			return err
		}
	}
	return nil
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
	pers.saveMutex.Lock()
	defer pers.saveMutex.Unlock()
	err := pers.tx.Commit()
	pers.tx = nil
	return err
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
		//log.Println("cached")
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
		stmt, err = pers.tx.Prepare(preparedString)
		if err != nil {
			return nil, err
		}
		pers.preparedStatements[record.tableMeta.name] = stmt
		pers.preparedStrings[record.tableMeta.name] = preparedString
	} else {
		//log.Println("Prepared statement cache hit")
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
	var err error
	for i := 0; i < len(record.relationRecords); i++ {
		relation := record.relationRecords[i].relation
		relRecord := record.relationRecords[i].record
		switch rel := relation.(type) {
		case *OneToMany:
			err = pers.saveOneToMany(rel, record, relRecord)
			if err != nil {
				return err
			}
		case *ManyToMany:
			err = pers.saveManyToMany(rel, record, relRecord)
			if err != nil {
				return err
			}
		}
	}
	return err

}

func (pers *Persister) saveManyToMany(m2m *ManyToMany, record *Record, relRecord *Record) error {
	_, exists, err := m2m.cache.GetJoinKey(relRecord)
	if err != nil {
		return err
	}

	if !exists {
		err := pers.Save(relRecord)
		if err != nil {
			return err
		}
	}
	recordPk, ok := record.values[0].(uint64)
	if !ok {
		return errors.New("Record value is not a uint64")
	}
	relRecordPk, ok := relRecord.values[0].(uint64)
	if !ok {
		return errors.New("Relation record value is not a uint64")
	}
	err = pers.saveJoinRecord(m2m, recordPk, relRecordPk)
	if err != nil {
		return err
	}
	return nil
}

func (pers *Persister) saveJoinRecord(m2m *ManyToMany, left, right uint64) error {

	rec, err := m2m.joinTable.NewRecord()
	if err != nil {
		return err
	}
	rec.values[0] = left
	rec.values[1] = right
	err = pers.save(rec)
	if err != nil {
		return err
	}
	return nil
}

func (pers *Persister) saveOneToMany(one2m *OneToMany, record *Record, relRecord *Record) error {
	k2, exists, err := one2m.cache.GetJoinKey(relRecord)
	if err != nil {
		return err
	}
	//relRecordValueIndex, ok := relRecord.fieldsMap[one2m.rightKeyField.Name()]
	//if !ok {
	//return errors.New("Cannot find relation record primary key")
	//} else {
	//record.SetByName(one2m.leftKeyField.Name(), relRecord.values[relRecordValueIndex])
	record.SetByName(one2m.leftKeyField.Name(), k2)
	//}
	if !exists {
		pers.Save(relRecord)
	}

	return nil
}

// Saves single record
func (pers *Persister) save(record *Record) error {

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
		log.Println("End TX; start new TX")
		pers.transactionCounter = 0
	}
	pers.transactionCounter++

	stmt, err := pers.preparedStatement(record)
	//log.Println("SAVE() saving record:")
	//log.Println(record.String())
	if err != nil {
		return err
	}
	pers.saveMutex.Lock()

	result, err := execStatement(stmt, record.Values())
	if err != nil {
		log.Println(err)
		log.Println(record.values[0])
		return err
	}
	pers.saveMutex.Unlock()
	// lastInsertId, err := result.LastInsertId()
	// if err != nil {
	// 	return err
	// }
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return errors.New("Wrong number of rows effected by single insert...")
	}
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
		preparedValueFormat, err := preparedValueFormat(pers.dialect, i)
		if err != nil {
			return "", err
		}
		values += preparedValueFormat
	}

	st += ")"
	values += ")"

	st = st + " VALUES " + values
	return st, nil
}
