CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "name" varchar,
  "phone" varchar UNIQUE NOT NULL,
  "login_code" varchar,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "drivers" (
  "id" bigserial PRIMARY KEY,
  "phone" varchar UNIQUE NOT NULL,
  "name" varchar,
  "login_code" varchar,
  "year" integer,
  "make" text,
  "model" text,
  "color" text,
  "license_plate" text,
  "status" integer NOT NULL DEFAULT (0),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "trips" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigserial NOT NULL,
  "driver_id" integer references drivers(id),
  "service_type" integer NOT NULL, -- 1 : bike
  "is_started" boolean NOT NULL DEFAULT (false),
  "is_completed" boolean NOT NULL DEFAULT (false),
  "origin_latitude" float NOT NULL,
  "origin_longitude" float NOT NULL,
  "destination_latitude" float NOT NULL,
  "destination_longitude" float NOT NULL,
  "destination_name" text NOT NULL,
  "driver_location_latitude" float,
  "driver_location_longitude" float,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "users" ("phone");
CREATE INDEX ON "trips" ("user_id");
CREATE INDEX ON "trips" ("driver_id");
CREATE INDEX ON "trips" ("user_id", "driver_id");

-- COMMENT ON COLUMN "entries"."amount" IS 'can be negative or positive';

-- COMMENT ON COLUMN "transfers"."amount" IS 'must be positive';

ALTER TABLE "trips" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

-- ALTER TABLE "trips" ADD FOREIGN KEY ("driver_id") REFERENCES "drivers" ("id");
