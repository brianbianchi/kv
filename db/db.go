package db

import (
	"github.com/syndtr/goleveldb/leveldb"
)

type Database struct {
	db *leveldb.DB
}

// Returns an instance of a database that we can work with.
func NewDatabase(dbPath string) (db *Database, closeFunc func() error, err error) {
	leveldb, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return nil, nil, err
	}

	d := &Database{db: leveldb}
	closeFunc = leveldb.Close

	return d, closeFunc, nil
}

// Gets the value from the given key
func (d *Database) GetKey(key []byte) ([]byte, error) {
	data, err := d.db.Get([]byte(key), nil)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Creates or replaces a key-value pair.
func (d *Database) PutKey(key []byte, value []byte) error {
	err := d.db.Put([]byte(key), []byte(value), nil)
	if err != nil {
		return err
	}

	return nil
}

// Deletes the given key-value pair.
func (d *Database) DeleteKey(key []byte) error {
	err := d.db.Delete([]byte(key), nil)
	if err != nil {
		return err
	}

	return nil
}
