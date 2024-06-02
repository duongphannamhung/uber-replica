-- name: CreateTrip :one
INSERT INTO trips (
    user_id,
    driver_id,
    service_type,
    departure_latitude,
    departure_longitude,
    departure_name,
    destination_latitude,
    destination_longitude,
    destination_name,
    fare
) VALUES (
    $1, NULL, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: GetTrip :one
SELECT * FROM trips
WHERE id = $1 LIMIT 1;

-- name: ListTrips :many
SELECT * FROM trips
where created_at <= NOW() AT TIME ZONE 'Asia/Bangkok'
ORDER BY created_at DESC
LIMIT $1
OFFSET $2
;


-- name: UpdateStartTrip :one
UPDATE trips
SET driver_id = $2,
    service_type = $3,
    is_started = TRUE,
    driver_location_latitude = $4,
    driver_location_longitude = $5
WHERE id = $1
RETURNING *;

-- name: UpdateTripFare :one
UPDATE trips
SET fare = $2
WHERE id = $1
RETURNING *;

-- name: TotalRevenue :one
SELECT SUM(CASE WHEN DATE(created_at) >= sqlc.arg('start_date')::text::timestamp THEN fare ELSE NULL END) as sum_revenue_in_period,
   SUM(CASE WHEN DATE(created_at) <= sqlc.arg('start_date')::text::timestamp THEN fare ELSE NULL END) as sum_revenue_previous_period
FROM trips
WHERE DATE(created_at) 
>= DATE_TRUNC('day', 
    sqlc.arg('start_date')::timestamp - CONCAT(DATE_PART('day', 
        sqlc.arg('end_date')::text::timestamp - sqlc.arg('start_date')::timestamp
    )::text, ' day')::interval) AND DATE(created_at) <= sqlc.arg('end_date')::timestamp
;

-- name: TotalTrip :one
SELECT COUNT(CASE WHEN DATE(created_at) >= sqlc.arg('start_date')::text::timestamp THEN id ELSE NULL END) as count_trip_in_period,
   COUNT(CASE WHEN DATE(created_at) <= sqlc.arg('start_date')::text::timestamp THEN id ELSE NULL END) as count_trip_previous_period
FROM trips
WHERE DATE(created_at) 
>= DATE_TRUNC('day', 
    sqlc.arg('start_date')::timestamp - CONCAT(DATE_PART('day', 
        sqlc.arg('end_date')::text::timestamp - sqlc.arg('start_date')::timestamp
    )::text, ' day')::interval) AND DATE(created_at) <= sqlc.arg('end_date')::timestamp
;

-- name: RevenueYear :many
SELECT row_number() over (order by month asc) as row_num,
	EXTRACT(MONTH FROM month) as month,
	sum_revenue
FROM
(SELECT 
	DATE_TRUNC('month', created_at) as month,
    SUM(fare) as sum_revenue
FROM trips
where DATE(created_at) >= CURRENT_DATE + INTERVAL '1 month' - INTERVAL '1 year'
GROUP BY 1) as s
order by row_number() over (order by month asc)
limit 12
;



-- -- name: UpdateTrip :one
-- UPDATE users
-- SET name = $2
-- WHERE id = $1
-- RETURNING *;

-- name: DeleteTrip :exec
DELETE FROM trips
WHERE id = $1;

-- name: CountAllTrips :one
SELECT COUNT(*) FROM trips;