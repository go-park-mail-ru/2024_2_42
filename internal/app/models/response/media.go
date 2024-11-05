package response

import "time"

type (
	MediaUploadResponse struct {
		Message string `json:"message"`
	}

	PinCreatedResponse struct {
		PinID   uint64 `json:"pin_id"`
		Message string `json:"message"`
	}

	PinPreviewResponse struct {
		AuthorName            string `json:"author_name"`
		AuthorAvatarUrl       string `json:"avatar_url"`
		AuthorFollowersNumber uint64 `json:"followers_count"`
		MediaUrl              string `json:"media_url"`
		ViewsNumber           uint64 `json:"views_count"`
		BookmarksNumber       uint64 `json:"bookmarks_count"`
	}

	PinPageResponse struct {
		AuthorName            string    `json:"author_name"`
		AuthorAvatarUrl       string    `json:"avatar_url"`
		AuthorFollowersNumber uint64    `json:"followers_count"`
		MediaUrl              string    `json:"media_url"`
		Title                 string    `json:"title"`
		Description           string    `json:"description"`
		RelatedLink           string    `json:"related_link"`
		Geolocation           string    `json:"geolocation"`
		CreationTime          time.Time `json:"creation_time"`
	}

	ResponseBookmarkExists struct {
		BookmarkID uint64 `json:"bookmark_id"`
	}

	BoardCreatedResponse struct {
		BoardID uint64 `json:"board_id"`
		Message string `json:"message"`
	}

	BoardResponse struct {
		BoardID      uint64    `json:"board_id"`
		OwnerID      uint64    `json:"owner_id"`
		Cover        string    `json:"board_cover"`
		Title        string    `json:"title"`
		Description  string    `json:"description"`
		Public       bool      `json:"public"`
		CreationTime time.Time `json:"creation_time"`
		UpdateTime   time.Time `json:"update_time"`
	}
)
