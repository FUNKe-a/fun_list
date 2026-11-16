-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE app_users (
	user_id INTEGER PRIMARY KEY,
	username TEXT UNIQUE NOT NULL,
	email TEXT UNIQUE,
	password_hash TEXT NOT NULL,
	role TEXT NOT NULL DEFAULT 'user' 
	CHECK (role IN ('user', 'moderator', 'admin'))
);

CREATE TABLE directors (
	director_id INTEGER PRIMARY KEY,
	name TEXT NOT NULL,
	date_of_birth INTEGER,
	date_of_death INTEGER,
	biography TEXT
);

CREATE TABLE media (
	media_id INTEGER PRIMARY KEY,
	title TEXT NOT NULL,
	director_id INTEGER NOT NULL,
    released_at INTEGER,
	synopsis TEXT,
    FOREIGN KEY (director_id) REFERENCES directors(director_id)
);

CREATE TABLE comments (
	comment_id INTEGER PRIMARY KEY,
	user_id INTEGER NOT NULL,
	media_id INTEGER NOT NULL,
	content TEXT NOT NULL,
	rating INTEGER NOT NULL CHECK(rating >= 1 AND RATING <= 10),
	created_at INTEGER NOT NULL DEFAULT (unixepoch()),
	FOREIGN KEY (user_id) REFERENCES app_users(user_id),
    FOREIGN KEY (media_id) REFERENCES media(media_id)
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE comments;
DROP TABLE media;
DROP TABLE directors;
DROP TABLE app_users;
