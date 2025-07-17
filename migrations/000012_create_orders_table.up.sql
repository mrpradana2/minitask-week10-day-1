CREATE TABLE IF NOT EXISTS orders (
	id serial4 NOT NULL,
	user_id int4 NOT NULL,
	schedule_id int4 NOT NULL,
	payment_methode_id int4 NOT NULL,
	"date" date NULL,
	"time" time NULL,
	total_price int4 NOT NULL,
	full_name varchar(255) NOT NULL,
	email varchar(255) NOT NULL,
	phone_number varchar(255) NOT NULL,
	paid bool NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT orders_pkey PRIMARY KEY (id),
	CONSTRAINT fk_orders_users FOREIGN KEY (user_id) REFERENCES users (id), 
	CONSTRAINT fk_orders_payment_methode FOREIGN KEY (payment_methode_id) REFERENCES payment_methode (id), 
	CONSTRAINT fk_orders_schedule FOREIGN KEY (schedule_id) REFERENCES schedule (id) 
);