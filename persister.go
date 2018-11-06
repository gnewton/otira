package otira

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strconv"
	"sync"
)

type Persister struct {
	wg                 *sync.WaitGroup
	ctx                context.Context
	cancelFunc         context.CancelFunc
	dialect            Dialect
	db                 *sql.DB
	tx                 *sql.Tx
	incoming           chan []*Record
	preparedStatements map[string]*sql.Stmt
	preparedStrings    map[string]string

	relationPKCacheMap map[Relation]map[string]int64
}

// needs to also have
//func NewPersister(ctx context.Context, db *sql.DB, dialect Dialect, size int) (*Persister, error) {
func NewPersister(db *sql.DB, dialect Dialect, size int) (*Persister, error) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if db == nil {
		return nil, errors.New("DB cannot be nil")
	}

	if dialect == nil {
		return nil, errors.New("Dialect cannot be nil")
	}

	if size < 0 {
		return nil, errors.New("Size must be > 0")
	}

	pers := new(Persister)
	pers.wg = &sync.WaitGroup{}

	ctx, cancel := context.WithCancel(context.Background())
	pers.ctx = ctx
	pers.cancelFunc = cancel
	//pers.ctx = ctx
	pers.db = db
	pers.dialect = dialect
	pers.incoming = make(chan []*Record, size)

	pers.preparedStatements = make(map[string]*sql.Stmt, 0)
	pers.preparedStrings = make(map[string]string, 0)
	pers.relationPKCacheMap = make(map[Relation]map[string]int64)

	pers.wg.Add(1)
	go pers.start()

	return pers, nil
}

// Saves record and associated relations records
func (pers *Persister) Save(record *Record) error {
	if record == nil {
		return errors.New("Record cannot be nil")
	}

	records, err := pers.prepareRelationRecords(record)
	if err != nil {
		return err
	}

	records = append(records, record)

	pers.incoming <- records

	return nil
}

func (pers *Persister) Done() error {
	// commit last transation
	// close db
	close(pers.incoming)
	pers.wg.Wait()
	return nil
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

func (pers *Persister) preparedStatement(record *Record) (*sql.Stmt, error) {
	var ok bool
	var stmt *sql.Stmt
	if stmt, ok = pers.preparedStatements[record.tableMeta.name]; !ok {
		preparedString, err := record.tableMeta.CreatePreparedStatementInsertAllFields(pers.dialect)
		if err != nil {
			return nil, err
		}
		stmt, err = pers.tx.Prepare(preparedString)
		if err != nil {
			return nil, err
		}
		pers.preparedStatements[record.tableMeta.name] = stmt
	}
	return stmt, nil
}

func (pers *Persister) save(record *Record) error {
	stmt, err := pers.preparedStatement(record)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(record.Values()...)
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

func (pers *Persister) start() {
	defer pers.wg.Done()
	n := 0

	for records := range pers.incoming {
		select {
		case <-pers.ctx.Done():
			log.Println("Cancel closing!!!!!!!!!!!")
			return // avoid leaking of this goroutine when ctx is done.
		default:
		}

		for i := 0; i < len(records); i++ {
			record := records[i]
			n++
			if record != nil {
				tableName := record.tableMeta.name
				// if preparedStatement, ok := pers.preparedStatements[tableName]; !ok {
				// 	if preparedString, ok2 := pers.preparedStrings[tableName]; !ok {
				// 		preparedString, _ = record.tableMeta.CreatePreparedStatementInsertFromRecord(pers.dialect, record)
				// 		pers.preparedStrings[tableName] = preparedString

				// 	}
				// 	preparedStatement, err = record.tx.Prepare(r.preparedString)
				// }
				log.Println(strconv.Itoa(n) + " FOO " + tableName)
			} else {
				log.Println(strconv.Itoa(n) + " A--nil--")
			}
			//pers.Save(record)

		}
	}
	// var err error
	// pers.tx, err = pers.db.Begin()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// A nil record means a record and all its associated relation record were just sent
}
