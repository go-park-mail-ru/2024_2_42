package models

import (
	"html"
	"pinset/internal/errors"
	"time"
)

type Pin struct {
	PinID        uint64     `json:"pin_id"`
	AuthorID     uint64     `json:"author_id"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	MediaUrl     string     `json:"media_url"`
	RelatedLink  string     `json:"related_link"`
	BoardID      uint64     `json:"board_id"`
	Commentaries []Comment  `json:"commentaries"`
	Bookmarks    []Bookmark `json:"bookmarks"`
	Views        uint64     `json:"views"`
	Geolocation  string     `json:"geolocation"`
	CreationTime time.Time  `json:"creation_time"`
	UpdateTime   time.Time  `json:"update_time"`
}

func (p *Pin) Sanitize() {
	p.Title = html.EscapeString(p.Title)
	p.Description = html.EscapeString(p.Description)
	p.MediaUrl = html.EscapeString(p.MediaUrl)
	p.RelatedLink = html.EscapeString(p.RelatedLink)
}

func (p Pin) Valid() error {
	if p.titleValid() && p.descriptionValid() {
		return nil
	}
	return errors.ErrPinDataInvalid
}

func (p Pin) titleValid() bool {
	return len(p.Title) > 0
}

func (p Pin) descriptionValid() bool {
	return len(p.Description) > 0
}
