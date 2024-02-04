-- name: CreateTrip :one
INSERT INTO trips (
    user_id,
    driver_id
) VALUES (
    $1, $2
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

-- -- name: UpdateTrip :one
-- UPDATE users
-- SET name = $2
-- WHERE id = $1
-- RETURNING *;

-- name: DeleteTrip :exec
DELETE FROM trips
WHERE id = $1;