package badger

import (
	"os"

	"github.com/czsilence/EvernoteTransfer/erro"
	"github.com/czsilence/EvernoteTransfer/storage"
	"github.com/czsilence/go/app"
	"github.com/czsilence/go/log"
	"github.com/czsilence/go/timer"
	"github.com/dgraph-io/badger"
)

func DB() storage.DB {
	if db != nil && db.db != nil {
		return db
	}

	var storage_path = app.GetDataPath("storage")
	if err := os.MkdirAll(storage_path, 0755); err != nil {
		log.E2("[db] create db path failed, err: %v", err)
	}

	opts := badger.DefaultOptions
	opts.Dir = storage_path
	opts.ValueDir = storage_path

	badger_db, err := badger.Open(opts)
	if err != nil {
		log.E2("[db] init badger db failed, err: %v", err)
	}

	db = &_BadgerDB{
		db: badger_db,
	}

	// timer to gc
	timer.SetInterval(func() {
		db.gc()
	}, 2000, false)

	return db
}

//////////////////////////////////////////////////////
type _BadgerDB struct {
	db *badger.DB
}

func (db *_BadgerDB) Put(key string, data []byte) erro.Error {
	err := db.db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(key), data)
		return err
	})

	if err != nil {
		return erro.E_Storage_PutFailed.With(err)
	}
	return nil
}

func (db *_BadgerDB) Get(key string) (data []byte, exist bool) {
	if db == nil || db.db == nil {
		return
	}

	db.db.View(func(txn *badger.Txn) error {
		if item, err := txn.Get([]byte(key)); err != nil {
			return err
		} else if data, err = item.Value(); err != nil {
			return err
		} else {
			exist = data != nil
		}
		return nil
	})

	return
}

func (db *_BadgerDB) Del(key string) erro.Error {
	err := db.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(key))
		return err
	})

	if err != nil {
		return erro.E_Storage_DelFailed.With(err)
	}
	return nil
}

func (db *_BadgerDB) Close() {
	if db != nil && db.db != nil {
		if err := db.db.Close(); err != nil {
			log.D("[db] badger db close with err: %v", err)
		} else {
			log.D("[db] badger db closed")
		}
		db.db = nil
	}
}

func (db *_BadgerDB) gc() {
	if db != nil && db.db != nil {
		db.db.RunValueLogGC(0.7)
	}
}

//////////////////////////////////////////////////////
var db *_BadgerDB
