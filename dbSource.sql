PRAGMA foreign_keys = ON;

CREATE TABLE "meetup"(
  "idmeetup" INTEGER PRIMARY KEY NOT NULL,
  "userhash" VARCHAR(64) NOT NULL,
  "adminhash" VARCHAR(64) NOT NULL,
  "description" VARCHAR(4096) NOT NULL
);
CREATE TABLE "date"(
  "iddate" INTEGER PRIMARY KEY NOT NULL,
  "meetup_idmeetup" INTEGER NOT NULL,
  "date" INTEGER NOT NULL,
  CONSTRAINT "fk_date_meetup"
    FOREIGN KEY("meetup_idmeetup")
    REFERENCES "meetup"("idmeetup")
    ON DELETE CASCADE
);
CREATE INDEX "date.fk_date_meetup_idx" ON "date" ("meetup_idmeetup");
CREATE TABLE "user"(
  "iduser" INTEGER PRIMARY KEY NOT NULL,
  "date_iddate" INTEGER NOT NULL,
  "name" VARCHAR(512) NOT NULL,
  "available" INTEGER NOT NULL,
  CONSTRAINT "fk_user_date"
    FOREIGN KEY("date_iddate")
    REFERENCES "date"("iddate")
    ON DELETE CASCADE
);
CREATE INDEX "user.fk_user_date_idx" ON "user" ("date_iddate");
CREATE TABLE "admin"(
  "idadmin" INTEGER PRIMARY KEY NOT NULL,
  "meetup_idmeetup" INTEGER NOT NULL,
  "email" VARCHAR(512),
  "alerts" INTEGER NOT NULL,
  CONSTRAINT "fk_admin_meetup"
    FOREIGN KEY("meetup_idmeetup")
    REFERENCES "meetup"("idmeetup")
    ON DELETE CASCADE
);
CREATE INDEX "admin.fk_admin_meetup_idx" ON "admin" ("meetup_idmeetup");