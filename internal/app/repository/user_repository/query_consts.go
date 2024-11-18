package userRepository

const (
	// User
	GetUserIDByEmail       = `SELECT user_id FROM "user" WHERE email = $1 LIMIT 1;`
	CreateUser             = `INSERT INTO "user" (nick_name, email, password) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING RETURNING user_id, nick_name;`
	CheckUserCredentials   = `SELECT password FROM "user" WHERE email = $1 LIMIT 1;`
	CheckUserByEmail       = `SELECT user_id FROM "user" WHERE email = $1 LIMIT 1;`
	GetUserAvatar          = `SELECT avatar_url FROM "user" WHERE user_id = $1 LIMIT 1;`
	GetUserInfoByID        = `SELECT user_name, nick_name, description, birth_time, gender, avatar_url FROM "user" WHERE user_id = $1 LIMIT 1;`
	UpdateUserInfoByID     = `UPDATE "user" SET user_name = $1, nick_name = $2, description = $3, birth_time = $4, gender = $5, update_time = NOW(), avatar_url = $6 WHERE user_id = $7 RETURNING user_id;`
	UpdateUserPasswordByID = `UPDATE "user" SET password = $1, update_time = NOW() WHERE user_id = $2 RETURNING user_id;`
	DeleteUserByID         = `DELETE FROM "user" WHERE user_id = $1;`

	// Follower
	FollowUser           = `INSERT INTO "follower" (owner_id, follower_id) VALUES ($1, $2) ON CONFLICT DO NOTHING;`
	UnfollowUser         = `DELETE FROM "follower" WHERE owner_id = $1, follower_id = $2;`
	GetAllFollowings     = `SELECT owner_id FROM "follower" WHERE follower_id = $1;`
	GetAllSubscriptions  = `SELECT following_id FROM "follower" WHERE owner_id = $1;`
	GetFollowingsCount   = `SELECT COUNT(owner_id) FROM "follower" WHERE follower_id = $1;`
	GetSubsriptionsCount = `SELECT COUNT(follower_id) FROM "follower" WHERE owner_id = $1;`

	// Content
	GetBoardsByUserID = `SELECT board_id FROM "saved_boards" WHERE user_id = $1;`
	GetPinsByUserID   = `SELECT pin_id FROM "saved_pins" WHERE user_id = $1`
)
