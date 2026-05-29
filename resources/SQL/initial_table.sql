create table transaction_types
(
	id SERIAL PRIMARY KEY,
	name VARCHAR(50) NOT NULL UNIQUE,
    inserted_at TIMESTAMP NOT NULL
);

create table categories
(
	id SERIAL PRIMARY KEY,
	name VARCHAR(50) NOT NULL UNIQUE,
    inserted_at TIMESTAMP NOT NULL
);

create table account_groups
(
	id SERIAL PRIMARY KEY,
	name VARCHAR(50) NOT NULL UNIQUE,
    inserted_at TIMESTAMP NOT NULL
);

create table transactions
(
	id SERIAL PRIMARY KEY,
	parent_id INTEGER NOT NULL,
	trade_date DATE NOT NULL,
	type_id INTEGER NOT NULL,
	category_id INTEGER NOT NULL,
	account_id INTEGER NOT NULL,
	account_dest_id INTEGER,
	amount NUMERIC(24,2) NOT NULL,
	description VARCHAR(1000),
    status VARCHAR(15) NOT NULL,
    source VARCHAR(15) NOT NULL,
	inserted_at TIMESTAMP NOT NULL,
	last_updated_at TIMESTAMP NOT NULL,
	FOREIGN KEY (type_id) REFERENCES transaction_types(id),
	FOREIGN KEY (category_id) REFERENCES categories(id),
	FOREIGN KEY (account_id) REFERENCES accounts(id),
	FOREIGN KEY (account_dest_id) REFERENCES accounts(id)
);

create table accounts
(
	id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL UNIQUE,
	balance numeric(24,2) NOT NULL,
	account_group_id INTEGER,
	status VARCHAR(1) NOT NULL,
    inserted_at TIMESTAMP NOT NULL,
    closed_at TIMESTAMP NOT NULL,
	FOREIGN KEY (account_group_id) REFERENCES account_groups(id)
);

create table owned_account(
	account_label varchar(30) PRIMARY KEY,
	account_number varchar(20) NOT NULL
);