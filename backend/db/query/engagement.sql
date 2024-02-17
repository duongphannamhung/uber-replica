-- name: CreateEngagement :one
INSERT INTO engagements (
    driver_id,
    status,
    latitude,
    longitude,
    geofence_id
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetEngagement :one
SELECT * FROM engagements
WHERE id = $1 LIMIT 1;

-- name: GetEngagementDriver :one
SELECT * FROM engagements
WHERE driver_id = $1 LIMIT 1;

-- name: GetActiveEngagementInGeo :one
SELECT * FROM engagements
WHERE geofence_id = $1 AND status = 2 LIMIT 1;

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

-- name: DeleteEngagement :exec
DELETE FROM engagements
WHERE id = $1;