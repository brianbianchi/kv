package db

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func createTempDb(t *testing.T) *Database {
	t.Helper()

	name, err := ioutil.TempDir(os.TempDir(), "kvdb")
	if err != nil {
		t.Fatalf("Could not create temp file: %v", err)
	}

	defer os.Remove(name)

	db, closeFunc, err := NewDatabase(name)
	if err != nil {
		t.Fatalf("Could not create a new database: %v", err)
	}
	t.Cleanup(func() { closeFunc() })

	return db
}

func TestGetSet(t *testing.T) {
	db := createTempDb(t)

	if err := db.PutKey([]byte("abc"), []byte("123")); err != nil {
		t.Fatalf("Could not write key: %v", err)
	}

	value, err := db.GetKey([]byte("abc"))
	if err != nil {
		t.Fatalf(`Could not get the key "abc": %v`, err)
	}

	if !bytes.Equal(value, []byte("123")) {
		t.Errorf(`Unexpected value: got %q, want %q`, value, "Great")
	}

	if err := db.DeleteKey([]byte("abc")); err != nil {
		t.Fatalf("Could not delete key: %v", err)
	}

	_, err = db.GetKey([]byte("abc"))
	if err == nil {
		t.Fatalf(`Got the key "abc". Delete failed`)
	}
}
