package models

import "time"

type Like struct {
	LikeID   uint64    `json:"like_id"`
	OwnerID  uint64    `json:"owner_id"`
	LikeTime time.Time `json:"like_time"`
}
