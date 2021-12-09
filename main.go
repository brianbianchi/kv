package main

import (
	"flag"
	"kv/db"
	"kv/web"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var (
	dbpath     = flag.String("dbpath", "./tmp/data", "Path to leveldb.")
	serverPort = flag.String("port", ":3000", "Port to serve API.")
)

func parseFlags() {
	flag.Parse()
}

func main() {
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 100
	rand.Seed(time.Now().Unix())

	parseFlags()

	db, close, err := db.NewDatabase(*dbpath)
	if err != nil {
		log.Fatalf("NewDatabase(%q): %v", *dbpath, err)
	}
	defer close()

	srv := web.NewServer(db)

	http.HandleFunc("/", srv.RouteHandler)
	log.Fatal(http.ListenAndServe(*serverPort, nil))
}
