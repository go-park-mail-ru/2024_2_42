package request

type UpdateBoardRequest struct {
	BoardID     uint64 `jsno:"board_id"`
	Cover       string `json:"board_cover"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
	PinAsACover bool   `json:"pin_as_a_cover"`
}

func (ubr UpdateBoardRequest) Valid() bool {
	return len(ubr.Title) > 0
}
