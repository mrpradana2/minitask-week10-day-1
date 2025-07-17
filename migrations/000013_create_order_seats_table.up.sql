CREATE TABLE IF NOT EXISTS order_seats (
	id serial4 NOT NULL,
	order_id int4 NOT NULL,
	seat_id int4 NOT NULL,
	CONSTRAINT order_seats_pkey PRIMARY KEY (id),
	CONSTRAINT fk_order_seats_orders FOREIGN KEY (order_id) REFERENCES orders (id), 
	CONSTRAINT fk_order_seats_seats FOREIGN KEY (seat_id) REFERENCES seats (id) 
);