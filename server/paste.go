package server

import (
	"github.com/sergivb01/paste/utils"
	"time"
)

type Paste struct {
	ID string `json:"id" db:"id"`

	Title   string `json:"title" db:"title"`
	Content string `json:"content" db:"content"`

	CreatedAt time.Time      `json:"createdAt" db:"created"`
	ExpiresAt utils.NullTime `json:"expiresAt" db:"expires"`
}

func (p Paste) ExpiresIn() time.Duration {
	if !p.ExpiresAt.Valid {
		return 0
	}

	return p.ExpiresAt.Time.Sub(p.CreatedAt)
}
