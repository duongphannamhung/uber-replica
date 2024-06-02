-- name: CreateUser :one
INSERT INTO users (
    phone
) VALUES (
    $1
)
RETURNING *;

-- name: CreateUserWithName :one
INSERT INTO users (
    name, phone
) VALUES (
    $1, $2
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByPhone :one
SELECT * FROM users
WHERE phone = $1 LIMIT 1;

-- name: UpdateUserLoginCode :one
UPDATE users
SET login_code = $2
WHERE id = $1
RETURNING *;

-- name: CountNewUsers :one
SELECT COUNT(distinct CASE WHEN DATE(created_at) >= sqlc.arg('start_date')::text::timestamp THEN id ELSE NULL END) as count_users_in_period,
   COUNT(distinct CASE WHEN DATE(created_at) <= sqlc.arg('start_date')::text::timestamp THEN id ELSE NULL END) as count_users_previous_period
FROM users
WHERE DATE(created_at) 
>= DATE_TRUNC('day', 
    sqlc.arg('start_date')::timestamp - CONCAT(DATE_PART('day', 
        sqlc.arg('end_date')::text::timestamp - sqlc.arg('start_date')::timestamp
    )::text, ' day')::interval) AND DATE(created_at) <= sqlc.arg('end_date')::timestamp
;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;