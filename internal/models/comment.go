package models

import (
	"html"
	"time"
	"youpin/internal/errors"
)

type Comment struct {
	CommentID    uint64    `json:"comment_id"`
	PinID        uint64    `json:"pin_id"`
	AuthorID     uint64    `json:"author_id"`
	Body         string    `json:"body"`
	Likes        []Like    `json:"likes"`
	CreationTime time.Time `json:"creation_time"`
	UpdateTime   time.Time `json:"update_time"`
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
