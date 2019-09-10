package utils

import (
	"database/sql/driver"
	"time"
)

type NullTime struct {
	Time  time.Time
	Valid bool `json:"-"`
}

func (n *NullTime) Scan(value interface{}) error {
	n.Time, n.Valid = value.(time.Time)
	return nil
}

func (n NullTime) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Time, nil
}
