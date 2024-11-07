package models

import "time"

type Bookmark struct {
	BookmarkID   uint64    `json:"bookmark_id"`
	OwnerID      uint64    `json:"owner_id"`
	PinID        uint64    `json:"pin_id"`
	BookmarkTime time.Time `json:"bookmark_time"`
}
