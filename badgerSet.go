package otira

import (
	"errors"
	"github.com/dgraph-io/badger"
	"os"
)

type badgerSet struct {
	db      *badger.DB
	txn     *badger.Txn
	dir     string
	counter uint32
}

func NewBadgerSet(dir string) (Set, error) {
	if dir == "" {
		return nil, errors.New("HashCache dir is empty")
	}

	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		err = os.RemoveAll(dir)
		if err != nil {
			return nil, err
		}
	}
	err := os.Mkdir(dir, 0755)
	if err != nil {
		return nil, err
	}

	opts := badger.DefaultOptions
	opts.Dir = dir
	opts.ValueDir = dir

	cache := new(badgerSet)
	cache.dir = dir
	cache.db, err = badger.Open(opts)
	if err != nil {
		return nil, err
	}

	cache.txn = cache.db.NewTransaction(true)
	cache.counter = 0
	return cache, nil

}

func (hc *badgerSet) Close() error {
	err := hc.txn.Commit()
	if err != nil {
		return err
	}
	hc.txn.Discard()
	err = hc.db.Close()

	if err != nil {
		return err
	}
	err = os.RemoveAll(hc.dir)
	if err != nil {
		return err
	}
	return nil
}

func (hc *badgerSet) Contains(key int64) (bool, error) {
	bk := int64ToByteArray(key)

	_, err := hc.txn.Get(bk)
	if err == badger.ErrKeyNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil

}

func (hc *badgerSet) Put(key int64) error {
	hc.counter++
	bk := int64ToByteArray(key)

	// Use the transaction...
	err := hc.txn.Set(bk, []byte(""))
	if err != nil {
		return err
	}

	// Commit the transaction and check for error.
	if hc.counter > 5000 {
		if err := hc.txn.Commit(); err != nil {
			return err
		}
		hc.counter = 0
		hc.txn = hc.db.NewTransaction(true)
	}
	return nil
}
