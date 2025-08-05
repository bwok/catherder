PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS meetup
(
    idmeetup    INTEGER PRIMARY KEY ASC,
    userhash    TEXT    NOT NULL,
    adminhash   TEXT    NOT NULL,
    dates       BLOB,
    description TEXT    NOT NULL
);

CREATE TABLE IF NOT EXISTS "user"
(
    iduser   INTEGER PRIMARY KEY ASC NOT NULL,
    idmeetup INTEGER                 NOT NULL,
    name     TEXT                    NOT NULL,
    dates    BLOB,
    FOREIGN KEY (idmeetup) REFERENCES meetup (idmeetup) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS "user.fk_user_meetup_idx" ON "user" ("idmeetup");
