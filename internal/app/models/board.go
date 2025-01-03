package models

import (
	"html"
	"time"

	"pinset/internal/errors"
)

type Board struct {
	BoardID       uint64    `json:"board_id"`
	OwnerID       uint64    `json:"owner_id"`
	Cover         string    `json:"board_cover"`
	Name          string    `json:"board_name"`
	Description   string    `json:"board_description"`
	Public        bool      `json:"public"`
	Pins          []Pin     `json:"pins"`
	Sections      []Section `json:"sections"`
	Collaborators []User    `json:"collaborators"`
	CreationTime  time.Time `json:"creation_time"`
	UpdateTime    time.Time `json:"update_time"`
}

type BoardPin struct {
	BoardID uint64 `json:"board_id"`
	Cover   string `json:"board_cover"`
	Name    string `json:"board_name"`
	Public  bool   `json:"public"`
}

func (b *Board) Sanitize() {
	b.Name = html.EscapeString(b.Name)
	b.Description = html.EscapeString(b.Description)
}

func (b Board) Valid() error {
	if b.nameValid() && b.descriptionValid() {
		return nil
	}
	return errors.ErrBoardDataInvalid
}

func (b Board) nameValid() bool {
	return len(b.Name) > 0
}

func (b Board) descriptionValid() bool {
	return len(b.Description) > 0
}
