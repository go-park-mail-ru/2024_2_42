package models

import (
	"html"
	"time"
	"youpin/internal/errors"
)

type Section struct {
	SectionID    uint64    `json:"section_id"`
	BoardID      uint64    `json:"board_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Pins         []Pin     `json:"pins"`
	CreationTime time.Time `json:"creation_time"`
	UpdateTime   time.Time `json:"update_time"`
}

func (s *Section) Sanitize() {
	s.Name = html.EscapeString(s.Name)
	s.Description = html.EscapeString(s.Description)
}

func (s Section) Valid() error {
	if s.nameValid() && s.descriptionValid() {
		return nil
	}
	return errors.ErrSectionDataInvalid
}

func (s Section) nameValid() bool {
	return len(s.Name) > 0
}

func (s Section) descriptionValid() bool {
	return len(s.Description) > 0
}
