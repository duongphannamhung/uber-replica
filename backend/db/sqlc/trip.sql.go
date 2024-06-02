// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0
// source: trip.sql

package db

import (
	"context"
	"database/sql"
)

const countAllTrips = `-- name: CountAllTrips :one
SELECT COUNT(*) FROM trips
`

func (q *Queries) CountAllTrips(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, countAllTrips)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createTrip = `-- name: CreateTrip :one
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
RETURNING id, user_id, driver_id, service_type, is_started, departure_latitude, departure_longitude, departure_name, destination_latitude, destination_longitude, destination_name, driver_location_latitude, driver_location_longitude, fare, created_at
`

type CreateTripParams struct {
	UserID               int64         `json:"user_id"`
	ServiceType          int32         `json:"service_type"`
	DepartureLatitude    float64       `json:"departure_latitude"`
	DepartureLongitude   float64       `json:"departure_longitude"`
	DepartureName        string        `json:"departure_name"`
	DestinationLatitude  float64       `json:"destination_latitude"`
	DestinationLongitude float64       `json:"destination_longitude"`
	DestinationName      string        `json:"destination_name"`
	Fare                 sql.NullInt32 `json:"fare"`
}

func (q *Queries) CreateTrip(ctx context.Context, arg CreateTripParams) (Trip, error) {
	row := q.db.QueryRowContext(ctx, createTrip,
		arg.UserID,
		arg.ServiceType,
		arg.DepartureLatitude,
		arg.DepartureLongitude,
		arg.DepartureName,
		arg.DestinationLatitude,
		arg.DestinationLongitude,
		arg.DestinationName,
		arg.Fare,
	)
	var i Trip
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.DriverID,
		&i.ServiceType,
		&i.IsStarted,
		&i.DepartureLatitude,
		&i.DepartureLongitude,
		&i.DepartureName,
		&i.DestinationLatitude,
		&i.DestinationLongitude,
		&i.DestinationName,
		&i.DriverLocationLatitude,
		&i.DriverLocationLongitude,
		&i.Fare,
		&i.CreatedAt,
	)
	return i, err
}

const deleteTrip = `-- name: DeleteTrip :exec

DELETE FROM trips
WHERE id = $1
`

// -- name: UpdateTrip :one
// UPDATE users
// SET name = $2
// WHERE id = $1
// RETURNING *;
func (q *Queries) DeleteTrip(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteTrip, id)
	return err
}

const getTrip = `-- name: GetTrip :one
SELECT id, user_id, driver_id, service_type, is_started, departure_latitude, departure_longitude, departure_name, destination_latitude, destination_longitude, destination_name, driver_location_latitude, driver_location_longitude, fare, created_at FROM trips
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetTrip(ctx context.Context, id int64) (Trip, error) {
	row := q.db.QueryRowContext(ctx, getTrip, id)
	var i Trip
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.DriverID,
		&i.ServiceType,
		&i.IsStarted,
		&i.DepartureLatitude,
		&i.DepartureLongitude,
		&i.DepartureName,
		&i.DestinationLatitude,
		&i.DestinationLongitude,
		&i.DestinationName,
		&i.DriverLocationLatitude,
		&i.DriverLocationLongitude,
		&i.Fare,
		&i.CreatedAt,
	)
	return i, err
}

const listTrips = `-- name: ListTrips :many
SELECT id, user_id, driver_id, service_type, is_started, departure_latitude, departure_longitude, departure_name, destination_latitude, destination_longitude, destination_name, driver_location_latitude, driver_location_longitude, fare, created_at FROM trips
where created_at <= NOW() AT TIME ZONE 'Asia/Bangkok'
ORDER BY created_at DESC
LIMIT $1
OFFSET $2
`

type ListTripsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListTrips(ctx context.Context, arg ListTripsParams) ([]Trip, error) {
	rows, err := q.db.QueryContext(ctx, listTrips, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Trip{}
	for rows.Next() {
		var i Trip
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.DriverID,
			&i.ServiceType,
			&i.IsStarted,
			&i.DepartureLatitude,
			&i.DepartureLongitude,
			&i.DepartureName,
			&i.DestinationLatitude,
			&i.DestinationLongitude,
			&i.DestinationName,
			&i.DriverLocationLatitude,
			&i.DriverLocationLongitude,
			&i.Fare,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const revenueYear = `-- name: RevenueYear :many
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
`

type RevenueYearRow struct {
	RowNum     int64  `json:"row_num"`
	Month      string `json:"month"`
	SumRevenue int64  `json:"sum_revenue"`
}

func (q *Queries) RevenueYear(ctx context.Context) ([]RevenueYearRow, error) {
	rows, err := q.db.QueryContext(ctx, revenueYear)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []RevenueYearRow{}
	for rows.Next() {
		var i RevenueYearRow
		if err := rows.Scan(&i.RowNum, &i.Month, &i.SumRevenue); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const totalRevenue = `-- name: TotalRevenue :one
SELECT SUM(CASE WHEN DATE(created_at) >= $1::text::timestamp THEN fare ELSE NULL END) as sum_revenue_in_period,
   SUM(CASE WHEN DATE(created_at) <= $1::text::timestamp THEN fare ELSE NULL END) as sum_revenue_previous_period
