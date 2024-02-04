-- name: CreateDriver :one
INSERT INTO drivers (
    year,
    make,
    model,
    color,
    license_plate
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetDriver :one
SELECT * FROM drivers
WHERE id = $1 LIMIT 1;

-- name: ListDrivers :many
SELECT * FROM drivers
ORDER BY id
LIMIT $1
OFFSET $2;

-- -- name: UpdateDriver :one
-- UPDATE drivers
-- SET name = $2
-- WHERE id = $1
-- RETURNING *;

-- name: DeleteDriver :exec
DELETE FROM drivers
WHERE id = $1;