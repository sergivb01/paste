package server

import (
	"time"
	"github.com/sergivb01/paste/utils"
)

type Paste struct {
	ID string `json:"id" db:"id"`

	Title string `json:"title" db:"title"`
	Content string `json:"content" db:"content"`

	CreatedAt time.Time `json:"createdAt" db:"created"`
	ExpiresAt utils.NullTime `json:"expiresAt" db:"expires"`
}
