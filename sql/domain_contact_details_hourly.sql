DROP TABLE IF EXISTS domain_contact_details_hourly;

CREATE TABLE domain_contact_details_hourly (
	domain_id int,
	domain_name varchar(128),
	domain_created_on varchar(27),
	registrar_ext_id varchar(32),
	registrant_client_id varchar(32),
	registrant_name varchar(255),
	registrant_org varchar(255),
	registrant_email varchar(255),
	created timestamp DEFAULT current_timestamp
);

CREATE UNIQUE INDEX dcdh__dn_dco_idx ON domain_contact_details_hourly ( domain_name, domain_created_on );
CREATE INDEX dcdh__dn_idx ON domain_contact_details_hourly ( domain_name );
CREATE INDEX dcdh__dco_idx ON domain_contact_details_hourly ( domain_created_on );
CREATE INDEX dcdh__rei_idx ON domain_contact_details_hourly ( registrar_ext_id );
