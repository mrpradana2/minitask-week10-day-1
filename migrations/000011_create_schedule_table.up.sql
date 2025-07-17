CREATE TABLE IF NOT EXISTS schedule (
	id serial4 NOT NULL,
	cinema_id int4 NOT NULL,
	movie_id int4 NOT NULL,
	"location" varchar(255) NOT NULL,
	"date" date NOT NULL,
	times _time NULL,
	price int4 NOT NULL,
	CONSTRAINT schedule_pkey PRIMARY KEY (id),
	CONSTRAINT fk_schedule_cinemas FOREIGN KEY (cinema_id) REFERENCES cinemas (id),
	CONSTRAINT fk_schedule_movies FOREIGN KEY (movie_id) REFERENCES movies (id)
);

