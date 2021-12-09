package web

import (
	"io/ioutil"
	"kv/db"
	"net/http"
	"strings"
)

type Server struct {
	db *db.Database
}

func NewServer(db *db.Database) *Server {
	return &Server{
		db: db,
	}
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
		data, err := s.db.GetKey(key)
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

		err = s.db.PutKey(key, value)
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
