package models

import "time"

type Pin struct {
	PinID        uint64    `json:"pin_id"`
	AuthorID     uint64    `json:"author_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	MediaUrl     string    `json:"media_url"`
	RelatedLink  string    `json:"related_link"`
	BoardID      uint64    `json:"board_id"`
	Commentaries []Comment `json:"commentaries"`
	Likes        []Like    `json:"likes"`
	CreationTime time.Time `json:"creation_time"`
}

func (p Pin) TitleValid() bool {
	return len(p.Title) > 0
}

func (p Pin) DescriptionValid() bool {
	return len(p.Description) > 0
}
