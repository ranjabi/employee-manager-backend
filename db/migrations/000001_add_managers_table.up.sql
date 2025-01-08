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

COMMIT;