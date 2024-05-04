CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "name" varchar,
  "phone" varchar UNIQUE NOT NULL,
  "login_code" varchar,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "vehicles" (
  "id" bigserial PRIMARY KEY ,
  "name" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

INSERT INTO vehicles (id, name) VALUES ('1', 'UrepBike');
INSERT INTO vehicles (id, name) VALUES ('2', 'UrepCar');
INSERT INTO vehicles (id, name) VALUES ('3', 'UrepCar 7');
INSERT INTO vehicles (id, name) VALUES ('4', 'UrepCar Plus');

CREATE TABLE "drivers" (
  "id" bigserial PRIMARY KEY,
  "phone" varchar NOT NULL,
  "login_code" varchar,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "trips" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigserial NOT NULL,
  "driver_id" integer references drivers(id),
  "service_type" integer NOT NULL, -- 1 : bike
  "is_started" boolean NOT NULL DEFAULT (false),
  "departure_latitude" float NOT NULL,
  "departure_longitude" float NOT NULL,
  "departure_name" text NOT NULL,
  "destination_latitude" float NOT NULL,
  "destination_longitude" float NOT NULL,
  "destination_name" text NOT NULL,
  "driver_location_latitude" float,
  "driver_location_longitude" float,
  "fare" integer,
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
