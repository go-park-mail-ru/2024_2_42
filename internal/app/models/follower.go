package models

import "time"

type Follower struct {
	OwnerID      uint64    `json:"owner_id"`
	FollowerID   uint64    `json:"follower_id"`
	CreationTime time.Time `json:"creation_time"`
}
