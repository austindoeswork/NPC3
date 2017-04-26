package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
	port string
	mux  *http.ServeMux
}

func (s *Server) AddHandler(path string, handler http.Handler) {
	s.mux.Handle(path, handler)
}

func (s *Server) AddStaticHandler(dirpath string) {
	s.mux.Handle("/", http.FileServer(http.Dir(dirpath)))
}

func New(port string) *Server {
	s := &Server{
		port: port,
		mux:  http.NewServeMux(),
	}
	s.mux.HandleFunc("/dawg", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "yo dawg")
	})
	return s
}

func (s *Server) Start() error {
	return http.ListenAndServe(s.port, s.mux)
}
