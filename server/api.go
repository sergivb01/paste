package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/teris-io/shortid"
)

func (s *Server) handleAPITestExpired() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ids, err := s.deleteExpiredPosts()
		if err != nil {
			s.handleError(err, http.StatusInternalServerError)(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(ids); err != nil {
			s.handleError(err, http.StatusInternalServerError)(w, r)
		}
	}
}

func (s *Server) handleAPIPasteGET() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rawID, ok := mux.Vars(r)["id"]
		if !ok {
			s.handleError(errNoID, http.StatusBadRequest)(w, r)
			return
		}

		paste, err := s.getPaste(rawID)
		if err != nil {
			code := http.StatusBadRequest
			if err == errPasteNotFoundOrExpired {
				code = http.StatusNotFound
			}

			s.handleError(err, code)(w, r)
			return
		}
		fmt.Printf("found paste: %+v\n", paste)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(paste); err != nil {
			s.handleError(err, http.StatusInternalServerError)(w, r)
		}
	}
}

func (s *Server) handleAPIPastePOST() http.HandlerFunc {
	type rawPaste struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Expires string `json:"expires"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		raw := &rawPaste{}

		if err := json.NewDecoder(r.Body).Decode(raw); err != nil {
			s.handleError(err, http.StatusBadRequest)(w, r)
			return
		}

		id, err := shortid.Generate()
		if err != nil {
			s.handleError(err, http.StatusInternalServerError)(w, r)
			return
		}

		paste := &Paste{
			ID:        id,
			Title:     raw.Title,
			Content:   raw.Content,
			CreatedAt: time.Now(),
		}

		if raw.Expires != "" {
			d, err := time.ParseDuration(raw.Expires)
			if err != nil {
				s.handleError(err, http.StatusBadRequest)(w, r)
				return
			}

			paste.ExpiresAt.Time = time.Now().Add(d)
			paste.ExpiresAt.Valid = true
		}

		if err := s.savePaste(paste); err != nil {
			s.handleError(err, http.StatusInternalServerError)(w, r)
			return
		}
		fmt.Printf("created paste with ID %s: %+v\n", paste.ID, paste)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(paste); err != nil {
			s.handleError(err, http.StatusInternalServerError)(w, r)
		}
	}
}
