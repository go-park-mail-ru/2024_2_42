package handlers

import (
	"encoding/json"
	"net/http"
	internal_errors "youpin/internal/errors"
	"youpin/internal/models"
)

var (
	pins = []models.Pin{
		{
			PinID:       1,
			AuthorID:    1,
			Title:       "Pin 1",
			Description: "Description 1",
			MediaUrl:    "https://images.unsplash.com/photo-1655635949384-f737c5133dfe?w=500&auto=format&fit=crop&q=60&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxzZWFyY2h8MTN8fG5ldXJhbCUyMG5ldHdvcmtzfGVufDB8MXwwfHx8Mg%3D%3D",
			BoardID:     1,
		},
		{
			PinID:       2,
			AuthorID:    1,
			Title:       "Pin 2",
			Description: "Description 2",
			MediaUrl:    "https://images.unsplash.com/photo-1596348158371-d3a25ec4dcf4?w=500&auto=format&fit=crop&q=60&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxzZWFyY2h8NXx8bmV1cmFsJTIwbmV0d29ya3N8ZW58MHwxfDB8fHwy",
			BoardID:     1,
		},
		{
			PinID:       3,
			AuthorID:    1,
			Title:       "Pin 3",
			Description: "Description 3",
			MediaUrl:    "https://images.unsplash.com/photo-1580618432485-1e08c5039909?w=500&auto=format&fit=crop&q=60&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxzZWFyY2h8NHx8bmV1cmFsJTIwbmV0d29ya3N8ZW58MHwxfDB8fHwy",
			BoardID:     1,
		},
		{
			PinID:       4,
			AuthorID:    1,
			Title:       "Pin 4",
			Description: "Description 4",
			MediaUrl:    "https://images.unsplash.com/photo-1593376893114-1aed528d80cf?w=500&auto=format&fit=crop&q=60&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxzZWFyY2h8M3x8bmV1cmFsJTIwbmV0d29ya3N8ZW58MHwxfDB8fHwy",
			BoardID:     1,
		},
		{
			PinID:       5,
			AuthorID:    1,
			Title:       "Pin 5",
			Description: "Description 5",
			MediaUrl:    "https://images.unsplash.com/photo-1680474569854-81216b34417a?w=500&auto=format&fit=crop&q=60&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxzZWFyY2h8N3x8bmV1cmFsJTIwbmV0d29ya3N8ZW58MHwxfDB8fHwy",
			BoardID:     1,
		},
		{
			PinID:       6,
			AuthorID:    1,
			Title:       "Pin 6",
			Description: "Description 6",
			MediaUrl:    "https://images.unsplash.com/photo-1668395093559-338fa935d929?w=500&auto=format&fit=crop&q=60&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxzZWFyY2h8MTR8fG5ldXJhbCUyMG5ldHdvcmtzfGVufDB8MXwwfHx8Mg%3D%3D",
			BoardID:     1,
		},
		{
			PinID:       7,
			AuthorID:    1,
			Title:       "Pin 7",
			Description: "Description 7",
			MediaUrl:    "https://images.unsplash.com/photo-1655737484103-193c4385e44d?q=80&w=2532&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D",
			BoardID:     1,
		},
		{
			PinID:       8,
			AuthorID:    1,
			Title:       "Pin 8",
			Description: "Description 8",
			MediaUrl:    "https://images.unsplash.com/photo-1530388684420-55a62e95ed82?w=500&auto=format&fit=crop&q=60&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxzZWFyY2h8MjR8fG5ldXJhbCUyMG5ldHdvcmtzfGVufDB8MXwwfHx8Mg%3D%3D",
			BoardID:     1,
		},
		{
			PinID:       9,
			AuthorID:    1,
			Title:       "Pin 9",
			Description: "Description 9",
			MediaUrl:    "https://images.unsplash.com/photo-1587383378486-83d683d9d02d?w=500&auto=format&fit=crop&q=60&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxzZWFyY2h8NDJ8fG5ldXJhbCUyMG5ldHdvcmtzfGVufDB8MXwwfHx8Mg%3D%3D",
			BoardID:     1,
		},
		{
			PinID:       10,
			AuthorID:    1,
			Title:       "Pin 10",
			Description: "Description 10",
			MediaUrl:    "https://images.unsplash.com/photo-1518276780006-c85b06fa3c11?w=500&auto=format&fit=crop&q=60&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxzZWFyY2h8NDB8fG5ldXJhbCUyMG5ldHdvcmtzfGVufDB8MXwwfHx8Mg%3D%3D",
			BoardID:     1,
		},
		{
			PinID:       11,
			AuthorID:    1,
			Title:       "Pin 11",
			Description: "Description 11",
			MediaUrl:    "https://images.unsplash.com/photo-1578259819688-2bf7b20a351a?w=500&auto=format&fit=crop&q=60&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxzZWFyY2h8NTN8fG5ldXJhbCUyMG5ldHdvcmtzfGVufDB8MXwwfHx8Mg%3D%3D",
			BoardID:     1,
		},
		{
			PinID:       12,
			AuthorID:    1,
			Title:       "Pin 12",
			Description: "Description 12",
			MediaUrl:    "https://images.unsplash.com/photo-1614029655965-2464911905a4?w=500&auto=format&fit=crop&q=60&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxzZWFyY2h8NTV8fG5ldXJhbCUyMG5ldHdvcmtzfGVufDB8MXwwfHx8Mg%3D%3D",
			BoardID:     1,
		},
		{
			PinID:       13,
			AuthorID:    1,
			Title:       "Pin 13",
			Description: "Description 13",
			MediaUrl:    "https://images.unsplash.com/photo-1603745871918-d756fb3c2c5e?w=500&auto=format&fit=crop&q=60&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxzZWFyY2h8NjB8fG5ldXJhbCUyMG5ldHdvcmtzfGVufDB8MXwwfHx8Mg%3D%3D",
			BoardID:     1,
		},
		{
			PinID:       14,
			AuthorID:    1,
			Title:       "Pin 14",
			Description: "Description 14",
			MediaUrl:    "https://images.unsplash.com/photo-1700075489227-47f36fb2709b?w=500&auto=format&fit=crop&q=60&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxzZWFyY2h8NjF8fG5ldXJhbCUyMG5ldHdvcmtzfGVufDB8MXwwfHx8Mg%3D%3D",
			BoardID:     1,
		},
		{
			PinID:       15,
			AuthorID:    1,
			Title:       "Pin 15",
			Description: "Description 15",
			MediaUrl:    "https://images.unsplash.com/photo-1613591876822-846e82526ee7?w=500&auto=format&fit=crop&q=60&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxzZWFyY2h8ODF8fG5ldXJhbCUyMG5ldHdvcmtzfGVufDB8MXwwfHx8Mg%3D%3D",
			BoardID:     1,
		},
	}
)

func Feed(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	header.Add("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
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

	feed := models.NewFeed(pins)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(feed); err != nil {
		internal_errors.SendErrorResponse(w, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrFeedNotAccessible,
		})
		return
	}
}
