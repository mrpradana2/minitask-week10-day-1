CREATE TABLE IF NOT EXISTS users (
	id serial4 NOT NULL,
	email varchar(255) NOT NULL,
	"password" varchar(255) NOT NULL,
	"role" varchar(255) NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
	modified_at timestamp NULL,
	CONSTRAINT users_email_key UNIQUE (email),
	CONSTRAINT users_pkey PRIMARY KEY (id)
);