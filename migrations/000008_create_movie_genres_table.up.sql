CREATE TABLE IF NOT EXISTS movie_genre (
	id serial4 NOT NULL,
	movie_id int4 NULL,
	genre_id int4 NULL,
	CONSTRAINT movie_genre_pkey PRIMARY KEY (id),
	CONSTRAINT fk_movie_genre_movies FOREIGN KEY (movie_id) REFERENCES movies (id),
	CONSTRAINT fk_movie_genre_genres FOREIGN KEY (genre_id) REFERENCES genres (id)
);
