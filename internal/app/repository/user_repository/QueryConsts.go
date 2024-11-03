package user_repository

const (
	//User
	CreateUser             = `INSERT INTO "user" (user_name, nick_name, email, password) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING RETURNING user_id;`
	CheckUserCredentials   = `SELECT INTO "user" (password) FROM "user" WHERE user_id = $1 LIMIT 1;`
	CheckUserByEmail       = `SELECT user_id FROM "user" WHERE email = $1 LIMIT 1;`
	GetUserInfoByID        = `SELECT user_name, nick_name, description, birth_time, gender, avatar_url FROM "user" WHERE user_id = $1 LIMIT 1;`
	UpdateUserInfoByID     = `UPDATE "user" SET nick_name = $1, description = $2, birth_time = $3, gender = $4, update_time = NOW() WHERE user_id = $5 RETURNING user_id;`
	UpdateUserPasswordByID = `UPDATE "user" SET password = $1, update_time = NOW() WHERE user_id = $2 RETURNING user_id;`
	DeleteUserByID         = `DELETE FROM "user" WHERE user_id = $1;`

	//Follower
	FollowUser          = `INSERT INTO follower (owner_id, following_id) VALUES ($1, $2) ON CONFLICT DO NOTHING;`
	UnfollowUser        = `DELETE FROM follower WHERE owner_id = $1, follower_id = $2;`
	GetAllFollowings    = `SELECT owner_id FROM follower WHERE following_id = $1;`
	GetAllSubscriptions = `SELECT following_id FROM follower WHERE owner_id = $1;`

	//Content
	GetBoardsByUserID = `SELECT board_id FROM saved_boards WHERE user_id = $1;`
	GetPinsByUserID   = `SELECT pin_id FROM saved_pins WHERE user_id = $1`
)
