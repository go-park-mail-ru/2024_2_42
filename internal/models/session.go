package models

import "time"

type Session struct {
	SessionID    uint64                 `json:"session_id"`
	Values       map[string]interface{} `json:"values"`
	CreationTime time.Time              `json:"creation_time"`
}
