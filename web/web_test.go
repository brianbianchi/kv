package web

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"kv/config"
	"kv/db"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func createShardDb(t *testing.T, idx int) *db.Database {
	t.Helper()

	name, err := ioutil.TempDir(os.TempDir(), fmt.Sprintf("db%d", idx))
	if err != nil {
		t.Fatalf("Could not create a temp db %d: %v", idx, err)
	}

	defer os.Remove(name)

	db, closeFunc, err := db.NewDatabase(name)
	if err != nil {
		t.Fatalf("Could not create new database %q: %v", name, err)
	}
	t.Cleanup(func() { closeFunc() })

	return db
}

func createShardServer(t *testing.T, idx int, addrs map[int]string) (*db.Database, *Server) {
	t.Helper()

	db := createShardDb(t, idx)

	cfg := &config.Shards{
		Addrs:  addrs,
		Count:  len(addrs),
		CurIdx: idx,
	}

	s := NewServer(db, cfg)
	return db, s
}

func TestWebServer(t *testing.T) {
	var ts1Handler func(w http.ResponseWriter, r *http.Request)
	var ts2Handler func(w http.ResponseWriter, r *http.Request)

	ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ts1Handler(w, r)
	}))
	defer ts1.Close()

	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ts2Handler(w, r)
	}))
	defer ts2.Close()

	addrs := map[int]string{
		0: strings.TrimPrefix(ts1.URL, "http://"),
		1: strings.TrimPrefix(ts2.URL, "http://"),
	}

	db1, web1 := createShardServer(t, 0, addrs)
	db2, web2 := createShardServer(t, 1, addrs)

	keys := map[string]int{
		"denver": 1,
		"SD":     0,
	}

	ts1Handler = web1.RouteHandler
	ts2Handler = web2.RouteHandler

	client := &http.Client{}
	for key := range keys {
		req, err := http.NewRequest(
			http.MethodPut,
			fmt.Sprintf(ts1.URL+"/%s", key),
			bytes.NewBuffer([]byte(fmt.Sprintf("value-%s", key))),
		)
		if err != nil {
			t.Fatalf("Could not set the key %q: %v", key, err)
		}
		_, err = client.Do(req)
		if err != nil {
			t.Fatalf("Could not set the key %q: %v", key, err)
		}
	}

	for key := range keys {
		resp, err := http.Get(fmt.Sprintf(ts1.URL+"/%s", key))
		if err != nil {
			t.Fatalf("Get key %q error: %v", key, err)
		}
		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Could read contents of the key %q: %v", key, err)
		}

		want := []byte("value-" + key)
		if !bytes.Contains(contents, want) {
			t.Errorf("Unexpected contents of the key %q: got %q, want the result to contain %q", key, contents, want)
		}

		log.Printf("Contents of key %q: %s", key, contents)
	}

	value1, err := db2.GetKey([]byte("denver"))
	if err != nil {
		t.Fatalf("Denver key error: %v", err)
	}

	want1 := "value-denver"
	if !bytes.Equal(value1, []byte(want1)) {
		t.Errorf("Unexpected value of  key: got %q, want %q", value1, want1)
	}

	value2, err := db1.GetKey([]byte("SD"))
	if err != nil {
		t.Fatalf("SD key error: %v", err)
	}

	want2 := "value-SD"
	if !bytes.Equal(value2, []byte(want2)) {
		t.Errorf("Unexpected value of SD key: got %q, want %q", value2, want2)
	}
}
