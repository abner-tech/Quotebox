package models

import (
	"errors"
	"time"
)

var ErrRecordNotFound = errors.New("models: no matching record found")

// A struct to hold a quote
type Quote struct {
	Quotations_id  int
	Insertion_date time.Time
	Author_name    string
	Category       string
	Quote          string
}
