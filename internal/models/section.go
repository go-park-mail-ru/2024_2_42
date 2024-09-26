package models

import "time"

type Section struct {
	SectionID    uint64    `json:"section_id"`
	BoardID      uint64    `json:"board_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Pins         []Pin     `json:"pins"`
	CreationTime time.Time `json:"CreationTime"`
}
