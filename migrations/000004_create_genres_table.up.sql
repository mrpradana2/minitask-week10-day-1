CREATE TABLE IF NOT EXISTS genres (
	id serial4 NOT NULL,
	"name" varchar(255) NOT NULL,
	CONSTRAINT genres_pkey PRIMARY KEY (id),
	CONSTRAINT unique_genre_name UNIQUE (name)
);