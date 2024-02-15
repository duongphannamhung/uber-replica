-- name: CreateDriver :one
INSERT INTO drivers (
    phone
) VALUES (
    $1
)
RETURNING *;

-- name: GetDriver :one
SELECT * FROM drivers
WHERE id = $1 LIMIT 1;

-- name: GetDriverByPhone :one
SELECT * FROM drivers
WHERE phone = $1 LIMIT 1;

-- name: ListDrivers :many
SELECT * FROM drivers
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateDriverLoginCode :one
UPDATE drivers
SET login_code = $2
WHERE id = $1
RETURNING *;

-- -- name: UpdateDriver :one
-- UPDATE drivers
-- SET name = $2
-- WHERE id = $1
-- RETURNING *;

-- name: DeleteDriver :exec
DELETE FROM drivers
WHERE id = $1;