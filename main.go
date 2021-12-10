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

	"github.com/BurntSushi/toml"
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

	var c config.Config
	if _, err := toml.DecodeFile(*configFile, &c); err != nil {
		log.Fatalf("toml.DecodeFile(%q): %v", *configFile, err)
	}

	var shardCount int
	var shardIdx int = -1
	var addrs = make(map[int]string)

	for _, s := range c.Shards {
		addrs[s.Idx] = s.Address

		if s.Idx+1 > shardCount {
			shardCount = s.Idx + 1
		}
		if s.Name == *shard {
			shardIdx = s.Idx
		}
	}

	if shardIdx < 0 {
		log.Fatalf("Shard %q was not found", *shard)
	}
	log.Printf("Shard count is %d, current shard: %d", shardCount, shardIdx)

	db, close, err := db.NewDatabase(*dbpath)
	if err != nil {
		log.Fatalf("NewDatabase(%q): %v", *dbpath, err)
	}
	defer close()

	srv := web.NewServer(db, shardIdx, shardCount, addrs)

	http.HandleFunc("/", srv.RouteHandler)
	log.Fatal(http.ListenAndServe(*serverPort, nil))
}
