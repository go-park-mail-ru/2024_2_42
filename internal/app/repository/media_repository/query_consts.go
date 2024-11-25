package mediarepository

// Pins CRUD
const (
	GetUserInfoForPin = `SELECT nick_name, avatar_url FROM "user" WHERE user_id = $1`

	CreatePin = `INSERT INTO pin (author_id, title, description, media_url, related_link) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING RETURNING pin_id;`

	AddPinToBoard            = `INSERT INTO saved_pin_to_board (board_id, pin_id) VALUES ($1, $2) ON CONFLICT DO NOTHING RETURNING pin_id;`
	DeletePinFromBoard       = `DELETE FROM saved_pin_to_board WHERE board_id = $1 AND pin_id = $2;`
	GetAllPins               = `SELECT pin_id, author_id, media_url, title, description, bookmarks, views FROM pin;`
	GetPinPreviewInfoByPinID = `SELECT pin_id, author_id, media_url, views FROM pin WHERE pin_id = $1;`
	GetPinPageInfoByPinID    = `SELECT pin_id, author_id, title, description, related_link, media_url, geolocation, creation_time FROM pin WHERE pin_id = $1;`
	GetPinAuthorByUserID     = `SELECT nick_name, avatar_url FROM "user" WHERE user_id = $1`
	GetPinsIDByBoardID       = `SELECT pin_id FROM saved_pin_to_board WHERE board_id = $1;`

	UpdatePinInfoByPinID       = `UPDATE pin SET title = $1, description = $2, board_id = $3, media_url = $4, related_link = $5, geolocation = $6 WHERE pin_id = $7`
	UpdatePinViewsByPinID      = `UPDATE pin SET views = views + 1 WHERE pin_id = $1;`
	UpdatePinUpdateTimeByPinID = `UPDATE pin SET update_time = $1 WHERE pin_id = $2;`

	DeletePinByPinID = `DELETE FROM pin WHERE pin_id = $1;`

	// Related things
	GetAllCommentariesByPinID = `SELECT * FROM comment WHERE pin_id = $1;`

	UpdateBookmarksCounter             = `UPDATE pin SET bookmarks = bookmarks + 1 WHERE pin_id = $1;`
	GetPinBookmarksNumberByPinID       = `SELECT COUNT(bookmark_id) FROM bookmark WHERE pin_id = $1;`
	GetBookmarkOnUserPin               = `SELECT bookmark_id FROM bookmark WHERE owner_id = $1 AND pin_id = $2`
	CreatePinBookmark                  = `INSERT INTO bookmark (owner_id, pin_id, bookmark_time) VALUES($1, $2, $3) ON CONFLICT DO NOTHING RETURNING bookmark_id;`
	UpdateBookmarksCountIncrease       = `UPDATE pin SET bookmarks = bookmarks + 1 WHERE pin_id = $1`
	UpdateBookmarksCountDecrease       = `UPDATE pin SET bookmarks = bookmarks - 1 WHERE pin_id = $1`
	DeletePinBookmarkByOwnerIDAndPinID = `DELETE FROM bookmark WHERE owner_id = $1 AND pin_id = $2;`
)

// Boards
const (
	GetAllBoardsByOwnerID = `SELECT * FROM BOARD WHERE owner_id = $1`
	GetBoardByBoardID     = `SELECT * FROM board WHERE board_id = $1`

	CreateBoard          = `INSERT INTO board (owner_id, name, description, public) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING RETURNING board_id;`
	UpdateBoardByBoardID = `UPDATE board SET name = $1, description = $2, public = $3 RETURNING board_id;`
	DeleteBoardByBoardID = `DELETE FROM board WHERE board_id = $1`
)
