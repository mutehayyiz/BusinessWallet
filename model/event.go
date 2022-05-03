package model

import (
	"errors"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"html"
	"strings"
	"time"
)

type Event struct {
	gorm.Model
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Organizer   int           `json:"organizer"`
	Attendees   pq.Int64Array `json:"attendees" gorm:"type:integer[]"`
	Start       time.Time     `json:"start"`
	Finish      time.Time     `json:"finish"`
}

type CreateEventRequest struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Start       time.Time `json:"start"`
	Finish      time.Time `json:"finish"`
}

func (e CreateEventRequest) Validate() error {
	e.Name = html.EscapeString(strings.TrimSpace(e.Name))
	e.Description = html.EscapeString(strings.TrimSpace(e.Description))

	if e.Name == "" {
		return errors.New("name is required")
	}

	if e.Description == "" {
		return errors.New("description is required")
	}

	return nil
}
