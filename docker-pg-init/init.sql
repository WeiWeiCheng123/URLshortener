
CREATE ROLE dcard_admin LOGIN PASSWORD 'admin_password' NOSUPERUSER INHERIT NOCREATEDB NOCREATEROLE NOREPLICATION;

CREATE DATABASE dcard_db with ENCODING = 'UTF8' LC_COLLATE = 'en_US.UTF-8' LC_CTYPE = 'en_US.UTF-8' CONNECTION LIMIT = -1 template=template0;
ALTER DATABASE dcard_db OWNER TO dcard_admin;

CREATE TABLE shortenerDB
(
    uid serial NOT NULL,
    shortID character varying(100) NOT NULL,
    originalURL character varying(500) NOT NULL,
    expireTime character varying(20) ,
    CONSTRAINT userinfo_pkey PRIMARY KEY (uid)
)
WITH (OIDS=FALSE);

ALTER TABLE shortenerDB OWNER TO dcard_admin;