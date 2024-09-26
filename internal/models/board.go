package models

import "time"

type Board struct {
	BoardID       uint64    `json:"board_id"`
	OwnerID       uint64    `json:"owner_id"`
	Name          string    `json:"board_name"`
	Description   string    `json:"board_description"`
	Public        bool      `json:"public"`
	Pins          []Pin     `json:"pins"`
	Sections      []Section `json:"sections"`
	Collaborators []User    `json:"collaborators"`
	CreationTime  time.Time `json:"creation_time"`
}
