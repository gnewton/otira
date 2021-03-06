package otira

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"
)

type Persister struct {
	TransactionSize        int
	createMutex            sync.Mutex
	db                     *sql.DB
	dialect                Dialect
	doneCreatingTables     bool
	preparedStatementCache map[string]*sql.Stmt
	preparedStrings        map[string]string
	relationPKCacheMap     map[Relation]map[string]int64
	saveMutex              sync.Mutex
	transactionCounter     int
	tx                     *sql.Tx
	SupportUpdates         bool
}

// needs to also have
//func NewPersister(ctx context.Context, db *sql.DB, dialect Dialect, size int) (*Persister, error) {
func NewPersister(db *sql.DB, dialect Dialect, txSize int) (*Persister, error) {
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
	pers.TransactionSize = txSize
	log.Println("txSize=", pers.TransactionSize)

	pers.preparedStatementCache = make(map[string]*sql.Stmt, 0)
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

func (pers *Persister) DeleteTables(tms ...*TableDef) error {
	pers.createMutex.Lock()
	defer pers.createMutex.Unlock()

	if tms == nil {
		return errors.New("CreateTables: Table list is nil")
	}
	if len(tms) == 0 {
		return errors.New("CreateTables: Table list is empty")
	}
	if pers.doneCreatingTables {
		return errors.New("Cannot delete tables that were created in this instance")
	}
	return nil
}

func verifyTableDef(pers *Persister, tms ...*TableDef) error {
	for i := 0; i < len(tms); i++ {
		tm := tms[i]
		if tm == nil {
			return errors.New("TableDef cannot be nil")
		}
		if tm.created {
			continue
		}

		err := pers.deleteTable(tm)
		//err := createTable(pers.db, pers.dialect, tm)

		if err != nil {
			return err
		}

	}
	for i := 0; i < len(tms); i++ {
		tm := tms[i]

		err := pers.deleteRelationTables(tm)
		if err != nil {
			return err
		}

	}
	return nil
}

func (pers *Persister) CreateTables(tms ...*TableDef) error {
	pers.createMutex.Lock()
	defer pers.createMutex.Unlock()
	if err := verifyTableDef(pers, tms...); err != nil {
		return err
	}

	for i := 0; i < len(tms); i++ {
		tm := tms[i]
		if tm == nil {
			return errors.New("TableDef cannot be nil")
		}
		if tm.created {
			continue
		}

		//err := pers.createTable(tm)
		err := createTable(pers.db, pers.dialect, tm)

		if err != nil {
			return err
		}

	}
	for i := 0; i < len(tms); i++ {
		tm := tms[i]

		err := pers.createRelationTables(tm)
		if err != nil {
			return err
		}

	}
	pers.doneCreatingTables = true
	err := pers.BeginTx()
	return err

}

func (pers *Persister) deleteTable(tm *TableDef) error {
	// Delete table
	sql, err := pers.dialect.DropTableIfExistsString(tm.name)
	if err != nil {
		return err
	}
	_, err = exec(pers.db, sql)
	if err != nil {
		return err
	}

	return nil
}

// func (pers *Persister) createTable(tm *TableDef) error {
// 	err := createTable(pers.db, pers.dialect, tm)
// 	return err
// }

func (pers *Persister) deleteRelationTables(tm *TableDef) error {
	err := pers.deleteOneToManyTables(tm)
	if err != nil {
		return err
	}
	err = pers.deleteManyToManyTables(tm)
	return err
}

func (pers *Persister) deleteOneToManyTables(tm *TableDef) error {
	//TODO

	log.Println("TODO")
	//return errors.New("TODO")
	return nil
}

func (pers *Persister) deleteManyToManyTables(tm *TableDef) error {
	log.Println("Delete M2M")
	log.Println(len(tm.manyToMany))
	for i := 0; i < len(tm.manyToMany); i++ {
		m2m := tm.manyToMany[i]
		if m2m == nil {
			return errors.New("m2m tabledef is nil")
		}
		log.Println(m2m.LeftTable.name, m2m.RightTable.name)
		log.Println(m2m.JoinTable)
		// FIXX this should be pers.CreateTable: this table migh have relation tables too
		err := pers.deleteTable(m2m.JoinTable)
		if err != nil {
			return err
		}
	}
	return nil
}

func (pers *Persister) createRelationTables(tm *TableDef) error {
	err := pers.createOneToManyTables(tm)
	if err != nil {
		return err
	}
	err = pers.createManyToManyTables(tm)
	return err
}

func (pers *Persister) createOneToManyTables(tm *TableDef) error {
	//TODO

	log.Println("TODO")
	//return errors.New("TODO")
	return nil
}

func (pers *Persister) createManyToManyTables(tm *TableDef) error {
	log.Println("Create M2M")
	log.Println(len(tm.manyToMany))
	for i := 0; i < len(tm.manyToMany); i++ {
		m2m := tm.manyToMany[i]
		log.Println("leftTable:", m2m.LeftTable.name, "rightTable:", m2m.RightTable.name)
		log.Println("jointTable", m2m.JoinTable.name)

		// FIXX this should be pers.CreateTable: this table migh have relation tables too
		//err := pers.createTable(m2m.JoinTable)
		//err := pers.createTable(m2m.JoinTable)
		err := createTable(pers.db, pers.dialect, m2m.JoinTable)
		if err != nil {
			return err
		}
		err = createTable(pers.db, pers.dialect, m2m.RightTable)
		if err != nil {
			return err
		}
		err = createTable(pers.db, pers.dialect, m2m.LeftTable)
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
	err := pers.closePreparedStatements()
	if err != nil {
		return err
	}

	err = pers.tx.Commit()
	if err != nil {
		return err
	}
	pers.tx = nil
	return err
}

func (pers *Persister) closePreparedStatements() error {
	for _, stmt := range pers.preparedStatementCache {
		if stmt == nil {
			return errors.New("Prepared statement should not be nil")
		}
		log.Println("Closing Prepared Statement", stmt)
		err := stmt.Close()
		if err != nil {
			return err
		}
	}
	return nil
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
	if preparedString, ok = pers.preparedStrings[record.tableDef.name]; !ok {
		preparedString, err = record.tableDef.CreatePreparedStatementInsertAllFields(pers.dialect)
		//log.Println("Prepared String=" + preparedString)
		if err != nil {
			log.Println("error=")
			log.Println(err)
			return "", err
		}
		pers.preparedStrings[record.tableDef.name] = preparedString
	} else {
		//log.Println("cached")
	}
	return preparedString, nil
}

func (pers *Persister) preparedStatement(record *Record) (*sql.Stmt, error) {
	var ok bool
	var stmt *sql.Stmt
	if stmt, ok = pers.preparedStatementCache[record.tableDef.name]; !ok {
		//preparedString, err := record.tableDef.CreatePreparedStatementInsertAllFields(pers.dialect)
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
		pers.preparedStatementCache[record.tableDef.name] = stmt
		pers.preparedStrings[record.tableDef.name] = preparedString
	} else {
		//log.Println("Prepared statement cache hit")
	}
	return stmt, nil
}

func (pers *Persister) commitAndBeginTx() error {
	err := pers.commit()
	if err != nil {
		return err
	}
	err = pers.BeginTx()
	if err != nil {
		return err
	}
	pers.preparedStatementCache = make(map[string]*sql.Stmt, 0)
	log.Println("End TX; start new TX")
	pers.transactionCounter = 0

	return nil
}

// saves record and all related records
func (pers *Persister) Save(rec *Record) error {
	log.Println("Start save: ", rec.tableDef.name, rec.String())

	if unsetPrimaryKey(rec) {
		return fmt.Errorf("Primary key not set for record %s %v", rec.String(), rec.valueIsSet)
	}

	update, err := pers.isUpdate(rec)
	if err != nil {
		return err
	}
	if update {
		return pers.Update(rec)
	}

	err = pers.saveRelations(rec)
	if err != nil {
		return err
	}

	err = pers.save(rec)
	if err != nil {
		return err
	}

	if pers.transactionCounter > pers.TransactionSize {
		err = pers.commitAndBeginTx()
		if err != nil {
			return err
		}

	}
	return nil
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
	_, exists, err := m2m.cache.MakeJoinKey(relRecord)
	if err != nil {
		return err
	}

	if !exists {
		// FIXXX
		err := pers.Save(relRecord)
		if err != nil {
			return err
		}
	}
	recordPk, ok := record.values[0].(int64)
	if !ok {
		return errors.New("Record value is not a int64")
	}
	relRecordPk, ok := relRecord.values[0].(int64)
	if !ok {
		return errors.New("Relation record value is not a int64")
	}
	err = pers.saveJoinRecord(m2m, recordPk, relRecordPk)
	if err != nil {
		return err
	}
	return nil
}

func (pers *Persister) saveJoinRecord(m2m *ManyToMany, left, right int64) error {

	rec, err := m2m.JoinTable.NewRecord()
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
	log.Println("saveOneToMany")
	k2, exists, err := one2m.cache.MakeJoinKey(relRecord)
	if err != nil {
		return err
	}
	//relRecordValueIndex, ok := relRecord.fieldsMap[one2m.rightKeyField.Name()]
	//if !ok {
	//return errors.New("Cannot find relation record primary key")
	//} else {
	//record.SetByName(one2m.leftKeyField.Name(), relRecord.values[relRecordValueIndex])
	record.SetByName(one2m.LeftKeyField.Name(), k2)
	//}
	if !exists {
		pers.Save(relRecord)
	} else {
		log.Println("In cache")
	}

	return nil
}

// Saves single record
func (pers *Persister) save(rec *Record) error {
	pers.saveMutex.Lock()
	pers.transactionCounter++

	stmt, err := pers.preparedStatement(rec)
	if err != nil {
		return err
	}

	defer pers.saveMutex.Unlock()

	if pers.SupportUpdates {
		// pk, err := rec.PrimaryKeyValue()
		// if err != nil {
		// 	return err
		// }
		// err = rec.tableDef.writeCache.Put(pk)
		// if err != nil {
		// 	return err
		// }
	}

	result, err := execStatement(stmt, rec.Values())
	if err != nil {
		if pers.dialect.IsUniqueContraintFailedError(err) {
			log.Println("Updating:" + rec.tableDef.name + "   Primary key=" + toString(rec.values[0]))
			err := pers.Update(rec)
			if err != nil {
				log.Println(err)
				log.Println("\tTable:" + rec.tableDef.name + "   Primary key=" + toString(rec.values[0]))
				return err
			}
			return nil
		} else {
			log.Println(err)
			log.Println("\tTable:" + rec.tableDef.name + "   Primary key=" + toString(rec.values[0]))
			return err
		}
	}

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
	a, b := pers.CreatePreparedStatementInsertSomeFields(record.tableDef.name, record.tableDef.fields...)
	return a, b
}

func (pers *Persister) CreatePreparedStatementInsertFromRecord(record *Record) (string, error) {
	if record == nil {
		return "", errors.New("Record cannot be nil")
	}
	fields := make([]FieldDef, 0)
	for i, _ := range record.values {
		fields = append(fields, record.tableDef.fields[i])
	}
	return pers.CreatePreparedStatementInsertSomeFields(record.tableDef.name, fields...)
}

func (pers *Persister) Exists(r *Record) (bool, error) {
	pkValue, err := r.PrimaryKeyValue()
	if err != nil {
		return false, err
	}

	s, err := pers.dialect.ExistsString(r.tableDef.name, r.tableDef.fields[0].Name(), pkValue)
	if err != nil {
		return false, err
	}

	stmt, err := pers.tx.Prepare(s)
	if err != nil {
		return false, err
	}
	var count int
	err = stmt.QueryRow().Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (pers *Persister) CreatePreparedStatementInsertSomeFields(tablename string, fields ...FieldDef) (string, error) {
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

func (pers *Persister) isUpdate(rec *Record) (bool, error) {

	if !pers.SupportUpdates {
		return false, nil
	}

	pk, err := rec.PrimaryKeyValue()
	if err != nil {
		return false, err
	}
	log.Println("Writing")
	log.Println(pk)

	ok, err := rec.tableDef.writeCache.Contains(pk)
	if ok {
		return true, nil
	}
	return false, nil
}

func (pers *Persister) Update(rec *Record) error {
	//TODO
	// - updateRelations
	// - updateRecord
	return pers.update(rec)
}

func (pers *Persister) update(rec *Record) error {
	log.Println("TODO update")
	updateString, err := pers.dialect.UpdateString(rec)
	if err != nil {
		return err
	}
	log.Println(updateString)

	stmt, err := pers.tx.Prepare(updateString)
	if err != nil {
		return err
	}
	log.Println(stmt)
	_, err = execStatement(stmt, rec.values)

	log.Println("DONE")
	if err != nil {
		return err
	}
	return nil
}
