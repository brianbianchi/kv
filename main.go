package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

type App struct {
	db *leveldb.DB
}

func (a *App) RouteHandler(w http.ResponseWriter, r *http.Request) {
	key := []byte(strings.TrimPrefix(r.URL.Path, "/"))

	if len(key) == 0 {
		w.Header().Set("Content-Length", "0")
		w.WriteHeader(404)
	}

	switch r.Method {
	case "GET":
		data, err := a.db.Get(key, nil)
		if err != nil {
			w.Header().Set("Content-Length", "0")
			w.WriteHeader(404)
			return
		}
		w.WriteHeader(201)
		w.Write(data)
	case "PUT":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.Header().Set("Content-Length", "0")
			w.WriteHeader(404)
			return
		}
		value := []byte(body)

		err = a.db.Put(key, value, nil)
		if err != nil {
			w.Header().Set("Content-Length", "0")
			w.WriteHeader(404)
			return
		}
		w.WriteHeader(201)
	case "DELETE":
		err := a.db.Delete(key, nil)
		if err != nil {
			w.Header().Set("Content-Length", "0")
			w.WriteHeader(404)
			return
		}
		w.WriteHeader(204)
	}
}

func main() {
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 100
	rand.Seed(time.Now().Unix())

	dbpath := flag.String("dbpath", "./tmp/data", "Path to leveldb.")
	serverPort := flag.String("port", ":3000", "Port to serve API.")
	flag.Parse()

	db, err := leveldb.OpenFile(*dbpath, nil)
	if err != nil {
		panic(fmt.Sprintf("LevelDB open failed: %s", err))
	}
	defer db.Close()

	a := &App{
		db: db,
	}

	http.HandleFunc("/", a.RouteHandler)
	log.Fatal(http.ListenAndServe(*serverPort, nil))
}
