CREATE TABLE IF NOT EXISTS profiles (
	user_id int4 NOT NULL,
	first_name varchar(255) NULL,
	last_name varchar(255) NULL,
	phone_number varchar(255) NULL,
	photo_path varchar NULL,
	title varchar(255) NULL,
	point int4 NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
	modified_at timestamp NULL,
	CONSTRAINT fk_profile_users FOREIGN KEY (user_id) REFERENCES users (id)
);