-- name: CreateTrip :one
INSERT INTO trips (
    user_id,
    driver_id,
    service_type,
    origin_latitude,
    origin_longitude,
    destination_latitude,
    destination_longitude,
    destination_name
) VALUES (
    $1, NULL, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetTrip :one
SELECT * FROM trips
WHERE id = $1 LIMIT 1;

-- name: ListTrips :many
SELECT * FROM trips
ORDER BY id
LIMIT $1
OFFSET $2;


-- name: UpdateStartTrip :one
UPDATE trips
SET driver_id = $2,
    service_type = $3,
    is_started = TRUE,
    driver_location_latitude = $4,
    driver_location_longitude = $5
WHERE id = $1
RETURNING *;

-- -- name: UpdateTrip :one
-- UPDATE users
-- SET name = $2
-- WHERE id = $1
-- RETURNING *;

-- name: DeleteTrip :exec
DELETE FROM trips
WHERE id = $1;