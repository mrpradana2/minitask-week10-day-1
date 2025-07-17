CREATE TABLE IF NOT EXISTS movie_casts (
	id serial4 NOT NULL,
	movie_id int4 NULL,
	cast_id int4 NULL,
	CONSTRAINT movie_casts_pkey PRIMARY KEY (id),
	CONSTRAINT fk_movie_casts_movies FOREIGN KEY (movie_id) REFERENCES movies (id),
	CONSTRAINT fk_casts_movies FOREIGN KEY (cast_id) REFERENCES casts (id)
);