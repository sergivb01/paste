package server

import (
	"database/sql"
	"errors"
)

var errPasteNotFoundOrExpired = errors.New("paste not found or has expired")

func (s *Server) savePaste(paste *Paste) error {
	_, err := s.db.NamedExec("INSERT INTO pastes (id, title, content, created, expires) VALUES (:id, :title, :content, :created, :expires)", paste)
	return err
}

func (s *Server) getPaste(ID string) (*Paste, error) {
	paste := &Paste{}

	if err := s.db.QueryRowx("SELECT * FROM pastes WHERE id = $1 AND (expires IS NULL OR expires > now())", ID).StructScan(paste); err != nil {
		if err == sql.ErrNoRows {
			return nil, errPasteNotFoundOrExpired
		}
		return nil, err
	}

	return paste, nil
}

func (s *Server) getLatestPastes(max int) ([]Paste, error) {
	var pastes []Paste

	if max < 1 {
		max = 1
	} else if max >= 25 {
		max = 25
	}

	return pastes, s.db.Select(&pastes, "SELECT * FROM pastes WHERE expires IS NULL OR expires > now() ORDER BY created DESC LIMIT $1", max)
}

func (s *Server) deleteExpiredPosts() ([]string, error) {
	var ids []string

	if err := s.db.Select(&ids, "DELETE FROM pastes WHERE expires < now() RETURNING id"); err != nil {
		return nil, err
	}

	return ids, nil
}
