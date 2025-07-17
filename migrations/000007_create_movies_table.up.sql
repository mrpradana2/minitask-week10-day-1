CREATE TABLE IF NOT EXISTS movies (
	id serial4 NOT NULL,
	title varchar(255) NOT NULL,
	image_path varchar NOT NULL,
	overview varchar NOT NULL,
	release_date date NOT NULL,
	director_name varchar(255) NOT NULL,
	duration int4 NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	modified_at timestamp DEFAULT now() NULL,
	CONSTRAINT movies_pkey PRIMARY KEY (id)
);
