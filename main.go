package main

import (
	"flag"
	"kv/config"
	"kv/db"
	"kv/web"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var (
	dbpath     = flag.String("dbpath", "./tmp/data", "Path to leveldb.")
	serverPort = flag.String("addr", "127.0.0.1:8080", "Address to serve API.")
	configFile = flag.String("config", "sharding.toml", "Config file for static sharding.")
	shard      = flag.String("shard", "", "The name of the shard.")
)

func parseFlags() {
	flag.Parse()

	if *shard == "" {
		log.Fatalf("Must provide shard name.")
	}
}

func main() {
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 100
	rand.Seed(time.Now().Unix())

	parseFlags()

	c, err := config.ParseFile(*configFile)
	if err != nil {
		log.Fatalf("Error parsing config %q: %v", *configFile, err)
	}

	shards, err := config.ParseShards(c.Shards, *shard)
	if err != nil {
		log.Fatalf("Error parsing shards config: %v", err)
	}

	log.Printf("Shard count is %d, current shard: %d", shards.Count, shards.CurIdx)

	db, close, err := db.NewDatabase(*dbpath)
	if err != nil {
		log.Fatalf("Error creating %q: %v", *dbpath, err)
	}
	defer close()

	srv := web.NewServer(db, shards)

	http.HandleFunc("/", srv.RouteHandler)
	log.Fatal(http.ListenAndServe(*serverPort, nil))
}
