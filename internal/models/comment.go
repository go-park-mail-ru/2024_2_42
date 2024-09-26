package models

import "time"

type Comment struct {
	CommentID   uint64    `json:"comment_id"`
	PinID       uint64    `json:"pin_id"`
	AuthorID    uint64    `json:"author_id"`
	Body        string    `json:"body"`
	Likes       []Like    `json:"likes"`
	PublishTime time.Time `json:"publish_time"`
}

func (c Comment) BodyValid() bool {
	return len(c.Body) > 0
}
