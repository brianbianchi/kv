package web

import (
	"fmt"
	"hash/fnv"
	"io"
	"io/ioutil"
	"kv/db"
	"net/http"
	"strings"
)

type Server struct {
	db         *db.Database
	shardIdx   int
	shardCount int
	addrs      map[int]string
}

func NewServer(db *db.Database, shardIdx, shardCount int, addrs map[int]string) *Server {
	return &Server{
		db:         db,
		shardIdx:   shardIdx,
		shardCount: shardCount,
		addrs:      addrs,
	}
}

func (s *Server) getShard(key string) int {
	h := fnv.New64()
	h.Write([]byte(key))
	return int(h.Sum64() % uint64(s.shardCount))
}

func (s *Server) redirect(shard int, w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("http://%s%s", s.addrs[shard], r.RequestURI)
	fmt.Printf("redirecting from shard %d to shard %d (%q)\n", s.shardIdx, shard, url)

	resp, err := http.Get(url)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error redirecting the request: %v", err)
		return
	}
	defer resp.Body.Close()

	io.Copy(w, resp.Body)
}

// Handles API requests
func (s *Server) RouteHandler(w http.ResponseWriter, r *http.Request) {
	key := []byte(strings.TrimPrefix(r.URL.Path, "/"))

	if len(key) == 0 {
		w.Header().Set("Content-Length", "0")
		w.WriteHeader(404)
		return
	}

	switch r.Method {
	case "GET":
		shard := s.getShard(string(key))
		data, err := s.db.GetKey(key)

		if shard != s.shardIdx {
			s.redirect(shard, w, r)
			return
		}

		fmt.Printf("Shard = %d, current shard = %d, addr = %q, Value = %q, error = %v\n", shard, s.shardIdx, s.addrs[shard], data, err)
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

		shard := s.getShard(string(key))
		if shard != s.shardIdx {
			s.redirect(shard, w, r)
			return
		}
		err = s.db.PutKey(key, value)
		fmt.Printf("Error = %v, shardIdx = %d, current shard = %d\n", err, shard, s.shardIdx)
		if err != nil {
			w.Header().Set("Content-Length", "0")
			w.WriteHeader(404)
			return
		}
		w.WriteHeader(201)
	case "DELETE":
		err := s.db.DeleteKey(key)
		if err != nil {
			w.Header().Set("Content-Length", "0")
			w.WriteHeader(404)
			return
		}
		w.WriteHeader(204)
	}
}
