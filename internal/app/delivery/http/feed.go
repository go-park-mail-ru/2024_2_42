package delivery

import (
	"encoding/json"
	"net/http"
	internal_errors "pinset/internal/errors"
)

func (fdc *FeedDeliveryController) Feed(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.Write([]byte("For now only GET method is allowed"))
		return
	}

	feed := fdc.Usecase.Feed()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(feed); err != nil {
		internal_errors.SendErrorResponse(w, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrFeedNotAccessible,
		})
		return
	}
}
