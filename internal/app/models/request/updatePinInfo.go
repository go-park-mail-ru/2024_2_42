package request

type UpdatePinRequest struct {
	PinID       uint64 `json:"pin_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	BoardID     uint64 `json:"board_id"`
	RelatedLink string `json:"related_link"`
	Geolocation string `json:"geolocation"`
}

func (upr UpdatePinRequest) Valid() bool {
	return len(upr.Title) > 0 && upr.BoardID > 0
}
