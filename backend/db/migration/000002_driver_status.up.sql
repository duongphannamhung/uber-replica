CREATE TABLE "engagements" (
  "id" bigserial PRIMARY KEY,
  "driver_id" bigserial references drivers(id),
  "status" integer NOT NULL,  -- 0: available, 1: enroute, 2: on_trip
  "latitude" float NOT NULL,
  "longitude" float NOT NULL,
  "geofence_id" integer NOT NULL DEFAULT (1), -- current, no geofence, so default 1
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "engagements" ("driver_id");