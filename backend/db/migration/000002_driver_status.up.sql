CREATE TABLE "engagements" (
  "id" bigserial PRIMARY KEY,
  "driver_id" integer references drivers(id) NOT NULL,
  "status" integer,  -- 1: off, 2: available, 3: enroute, 4: on_trip
  "in_trip" integer references trips(id),

  "vehicle_id" integer references vehicles(id) NOT NULL, -- info vehicle of driver
  "name" varchar NOT NULL, -- tên tài xế
  "label" text NOT NULL, -- hãng
  "model" text NOT NULL, -- dòng xe
  "color" text NOT NULL, -- màu xe
  "license_plate" text NOT NULL, -- biển số xe

  "latitude" float,
  "longitude" float,
  "geofence_id" integer DEFAULT (1), -- current, no geofence, so default 1
  "created_at" timestamptz NOT NULL DEFAULT (now()),

  UNIQUE ("driver_id", "vehicle_id")
);

CREATE INDEX ON "engagements" ("driver_id");
CREATE INDEX ON "engagements" ("driver_id", "vehicle_id");