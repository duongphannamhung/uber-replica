CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "name" varchar,
  "phone" varchar UNIQUE NOT NULL,
  "login_code" varchar,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "trips" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigserial NOT NULL,
  "driver_id" bigserial NOT NULL,
  "is_started" boolean NOT NULL DEFAULT (false),
  "is_completed" boolean NOT NULL DEFAULT (false),
  "origin" json,
  "destination" json,
  "destination_name" text,
  "driver_location" json,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "drivers" (
  "id" bigserial PRIMARY KEY,
  "year" integer,
  "make" text,
  "model" text,
  "color" text,
  "license_plate" text,
  "status" integer NOT NULL DEFAULT (0),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "users" ("phone");
CREATE INDEX ON "trips" ("user_id");
CREATE INDEX ON "trips" ("driver_id");
CREATE INDEX ON "trips" ("user_id", "driver_id");

-- COMMENT ON COLUMN "entries"."amount" IS 'can be negative or positive';

-- COMMENT ON COLUMN "transfers"."amount" IS 'must be positive';

ALTER TABLE "trips" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "trips" ADD FOREIGN KEY ("driver_id") REFERENCES "drivers" ("id");
