package model

import (
	"time"
)

type URL struct {
	Id          int
	Code        string
	OriginalURL string
	CreatedAt   time.Time
}
