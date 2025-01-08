BEGIN;

CREATE TABLE managers (
	id 					uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	email				text NOT NULL UNIQUE,
	password			text NOT NULL,
	name				text,
	user_image_uri		text,
	company_name		text,
	company_image_uri	text
);

INSERT INTO managers (email, password) VALUES ('a@a.a', '$2a$10$blk0sKGaU8pNVuTXWqHrJOTW81/VaPDMNUJsww1Ai3ZpzPBcZbUoO');

COMMIT;