CREATE TABLE IF NOT EXISTS casts (
	id serial4 NOT NULL,
	"name" varchar(255) NOT NULL,
	CONSTRAINT casts_pkey PRIMARY KEY (id),
	CONSTRAINT unique_cast_name UNIQUE (name)
);