PRAGMA foreign_keys = ON;

CREATE TABLE "meetup"(
  "idmeetup" INTEGER PRIMARY KEY NOT NULL,
  "userhash" VARCHAR(64) NOT NULL,
  "adminhash" VARCHAR(64) NOT NULL,
  "adminemail" VARCHAR(512),
  "sendalerts" INTEGER NOT NULL,
  "dates" BLOB,
  "description" VARCHAR(4096) NOT NULL
);
CREATE TABLE "user"(
  "iduser" INTEGER PRIMARY KEY NOT NULL,
  "meetup_idmeetup" INTEGER NOT NULL,
  "name" VARCHAR(512) NOT NULL,
  "dates" BLOB,
  CONSTRAINT "fk_user_date"
    FOREIGN KEY("meetup_idmeetup")
    REFERENCES "meetup"("idmeetup")
    ON DELETE CASCADE
);
CREATE INDEX "user.fk_user_meetup_idx" ON "user" ("meetup_idmeetup");