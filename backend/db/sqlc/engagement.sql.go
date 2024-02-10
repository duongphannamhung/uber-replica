// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0
// source: engagement.sql

package db

import (
	"context"
	"database/sql"
)

const createEngagement = `-- name: CreateEngagement :one
INSERT INTO engagements (
    driver_id,
    status,
    latitude,
    longitude,
    geofence_id
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING id, driver_id, status, latitude, longitude, geofence_id, created_at
`

type CreateEngagementParams struct {
	DriverID   sql.NullInt64 `json:"driver_id"`
	Status     int32         `json:"status"`
	Latitude   float64       `json:"latitude"`
	Longitude  float64       `json:"longitude"`
	GeofenceID int32         `json:"geofence_id"`
}

func (q *Queries) CreateEngagement(ctx context.Context, arg CreateEngagementParams) (Engagement, error) {
	row := q.db.QueryRowContext(ctx, createEngagement,
		arg.DriverID,
		arg.Status,
		arg.Latitude,
		arg.Longitude,
		arg.GeofenceID,
	)
	var i Engagement
	err := row.Scan(
		&i.ID,
		&i.DriverID,
		&i.Status,
		&i.Latitude,
		&i.Longitude,
		&i.GeofenceID,
		&i.CreatedAt,
	)
	return i, err
}

const deleteEngagement = `-- name: DeleteEngagement :exec
DELETE FROM engagements
WHERE id = $1
`

func (q *Queries) DeleteEngagement(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteEngagement, id)
	return err
}

const getEngagement = `-- name: GetEngagement :one
SELECT id, driver_id, status, latitude, longitude, geofence_id, created_at FROM engagements
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetEngagement(ctx context.Context, id int64) (Engagement, error) {
	row := q.db.QueryRowContext(ctx, getEngagement, id)
	var i Engagement
	err := row.Scan(
		&i.ID,
		&i.DriverID,
		&i.Status,
		&i.Latitude,
		&i.Longitude,
		&i.GeofenceID,
		&i.CreatedAt,
	)
	return i, err
}

const listEngagements = `-- name: ListEngagements :many
SELECT id, driver_id, status, latitude, longitude, geofence_id, created_at FROM engagements
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListEngagementsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListEngagements(ctx context.Context, arg ListEngagementsParams) ([]Engagement, error) {
	rows, err := q.db.QueryContext(ctx, listEngagements, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Engagement{}
	for rows.Next() {
		var i Engagement
		if err := rows.Scan(
			&i.ID,
			&i.DriverID,
			&i.Status,
			&i.Latitude,
			&i.Longitude,
			&i.GeofenceID,
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

const updateEngagementLatLng = `-- name: UpdateEngagementLatLng :one
UPDATE engagements
SET latitude = $2, longitude = $3
WHERE driver_id = $1
RETURNING id, driver_id, status, latitude, longitude, geofence_id, created_at
`

type UpdateEngagementLatLngParams struct {
	DriverID  sql.NullInt64 `json:"driver_id"`
	Latitude  float64       `json:"latitude"`
	Longitude float64       `json:"longitude"`
}

func (q *Queries) UpdateEngagementLatLng(ctx context.Context, arg UpdateEngagementLatLngParams) (Engagement, error) {
	row := q.db.QueryRowContext(ctx, updateEngagementLatLng, arg.DriverID, arg.Latitude, arg.Longitude)
	var i Engagement
	err := row.Scan(
		&i.ID,
		&i.DriverID,
		&i.Status,
		&i.Latitude,
		&i.Longitude,
		&i.GeofenceID,
		&i.CreatedAt,
	)
	return i, err
}

const updateEngagementStatus = `-- name: UpdateEngagementStatus :one
UPDATE engagements
SET status = $2
WHERE driver_id = $1
RETURNING id, driver_id, status, latitude, longitude, geofence_id, created_at
`

type UpdateEngagementStatusParams struct {
	DriverID sql.NullInt64 `json:"driver_id"`
	Status   int32         `json:"status"`
}

func (q *Queries) UpdateEngagementStatus(ctx context.Context, arg UpdateEngagementStatusParams) (Engagement, error) {
	row := q.db.QueryRowContext(ctx, updateEngagementStatus, arg.DriverID, arg.Status)
	var i Engagement
	err := row.Scan(
		&i.ID,
		&i.DriverID,
		&i.Status,
		&i.Latitude,
		&i.Longitude,
		&i.GeofenceID,
		&i.CreatedAt,
	)
	return i, err
}