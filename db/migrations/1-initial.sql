-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE app_users (
	user_id INTEGER PRIMARY KEY,
	username TEXT UNIQUE NOT NULL,
	email TEXT UNIQUE,
	password_hash TEXT NOT NULL,
	password_salt TEXT NOT NULL
);

CREATE TABLE directors (
	director_id INTEGER PRIMARY KEY,
	name TEXT NOT NULL,
	date_of_birth TEXT,
	date_of_death TEXT,
	biography TEXT
);

CREATE TABLE media (
	media_id INTEGER PRIMARY KEY,
	title TEXT NOT NULL,
	director_id INTEGER NOT NULL,
    release_year INTEGER,
	synopsis TEXT,
    FOREIGN KEY (director_id) REFERENCES directors(director_id)
);

CREATE TABLE comments (
	comment_id INTEGER PRIMARY KEY,
	user_id INTEGER NOT NULL,
	media_id INTEGER NOT NULL,
	content TEXT NOT NULL,
	rating INTEGER NOT NULL CHECK(rating >= 1 AND RATING <= 10),
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES app_users(user_id),
    FOREIGN KEY (media_id) REFERENCES media(media_id)
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE comments;
DROP TABLE media;
DROP TABLE directors;
DROP TABLE app_users;
