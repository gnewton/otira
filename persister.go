package otira

import (
	"database/sql"
	"errors"
	"log"
)

type Persister struct {
	dialect  Dialect
	tx       *sql.Tx
	incoming chan *Record
}

// needs to also have
func NewPersister(db *sql.DB, dialect Dialect, size int) (*Persister, error) {
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
	pers.incoming = make(chan *Record, size)

	go pers.start()

	return pers, nil
}

func (pers *Persister) Save(record *Record) error {
	if record == nil {
		return errors.New("Record cannot be nil")
	}
	pers.incoming <- record
	return nil
}

func (pers *Persister) start() {

	for record := range pers.incoming {
		log.Println(record.tableMeta.name)
	}

}
