CREATE ROLE dcard_admin LOGIN PASSWORD 'password123' NOSUPERUSER INHERIT NOCREATEDB NOCREATEROLE NOREPLICATION;

CREATE DATABASE dcard_db with ENCODING = 'UTF8' LC_COLLATE = 'en_US.UTF-8' LC_CTYPE = 'en_US.UTF-8' CONNECTION LIMIT = -1 template=template0;
ALTER DATABASE dcard_db OWNER TO dcard_admin;

\connect dcard_db;

CREATE TABLE shortenerdb
(
    id serial NOT NULL,
    shortID  character(11) NOT NULL,
    originalURL character varying(500) NOT NULL,
    expireTime character varying(20) NOT NULL,
    CONSTRAINT "short_url" PRIMARY KEY (id)
);

ALTER TABLE shortenerdb OWNER TO dcard_admin;