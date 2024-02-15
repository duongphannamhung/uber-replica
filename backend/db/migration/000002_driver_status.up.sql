CREATE TABLE "engagements" (
  "id" bigserial PRIMARY KEY,
  "driver_id" integer references drivers(id) NOT NULL,
  "status" integer NOT NULL,  -- 0: off, 1: available, 2: enroute, 3: on_trip
  "latitude" float NOT NULL,
  "longitude" float NOT NULL,
  "geofence_id" integer NOT NULL DEFAULT (1), -- current, no geofence, so default 1
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "engagements" ("driver_id");