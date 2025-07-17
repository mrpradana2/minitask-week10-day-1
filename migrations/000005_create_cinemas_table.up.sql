CREATE TABLE IF NOT EXISTS cinemas (
	id serial4 NOT NULL,
	"name" varchar(255) NOT NULL,
	image_path varchar NULL,
	CONSTRAINT cinemas_pkey PRIMARY KEY (id)
);