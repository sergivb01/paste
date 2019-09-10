package server

import (
	"time"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/teris-io/shortid"
)

func (s *Server) routes() {
	r := mux.NewRouter()

	r.HandleFunc("/", s.handleIndex()).Methods("GET")
	r.HandleFunc("/{id}", s.handleGETPaste()).Methods("GET")
	r.HandleFunc("/paste", s.handlePostPaste()).Methods("POST")

	s.router = r
}

func (s *Server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pastes, err := s.getLatestPastes()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(pastes); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (s *Server) handleGETPaste() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rawID, ok := mux.Vars(r)["id"]
		if !ok {
			http.Error(w, "No id provided", http.StatusBadRequest)
			return
		}

		paste, err := s.getPaste(rawID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(paste); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (s *Server) handlePostPaste() http.HandlerFunc {
	type rawPaste struct {
		Title string `json:"title"`
		Content string `json:"content"`
		Expires string `json:"expires"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		raw := &rawPaste{}

		if err := json.NewDecoder(r.Body).Decode(raw); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := shortid.Generate()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		paste := &Paste{
			ID: id,
			Title: raw.Title,
			Content: raw.Content,
			CreatedAt: time.Now(),
		}

		if raw.Expires != "" {
			d, err := time.ParseDuration(raw.Expires)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			paste.ExpiresAt.Time = time.Now().Add(d)
			paste.ExpiresAt.Valid = true
		}

		if err := s.savePaste(paste); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(paste); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
