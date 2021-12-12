package web

import (
	"fmt"
	"io"
	"io/ioutil"
	"kv/config"
	"kv/db"
	"net/http"
	"strings"
)

type Server struct {
	db     *db.Database
	shards *config.Shards
}

func NewServer(db *db.Database, s *config.Shards) *Server {
	return &Server{
		db:     db,
		shards: s,
	}
}

func (s *Server) redirect(shard int, w http.ResponseWriter, r *http.Request) {
	url := "http://" + s.shards.Addrs[shard] + r.RequestURI
	fmt.Printf("redirecting from shard %d to shard %d (%q)\n", s.shards.CurIdx, shard, url)

	client := &http.Client{}
	var req *http.Request
	var resp *http.Response
	var err error = nil
	switch r.Method {
	case "GET":
		req, err = http.NewRequest(http.MethodGet, url, nil)
	case "PUT":
		req, err = http.NewRequest(http.MethodPut, url, r.Body)
	case "DELETE":
		req, err = http.NewRequest(http.MethodDelete, url, nil)
	}
	if err != nil {
		w.WriteHeader(500)
		fmt.Printf("Error redirecting the request: %v", err)
		return
	}
	resp, err = client.Do(req)
	if err != nil {
		w.WriteHeader(500)
		fmt.Printf("Error redirecting the request: %v", err)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
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
		shard := s.shards.Index(string(key))
		if shard != s.shards.CurIdx {
			s.redirect(shard, w, r)
			return
		}

		data, err := s.db.GetKey(key)

		fmt.Printf("Shard = %d, current shard = %d, addr = %q, Value = %q, error = %v\n", shard, s.shards.CurIdx, s.shards.Addrs[shard], data, err)
		if err != nil {
			w.Header().Set("Content-Length", "0")
			w.WriteHeader(404)
			return
		}
		w.Write(data)
	case "PUT":
		shard := s.shards.Index(string(key))
		if shard != s.shards.CurIdx {
			s.redirect(shard, w, r)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.Header().Set("Content-Length", "0")
			w.WriteHeader(404)
			return
		}

		err = s.db.PutKey(key, []byte(body))
		fmt.Printf("Shard = %d, current shard = %d, addr = %q, Value = %q, error = %v\n", shard, s.shards.CurIdx, s.shards.Addrs[shard], []byte(body), err)
		if err != nil {
			w.Header().Set("Content-Length", "0")
			w.WriteHeader(404)
			return
		}
		w.WriteHeader(201)
	case "DELETE":
		shard := s.shards.Index(string(key))
		if shard != s.shards.CurIdx {
			s.redirect(shard, w, r)
			return
		}
		err := s.db.DeleteKey(key)
		fmt.Printf("Shard = %d, current shard = %d, addr = %q, error = %v\n", shard, s.shards.CurIdx, s.shards.Addrs[shard], err)

		if err != nil {
			w.Header().Set("Content-Length", "0")
			w.WriteHeader(404)
			return
		}
		w.WriteHeader(204)
	}
}
