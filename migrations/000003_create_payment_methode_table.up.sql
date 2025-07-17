CREATE TABLE IF NOT EXISTS payment_methode (
	id serial4 NOT NULL,
	"name" varchar(255) NOT NULL,
	image_path varchar(255) NOT NULL,
	CONSTRAINT payment_methode_pkey PRIMARY KEY (id)
);