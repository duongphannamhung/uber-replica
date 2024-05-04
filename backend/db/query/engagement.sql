-- name: CreateEngagement :one
INSERT INTO engagements (
    driver_id,
    name,
    vehicle_id,
    label,
    model,
    color,
    license_plate
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetEngagement :one
SELECT * FROM engagements
WHERE id = $1 LIMIT 1;

-- name: GetEngagementDriver :one
SELECT * FROM engagements
WHERE driver_id = $1 LIMIT 1;

-- name: GetActiveEngagementInGeoWithVehicle :one
SELECT * FROM engagements
WHERE geofence_id = $1 
AND status = 2 
AND vehicle_id = $2
LIMIT 1;

-- name: ListEngagements :many
SELECT * FROM engagements
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateEngagementLatLng :one
UPDATE engagements
SET latitude = $2, longitude = $3, geofence_id = $4
WHERE driver_id = $1
RETURNING *;

-- name: UpdateEngagementStatus :one
UPDATE engagements
SET status = $2
WHERE driver_id = $1
RETURNING *;

-- name: UpdateEngagementTrip :one
UPDATE engagements
SET in_trip = $2
WHERE driver_id = $1
RETURNING *;


-- name: GetDriverInfo :one
SELECT driver_id,
    name,
    label,
    model,
    color,
    license_plate
FROM engagements
WHERE driver_id = $1
AND vehicle_id = $2;

-- -- name: UpdateDriverInfo :one
-- UPDATE engagements
-- SET name = $2,
--     vehicle_id = $3,
--     label = $4,
--     model = $5,
--     color = $6,
--     license_plate = $7
-- WHERE id = $1
-- RETURNING *;

-- name: DeleteEngagement :exec
DELETE FROM engagements
WHERE id = $1;