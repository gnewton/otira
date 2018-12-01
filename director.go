package otira

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strconv"
	"sync"
)

type Director struct {
	wg         *sync.WaitGroup
	ctx        context.Context
	cancelFunc context.CancelFunc
	incoming   chan []*Record
	persister  *Persister
}

func NewDirector(db *sql.DB, dialect Dialect, size int) (*Director, error) {
	if db == nil {
		return nil, errors.New("DB cannot be nil")
	}

	if dialect == nil {
		return nil, errors.New("Dialect cannot be nil")
	}

	if size < 0 {
		return nil, errors.New("Size must be > 0")
	}

	pers, err := NewPersister(db, dialect)
	if err != nil {
		return nil, err
	}

	dir := new(Director)
	dir.persister = pers
	dir.wg = &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	dir.ctx = ctx
	dir.cancelFunc = cancel
	dir.incoming = make(chan []*Record, size)
	dir.wg.Add(1)
	go dir.start()

	return dir, nil
}

// Saves record and associated relations records
func (dir *Director) Save(record *Record) error {
	if record == nil {
		return errors.New("Record cannot be nil")
	}

	records, err := dir.persister.prepareRelationRecords(record)
	if err != nil {
		return err
	}

	records = append(records, record)

	dir.incoming <- records

	return nil
}

func (dir *Director) Done() error {
	var err error = nil

	err = dir.persister.Done()
	if err != nil {
		log.Println(err)
	}

	close(dir.incoming)
	dir.wg.Wait()
	return nil
}

func (dir *Director) start() {
	defer dir.wg.Done()
	n := 0

	for records := range dir.incoming {
		select {
		case <-dir.ctx.Done():
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
