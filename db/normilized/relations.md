/*
Условие 1НФ:
- Все атрибуты являются простыми
- Нет повторений строк в таблице

Условие 2НФ:
- Таблица находится в 1НФ
- Каждый атрибут целиком зависит от первичного ключа

Условие 3НФ:
- Таблица находится во 2НФ
- Отсутствуют транзитивные зависимости

Условие НФБК:
- Таблица находится в 3НФ
- Ключевые атрибуты не должны зависеть от не-ключевых
*/

/*
User DB:
Таблица-хранилище пользователей.
*/
CREATE TABLE IF NOT EXISTS "user" (
    user_id SERIAL PRIMARY KEY,
    user_name TEXT
        CONSTRAINT user_name_length CHECK (CHAR_LENGTH(user_name) <= 255)
        NOT NULL,
    nick_name TEXT
        CONSTRAINT nick_name_length CHECK (CHAR_LENGTH(nick_name) <= 255)
        UNIQUE
        NOT NULL,
    email TEXT
        CONSTRAINT email_length CHECK (CHAR_LENGTH(email) <= 255)
        UNIQUE
        NOT NULL,
	password TEXT
        CONSTRAINT password_length CHECK (
            CHAR_LENGTH(password) >= 8 AND
            CHAR_LENGTH(password) <= 24),
    birth_time TIMESTAMPTZ,
    gender TEXT
        CONSTRAINT gender_length CHECK (CHAR_LENGTH(gender) <= 20)
        NOT NULL,
    avatar_url TEXT,
    creation_time TIMESTAMPTZ DEFAULT NOW(),
    update_time TIMESTAMPTZ DEFAULT NOW()
);

/*
Board DB:
Таблица-хранилище досок, в которые можно сохранять пины.
*/
CREATE TABLE IF NOT EXISTS board (
    board_id SERIAL PRIMARY KEY,
    owner_id INT REFERENCES "user"(user_id) 
        ON DELETE CASCADE
        NOT NULL,
	name TEXT 
        CONSTRAINT board_name_length CHECK (CHAR_LENGTH(name) <= 255)
        NOT NULL,
	description TEXT
        CONSTRAINT description_text CHECK (CHAR_LENGTH(description) <= 500),
    public BOOLEAN 
        NOT NULL
        DEFAULT false,
    creation_time TIMESTAMPTZ DEFAULT NOW(),
    update_time TIMESTAMPTZ DEFAULT NOW()
);

/*
Comment DB:
Таблица-хранилище комментов пользователей под пинами.
*/
CREATE TABLE IF NOT EXISTS comment (
    comment_id SERIAL PRIMARY KEY,
    pin_id INT REFERENCES pin(pin_id)
        ON DELETE CASCADE
        NOT NULL,
    author_id INT REFERENCES "user"(user_id)
        ON DELETE CASCADE
        NOT NULL,
    body TEXT
        CONSTRAINT body_textlength CHECK(CHAR_LENGTH(body) <= 500)
        NOT NULL,
    creation_time TIMESTAMPTZ DEFAULT NOW(),
    update_time TIMESTAMPTZ DEFAULT NOW()
);

/*
Section DB:
Таблица-хранилище разделов в досках для хранения пинов.
*/
CREATE TABLE IF NOT EXISTS section (
    section_id SERIAL PRIMARY KEY,
    board_id INT REFERENCES board(board_id)
        ON DELETE CASCADE
        NOT NULL,
	name TEXT
        CONSTRAINT name_length CHECK(CHAR_LENGTH(name) <= 255)
        NOT NULL,
	description TEXT
        CONSTRAINT decription_length CHECK(CHAR_LENGTH(description) <= 500),
    creation_time TIMESTAMPTZ DEFAULT NOW(),
    update_time TIMESTAMPTZ DEFAULT NOW()
);

/*
Pin DB:
Таблица-хранилище пинов.
board_id указанный в атрибутах ссылается на закрытую доску пользователя.
*/
CREATE TABLE IF NOT EXISTS pin (
    pin_id SERIAL PRIMARY KEY,
	author_id INT REFERENCES "user"(user_id)
        ON DELETE CASCADE
        NOT NULL,
    title TEXT
        CONSTRAINT name_length CHECK(CHAR_LENGTH(title) <= 255)
        NOT NULL,
	description TEXT
        CONSTRAINT decription_length CHECK(CHAR_LENGTH(description) <= 500),
    board_id INT REFERENCES board(board_id)
        ON DELETE CASCADE
        NOT NULL,
    media_url TEXT
        NOT NULL,
    related_link TEXT
        NOT NULL,
    creation_time TIMESTAMPTZ DEFAULT NOW(),
    update_time TIMESTAMPTZ DEFAULT NOW()
);

/*
Bookmark DB:
Таблица хранилище пинов, сохраненных в закладки пользователя.
*/
CREATE TABLE IF NOT EXISTS bookmark (
    bookmark_id SERIAL PRIMARY KEY,
    pin_id INT REFERENCES pin(pin_id)
        ON DELETE CASCADE
        NOT NULL,
    bookmark_time TIMESTAMPTZ
);

/*
Saved pin to board DB:
Таблица-хранилище соответствия досок-сохраненных пинов.
*/
CREATE TABLE IF NOT EXISTS saved_pin_to_board (
    board_id INT REFERENCES board(board_id)
        ON DELETE CASCADE
        NOT NULL,
    pin_id INT REFERENCES pin(pin_id)
        ON DELETE CASCADE
        NOT NULL,
    PRIMARY KEY (board_id, pin_id)
);

/*
Saved pin to section DB:
Таблица-хранилище соответствия разделов-сохраненных пинов.
*/
CREATE TABLE IF NOT EXISTS saved_pin_to_section (
    section_id INT REFERENCES section(section_id)
        ON DELETE CASCADE
        NOT NULL,
    pin_id INT REFERENCES pin(pin_id)
        ON DELETE CASCADE
        NOT NULL,
    PRIMARY KEY (section_id, pin_id)
);

/*
Follower DB:
Таблица-хранилище подписчиков и подписок.
*/
CREATE TABLE IF NOT EXISTS follower (
    follower_id INT REFERENCES "user"(user_id)
        ON DELETE CASCADE
        NOT NULL,
    following_id INT REFERENCES "user"(user_id)
        ON DELETE CASCADE
        NOT NULL,
    PRIMARY KEY (follower_id, following_id)
);

/*
Saved boards DB:
Таблица-хранилище сохраненных досок пользователя.
*/
CREATE TABLE IF NOT EXISTS saved_board (
    user_id INT REFERENCES "user"(user_id)
        ON DELETE CASCADE
        NOT NULL,
    board_id INT REFERENCES board(board_id)
        ON DELETE CASCADE
        NOT NULL,
    PRIMARY KEY (user_id, board_id)
);

/*
Relation user:
{user_id} -> user_name, nick_name, email, password, birth_time, gender, avatar_url, creation_time, update_time 

Relation board:
{board_id} -> owner_id, name, decription, public, creation_time, update_time

Relation comment:
{comment_id} -> owner_id, pin_id, author_id, body, creation_time, update_time

Relation section:
{section_id} -> board_id, name, description, creation_time, update_time

Relation pin:
{pin_id} -> author_id, title, description, board_id, media_url, related_link, creation_time, update_time

Relation bookmark:
{bookmark_id} -> pin_id, creation_time

Relation saved_pin_to_board:
{board_id, pin_id} -> ∅ 

Relation saved_pin_to_section:
{section_id, pin_id} -> ∅

Relation follower:
{follower_id, following_id} -> ∅

Relation saved_board:
{user_id, board_id} -> ∅
*/
