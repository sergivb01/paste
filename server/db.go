package server

func (s *Server) savePaste(paste *Paste) error {
	_, err := s.db.NamedExec("INSERT INTO pastes (id, title, content, created, expires) VALUES (:id, :title, :content, :created, :expires)", paste)
	return err
}

func (s *Server) getPaste(ID string) (*Paste, error) {
	paste := &Paste{}

	r := s.db.QueryRowx("SELECT * FROM pastes WHERE id = $1 AND (expires IS NULL OR expires > now())", ID)
	if err := r.StructScan(paste); err != nil {
		return nil, err
	}

	return paste, nil
}

func (s *Server) getLatestPastes() ([]Paste, error) {
	var pastes []Paste

	if err := s.db.Select(&pastes, "SELECT * FROM pastes WHERE expires IS NULL OR expires > now() ORDER BY created DESC LIMIT 15"); err != nil {
		return nil, err
	}

	return pastes, nil
}