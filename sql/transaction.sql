DROP TABLE IF EXISTS transaction;

CREATE TABLE transaction (
	tld varchar(128),
	registrar_ext_id varchar(32),
	registrar_name varchar(255),
	server_transaction_id int,
	command varchar(128),
	object_type varchar(128),
	object_name varchar(255),
	transaction_date varchar(20),
	created timestamp DEFAULT current_timestamp
);

CREATE UNIQUE INDEX transaction__trid_idx ON transaction ( server_transaction_id );
CREATE INDEX transaction__tld_idx ON transaction ( tld );
CREATE INDEX transaction__ot_idx ON transaction ( object_type );
CREATE INDEX transaction__on_idx ON transaction ( object_name );
CREATE INDEX transaction__rei_idx ON transaction ( registrar_ext_id );