FROM trips
WHERE DATE(created_at) 
>= DATE_TRUNC('day', 
    $1::timestamp - CONCAT(DATE_PART('day', 
        $2::text::timestamp - $1::timestamp
    )::text, ' day')::interval) AND DATE(created_at) <= $2::timestamp
`

type TotalRevenueParams struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type TotalRevenueRow struct {
	SumRevenueInPeriod       int64 `json:"sum_revenue_in_period"`
	SumRevenuePreviousPeriod int64 `json:"sum_revenue_previous_period"`
}

func (q *Queries) TotalRevenue(ctx context.Context, arg TotalRevenueParams) (TotalRevenueRow, error) {
	row := q.db.QueryRowContext(ctx, totalRevenue, arg.StartDate, arg.EndDate)
	var i TotalRevenueRow
	err := row.Scan(&i.SumRevenueInPeriod, &i.SumRevenuePreviousPeriod)
	return i, err
}

const totalTrip = `-- name: TotalTrip :one
SELECT COUNT(CASE WHEN DATE(created_at) >= $1::text::timestamp THEN id ELSE NULL END) as count_trip_in_period,
   COUNT(CASE WHEN DATE(created_at) <= $1::text::timestamp THEN id ELSE NULL END) as count_trip_previous_period
FROM trips
WHERE DATE(created_at) 
>= DATE_TRUNC('day', 
    $1::timestamp - CONCAT(DATE_PART('day', 
        $2::text::timestamp - $1::timestamp
    )::text, ' day')::interval) AND DATE(created_at) <= $2::timestamp
`

type TotalTripParams struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type TotalTripRow struct {
	CountTripInPeriod       int64 `json:"count_trip_in_period"`
	CountTripPreviousPeriod int64 `json:"count_trip_previous_period"`
}

func (q *Queries) TotalTrip(ctx context.Context, arg TotalTripParams) (TotalTripRow, error) {
	row := q.db.QueryRowContext(ctx, totalTrip, arg.StartDate, arg.EndDate)
	var i TotalTripRow
	err := row.Scan(&i.CountTripInPeriod, &i.CountTripPreviousPeriod)
	return i, err
}

const updateStartTrip = `-- name: UpdateStartTrip :one
UPDATE trips
SET driver_id = $2,
    service_type = $3,
    is_started = TRUE,
    driver_location_latitude = $4,
    driver_location_longitude = $5
WHERE id = $1
RETURNING id, user_id, driver_id, service_type, is_started, departure_latitude, departure_longitude, departure_name, destination_latitude, destination_longitude, destination_name, driver_location_latitude, driver_location_longitude, fare, created_at
`

type UpdateStartTripParams struct {
	ID                      int64           `json:"id"`
	DriverID                sql.NullInt32   `json:"driver_id"`
	ServiceType             int32           `json:"service_type"`
	DriverLocationLatitude  sql.NullFloat64 `json:"driver_location_latitude"`
	DriverLocationLongitude sql.NullFloat64 `json:"driver_location_longitude"`
}

func (q *Queries) UpdateStartTrip(ctx context.Context, arg UpdateStartTripParams) (Trip, error) {
	row := q.db.QueryRowContext(ctx, updateStartTrip,
		arg.ID,
		arg.DriverID,
		arg.ServiceType,
		arg.DriverLocationLatitude,
		arg.DriverLocationLongitude,
	)
	var i Trip
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.DriverID,
		&i.ServiceType,
		&i.IsStarted,
		&i.DepartureLatitude,
		&i.DepartureLongitude,
		&i.DepartureName,
		&i.DestinationLatitude,
		&i.DestinationLongitude,
		&i.DestinationName,
		&i.DriverLocationLatitude,
		&i.DriverLocationLongitude,
		&i.Fare,
		&i.CreatedAt,
	)
	return i, err
}

const updateTripFare = `-- name: UpdateTripFare :one
UPDATE trips
SET fare = $2
WHERE id = $1
RETURNING id, user_id, driver_id, service_type, is_started, departure_latitude, departure_longitude, departure_name, destination_latitude, destination_longitude, destination_name, driver_location_latitude, driver_location_longitude, fare, created_at
`

type UpdateTripFareParams struct {
	ID   int64         `json:"id"`
	Fare sql.NullInt32 `json:"fare"`
}

func (q *Queries) UpdateTripFare(ctx context.Context, arg UpdateTripFareParams) (Trip, error) {
	row := q.db.QueryRowContext(ctx, updateTripFare, arg.ID, arg.Fare)
	var i Trip
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.DriverID,
		&i.ServiceType,
		&i.IsStarted,
		&i.DepartureLatitude,
		&i.DepartureLongitude,
		&i.DepartureName,
		&i.DestinationLatitude,
		&i.DestinationLongitude,
		&i.DestinationName,
		&i.DriverLocationLatitude,
		&i.DriverLocationLongitude,
		&i.Fare,
		&i.CreatedAt,
	)
	return i, err
}
