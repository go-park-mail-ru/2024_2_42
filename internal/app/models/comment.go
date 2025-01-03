package models

import (
	"html"
	"pinset/internal/errors"
	"time"
)

type Comment struct {
	CommentID    uint64     `json:"comment_id"`
	PinID        uint64     `json:"pin_id"`
	AuthorID     uint64     `json:"author_id"`
	Body         string     `json:"body"`
	CreationTime time.Time  `json:"creation_time"`
	UpdateTime   time.Time  `json:"update_time"`
}

func (c *Comment) Sanitize() {
	c.Body = html.EscapeString(c.Body)
}

func (c Comment) Valid() error {
	if len(c.Body) > 0 {
		return nil
	}
	return errors.ErrCommentDataInvalid
}
