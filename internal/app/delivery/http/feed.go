package delivery

import (
	"encoding/json"
	"net/http"
	internal_errors "pinset/internal/errors"
)

func (fdc *FeedDeliveryController) Feed(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	header.Add("Access-Control-Allow-Methods", "GET")
	header.Add("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

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
