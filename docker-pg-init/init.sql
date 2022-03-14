CREATE ROLE dcard_db_admin LOGIN PASSWORD 'admin_password' NOSUPERUSER INHERIT NOCREATEDB NOCREATEROLE NOREPLICATION;
CREATE ROLE dcard_db_user LOGIN PASSWORD 'user_password' NOSUPERUSER INHERIT NOCREATEDB NOCREATEROLE NOREPLICATION;

CREATE DATABASE dcard_db WITH ENCODING = 'UTF8' LC_COLLATE = 'en_US.UTF-8' LC_CTYPE = 'en_US.UTF-8' CONNECTION LIMIT = -1 template=template0;
ALTER DATABASE dcard_db OWNER TO dcard_db_admin;

\connect dcard_db;
REVOKE USAGE ON SCHEMA public FROM PUBLIC;
REVOKE CREATE ON SCHEMA public FROM PUBLIC;

GRANT USAGE ON SCHEMA public TO dcard_db_admin;
GRANT CREATE ON SCHEMA public TO dcard_db_admin;
GRANT USAGE ON SCHEMA public TO dcard_db_user;

CREATE TABLE shortener
(
    id SERIAL NOT NULL,
    shortID  CHARACTER(7) NOT NULL,
    originalURL CHARACTER varying(500) NOT NULL,
    expireTime TIMESTAMP,
    CONSTRAINT "short_url" PRIMARY KEY (id)
);

ALTER TABLE shortener OWNER TO dcard_db_admin;
GRANT SELECT, INSERT, UPDATE, DELETE, REFERENCES ON TABLE shortener TO dcard_db_user;
GRANT ALL ON SEQUENCE shortener_id_seq TO dcard_db_user;
