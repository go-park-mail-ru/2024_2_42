package models

type Feed struct {
	Pins []Pin `json:"pins"`
}

func NewFeed(pinList []Pin) Feed {
	return Feed{Pins: pinList}
}
