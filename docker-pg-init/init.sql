/*
        dcard_db_admin is the database owner.
        It owns all objects in the database.

        dcard_db_user is the account used by the golang executable.
        It is not allow to create /change any object in database.
        For normal table, only CRUD privilege is granted.

        dcard_db_ramdonly is used during debugging.
        It should have select privilege.
*/

CREATE ROLE dcard_db_admin LOGIN PASSWORD 'admin_password' NOSUPERUSER INHERIT NOCREATEDB NOCREATEROLE NOREPLICATION;
CREATE ROLE dcard_db_user LOGIN PASSWORD 'user_password' NOSUPERUSER INHERIT NOCREATEDB NOCREATEROLE NOREPLICATION;
CREATE ROLE dcard_db_ramdonly LOGIN PASSWORD 'ramdonly_password' NOSUPERUSER INHERIT NOCREATEDB NOCREATEROLE NOREPLICATION;

CREATE DATABASE dcard_db WITH ENCODING = 'UTF8' LC_COLLATE = 'en_US.UTF-8' LC_CTYPE = 'en_US.UTF-8' CONNECTION LIMIT = -1 template=template0;
ALTER DATABASE dcard_db OWNER TO dcard_db_admin;

\connect dcard_db;
REVOKE USAGE ON SCHEMA public FROM PUBLIC;
REVOKE CREATE ON SCHEMA public FROM PUBLIC;

GRANT USAGE ON SCHEMA public TO dcard_db_admin;
GRANT CREATE ON SCHEMA public TO dcard_db_admin;

/* grant the schema access privilege to normal users. Without schema right, user will unable to see the tables. */
GRANT USAGE ON SCHEMA public TO dcard_db_user;
GRANT USAGE ON SCHEMA public to demo_readonly;

CREATE TABLE shortener
(
    short_id  INT NOT NULL,             /* URL short ID */
    original_url VARCHAR(500) NOT NULL,     /* URL */
    expire_time TIMESTAMP NOT NULL,         /* Expire time */
    CONSTRAINT "short_url" PRIMARY KEY (short_id)
);

ALTER TABLE shortener OWNER TO dcard_db_admin;

/* grant the CRUD privilege to normal users. grant the select privilege to ramdon users. */
GRANT SELECT, INSERT, UPDATE, DELETE, REFERENCES ON TABLE shortener TO dcard_db_user;
GRANT SELECT ON TABLE shortener TO dcard_db_ramdonly;