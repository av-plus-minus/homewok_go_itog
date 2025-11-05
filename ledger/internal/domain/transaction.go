package domain

import (
	"errors"
	"time"
)

type Transaction struct {
	ID          int
	Amount      float64
	Category    string
	Description string
	Date        time.Time
}

func (t Transaction) Validate() error {
	if t.Amount <= 0 {
		return errors.New("amount must be > 0")
	}
	if t.Category == "" {
		return errors.New("category must be non-empty")
	}
	return nil
}

func (t *Transaction) SetDefaultDate() {
	if t.Date.IsZero() {
		t.Date = time.Now()
	}
}
