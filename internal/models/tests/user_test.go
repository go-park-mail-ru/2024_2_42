package tests

import (
	"errors"
	"testing"
	internal_errors "youpin/internal/errors"
	"youpin/internal/models"
)

func TestPinValid(t *testing.T) {
	testData := []struct {
		name   string
		input  models.Pin
		result error
	}{
		{
			name: "Valid",
			input: models.Pin{
				PinID:       1,
				AuthorID:    1,
				Title:       "Pin 1",
				Description: "Description 1",
				MediaUrl:    "https://images.unsplash.com/photo-1655635949384-f737c5133dfe?w=500&auto=format&fit=crop&q=60&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxzZWFyY2h8MTN8fG5ldXJhbCUyMG5ldHdvcmtzfGVufDB8MXwwfHx8Mg%3D%3D",
				BoardID:     1,
			},
			result: nil,
		},
		{
			name: "Invalid title",
			input: models.Pin{
				PinID:       1,
				AuthorID:    1,
				Title:       "Pin 1",
				Description: "Description 1",
				MediaUrl:    "https://images.unsplash.com/photo-1655635949384-f737c5133dfe?w=500&auto=format&fit=crop&q=60&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxzZWFyY2h8MTN8fG5ldXJhbCUyMG5ldHdvcmtzfGVufDB8MXwwfHx8Mg%3D%3D",
				BoardID:     1,
			},
			result: internal_errors.ErrPinDataInvalid,
		},
	}
	for _, test := range testData {
		t.Run(test.name, func(t *testing.T) {
			resultErr := test.input.Valid()

			if resultErr != nil {
				if !errors.Is(resultErr, test.result) {
					t.Errorf("expected error %v , got %v", test.result, resultErr)
				}
			}
		})
	}
}
