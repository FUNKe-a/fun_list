-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE app_users (
	user_id BIGINT PRIMARY KEY,
	username VARCHAR(32) UNIQUE NOT NULL,
	email TEXT UNIQUE,
	password_hash TEXT NOT NULL,
	password_salt TEXT NOT NULL
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE app_users;
