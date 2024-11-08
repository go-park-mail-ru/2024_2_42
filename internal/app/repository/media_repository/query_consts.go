package mediarepository

// Pins CRUD
const (
	GetUserInfoForPin = `SELECT nick_name, avatar_url FROM "user" WHERE user_id = $1`

	CreatePin = `INSERT INTO pin (author_id, title, description, media_url, related_link) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING RETURNING pin_id;`

	GetAllPins               = `SELECT pin_id, author_id, media_url, title, description FROM pin;`
	GetPinPreviewInfoByPinID = `SELECT pin_id, author_id, media_url, views FROM pin WHERE pin_id = $1;`
	GetPinPageInfoByPinID    = `SELECT pin_id, author_id, title, description, related_link, geolocation, creation_time FROM pin WHERE pin_id = $1;`
	GetPinAuthorByUserID     = `SELECT user_name, avatar_url FROM "user" WHERE user_id = $1`

	UpdatePinInfoByPinID       = `UPDATE pin SET title = $1, description = $2, board_id = $3, media_url = $4, related_link = $5, geolocation = $6 WHERE pin_id = $7`
	UpdatePinViewsByPinID      = `UPDATE pin SET views = views + 1 WHERE pin_id = $1;`
	UpdatePinUpdateTimeByPinID = `UPDATE pin SET update_time = $1 WHERE pin_id = $2;`

	DeletePinByPinID = `DELETE FROM pin WHERE pin_id = $1;`

	// Related things
	GetAllCommentariesByPinID     = `SELECT * FROM comment WHERE pin_id = $1;`
	GetPinBookmarksNumberByPinID  = `SELECT COUNT(bookmark_id) FROM bookmark WHERE pin_id = $1;`
	GetBookmarkOnUserPin          = `SELECT bookmark_id FROM bookmark WHERE owner_id = $1 AND pin_id = $2`
	CreatePinBookmark             = `INSERT INTO bookmark (owner_id, pind_id, bookmark_time) VALUES($1, $2, $3) ON CONFLICT DO NOTHING RETURNING bookmark_id;`
	DeletePinBookmarkByBookmarkID = `DELETE FROM bookmark WHERE bookmark_id = $1;`
)

// Boards
const (
	GetAllBoardsByOwnerID = `SELECT * FROM BOARD WHERE owner_id = $1`
	GetBoardByBoardID     = `SELECT * FROM board WHERE board_id = $1`

	CreateBoard           = `INSERT INTO board (owner_id, name, description, public) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING RETURNING board_id;`
	UpdateBoardByBoardID  = `UPDATE board SET name = $1, description = $2, public = $3 RETURNING board_id;`
	DeleteBoardByBoardID  = `DELETE FROM board WHERE board_id = $1`
	GetFirstUserBoardByID = `SELECT board_id FROM saved_board WHERE user_id = $1 ORDER BY board_id ASC LIMIT 1`
)
