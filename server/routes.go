package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/teris-io/shortid"

	"github.com/sergivb01/paste/utils"
)

var (
	errNoID        = errors.New("missing ID parameter in URL")
	errInvalidForm = errors.New("invalid form. please check inputs")
	errInvalidTime = errors.New("invalid time from form select")

	customDurations = map[string]time.Duration{
		"Never":     0,
		"5 minutes": time.Minute * 5,
		"1 hour": time.Hour,
		"6 hours": time.Hour * 6,
		"1 day": time.Hour * 24,
		"1 week": time.Hour * 24 * 7,
		"2 weeks": time.Hour * 24 * 14,
		"1 month": time.Hour * 24 * 30,
	}
)

func (s *Server) routes() {
	r := mux.NewRouter()

	r.HandleFunc("/", s.handleIndex()).Methods("GET")
	r.HandleFunc("/", s.handlePastePOST()).Methods("POST")
	r.HandleFunc("/latests", s.handleLatests()).Methods("GET")

	r.HandleFunc("/{id}", s.handlePasteGET()).Methods("GET")
	r.HandleFunc("/api/paste", s.handleAPIPasteGET()).Methods("GET")
	r.HandleFunc("/api/paste", s.handleAPIPastePOST()).Methods("POST")

	s.router = r
}

func (s *Server) renderTpl(w http.ResponseWriter, name string, args interface{}) {
	if err := s.tpl.ExecuteTemplate(w, name, args); err != nil {
		s.handleError(err, http.StatusInternalServerError)
	}
}

func (s *Server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pastes, err := s.getLatestPastes(5)
		if err != nil {
			s.handleError(err, http.StatusInternalServerError)(w, r)
			return
		}

		s.renderTpl(w, "index", pastes)
	}
}

func (s *Server) handlePasteGET() http.HandlerFunc {
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

		s.renderTpl(w, "paste", paste)
	}
}


func (s *Server) handleLatests() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pastes, err := s.getLatestPastes(25)
		if err != nil {
			s.handleError(err, http.StatusInternalServerError)(w, r)
			return
		}

		s.renderTpl(w, "latests", pastes)
	}
}

func (s *Server) handlePastePOST() http.HandlerFunc {
	type rawPaste struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Expires string `json:"expires"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			s.handleError(err, http.StatusBadRequest)(w, r)
			return
		}

		id, err := shortid.Generate()
		if err != nil {
			s.handleError(err, http.StatusInternalServerError)(w, r)
			return
		}

		var (
			title      = r.FormValue("title")
			content    = r.FormValue("content")
			rawExpires = r.FormValue("expires")
		)

		if len(title) >= 25 || len(title) == 0 || len(content) == 0 || len(content) >= 1024 || len(rawExpires) == 0 {
			s.handleError(errInvalidForm, http.StatusBadRequest)(w, r)
			return
		}

		d, ok := customDurations[rawExpires]
		if !ok {
			s.handleError(errInvalidTime, http.StatusBadRequest)(w, r)
			return
		}

		paste := &Paste{
			ID:        id,
			Title:     title,
			Content:   content,
			CreatedAt: time.Now(),
			ExpiresAt: utils.NullTime{
				Time: time.Now().Add(d),
				Valid: d != 0,
			},
		}

		// if rawExpires != "" {
		// 	t, err := time.Parse(time.RFC3339[:16], rawExpires)
		// 	if err != nil {
		// 		s.handleError(err, http.StatusBadRequest)(w, r)
		// 		return
		// 	}

		// 	paste.ExpiresAt.Time = t
		// 	paste.ExpiresAt.Valid = true
		// }

		if err := s.savePaste(paste); err != nil {
			s.handleError(err, http.StatusInternalServerError)(w, r)
			return
		}

		http.Redirect(w, r, "/"+paste.ID, http.StatusSeeOther)
	}
}
