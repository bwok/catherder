CREATE DATABASE meetupdatabase
  WITH OWNER = meetupuser
  ENCODING = 'UTF8'
  TABLESPACE = pg_default
  LC_COLLATE = 'en_NZ.UTF-8'
  LC_CTYPE = 'en_NZ.UTF-8'
  CONNECTION LIMIT = -1;

-- Or add to an already existing database.

CREATE SCHEMA meetup
  AUTHORIZATION meetupuser;


CREATE TABLE meetup.meetup
(
  idmeetup    bigserial               NOT NULL,
  userhash    character varying(64)   NOT NULL,
  adminhash   character varying(64)   NOT NULL,
  adminemail  character varying(512)  NOT NULL,
  sendalerts  boolean                 NOT NULL,
  dates       bigint[]                NOT NULL,
  description character varying(4096) NOT NULL,
  CONSTRAINT meetup_pk PRIMARY KEY (idmeetup)
)
  WITH (
    OIDS= FALSE
  );
ALTER TABLE meetup.meetup
  OWNER TO meetupuser;

CREATE TABLE meetup."user"
(
  iduser   bigserial              NOT NULL,
  idmeetup bigint                 NOT NULL,
  name     character varying(512) NOT NULL,
  dates    bigint[]               NOT NULL,
  CONSTRAINT user_pk PRIMARY KEY (iduser),
  CONSTRAINT meetup_user_fk FOREIGN KEY (idmeetup)
    REFERENCES meetup.meetup (idmeetup) MATCH SIMPLE
    ON UPDATE NO ACTION ON DELETE CASCADE
)
  WITH (
    OIDS= FALSE
  );
ALTER TABLE meetup."user"
  OWNER TO meetupuser;
