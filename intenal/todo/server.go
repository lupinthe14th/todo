package todo

import (
	"context"
	"encoding/json"
	"net/http"
)

type Server struct {
	server *http.Server
	db     DB
}

func NewServer(addr string, db DB) *Server {
	return &Server{
		server: &http.Server{Addr: addr},
		db:     db,
	}
}

func (s *Server) Start() error {
	s.initHandlers()
	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.server.Shutdown(ctx)
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) initHandlers() {
	mux := http.NewServeMux()
	s.server.Handler = mux

	mux.HandleFunc("/create", s.HandleCreate)
	mux.HandleFunc("/getall", s.HandleGetAll)
}

func (s *Server) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var todo TODO
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.db.Put(r.Context(), &todo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func (s *Server) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	todos, err := s.db.GetAll(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
