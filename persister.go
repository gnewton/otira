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
	incoming           chan *Record
	preparedStatements map[string]*sql.Stmt
	preparedStrings    map[string]string
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
	pers.incoming = make(chan *Record, size)

	pers.preparedStatements = make(map[string]*sql.Stmt, 0)
	pers.preparedStrings = make(map[string]string, 0)

	pers.wg.Add(1)
	go pers.start()

	return pers, nil
}

func (pers *Persister) Save(record *Record) error {
	if record == nil {
		return errors.New("Record cannot be nil")
	}

	pers.saveOneToManyRecords(record)
	pers.savemanyTomanyRecords(record)

	pers.incoming <- record

	return nil
}

func (pers *Persister) Done() error {
	// commit last transation
	// close db
	close(pers.incoming)
	pers.wg.Wait()
	return nil
}

func (pers *Persister) saveOneToManyRecords(record *Record) {

}
func (pers *Persister) savemanyTomanyRecords(record *Record) {

}

func (pers *Persister) start() {
	defer pers.wg.Done()
	n := 0

	for record := range pers.incoming {
		select {
		case <-pers.ctx.Done():
			log.Println("Cancel closing!!!!!!!!!!!")
			return // avoid leaking of this goroutine when ctx is done.
		default:
		}

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
	// var err error
	// pers.tx, err = pers.db.Begin()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// A nil record means a record and all its associated relation record were just sent
}
