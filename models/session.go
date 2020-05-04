package models

import "time"

type session struct {
	Username     string
	LastActivity time.Time
}
