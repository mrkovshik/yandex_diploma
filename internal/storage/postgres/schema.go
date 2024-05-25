package postgres

const Schema = `
CREATE TABLE IF NOT EXISTS users (
	id serial4 NOT NULL,
	login varchar NOT NULL,
	"password" varchar NOT NULL,
	created_at timestamptz NOT NULL,
	balance float4 DEFAULT 0 NOT NULL,
	CONSTRAINT users_pk PRIMARY KEY (id)
);
CREATE TABLE IF NOT EXISTS orders (
	id serial4 NOT NULL,
	order_number varchar NOT NULL,
	user_id int4 NOT NULL,
	uploaded_at timestamptz NOT NULL,
	status varchar DEFAULT 'NEW'::character varying NOT NULL,
	accrual float4 DEFAULT 0 NOT NULL,
	CONSTRAINT orders_pk PRIMARY KEY (id),
	CONSTRAINT orders_unique UNIQUE (order_number)
);                                  
    CREATE TABLE IF NOT EXISTS withdrawals (
		id serial4 NOT NULL,
	amount float4 NOT NULL,
	processed_at timestamptz NOT NULL,
	order_number varchar NOT NULL,
	user_id int4 NOT NULL,
	CONSTRAINT withdrawals_pk PRIMARY KEY (id)
)`

const SchemaDrop = `
Drop TABLE IF  EXISTS users;
Drop TABLE IF  EXISTS orders;
Drop TABLE IF  EXISTS withdrawals;`
