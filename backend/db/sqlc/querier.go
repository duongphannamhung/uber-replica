// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0

package db

import (
	"context"
)

type Querier interface {
	CreateDriver(ctx context.Context, phone string) (Driver, error)
	CreateEngagement(ctx context.Context, arg CreateEngagementParams) (Engagement, error)
	CreateTrip(ctx context.Context, arg CreateTripParams) (Trip, error)
	CreateUser(ctx context.Context, phone string) (User, error)
	DeleteDriver(ctx context.Context, id int64) error
	// -- name: UpdateDriverInfo :one
	// UPDATE engagements
	// SET name = $2,
	//     vehicle_id = $3,
	//     label = $4,
	//     model = $5,
	//     color = $6,
	//     license_plate = $7
	// WHERE id = $1
	// RETURNING *;
	DeleteEngagement(ctx context.Context, id int64) error
	// -- name: UpdateTrip :one
	// UPDATE users
	// SET name = $2
	// WHERE id = $1
	// RETURNING *;
	DeleteTrip(ctx context.Context, id int64) error
	DeleteUser(ctx context.Context, id int64) error
	GetActiveEngagementInGeoWithVehicle(ctx context.Context, arg GetActiveEngagementInGeoWithVehicleParams) (Engagement, error)
	GetDriver(ctx context.Context, id int64) (Driver, error)
	GetDriverByPhone(ctx context.Context, phone string) (Driver, error)
	GetDriverInfo(ctx context.Context, arg GetDriverInfoParams) (GetDriverInfoRow, error)
	GetEngagement(ctx context.Context, id int64) (Engagement, error)
	GetEngagementDriver(ctx context.Context, driverID int32) (Engagement, error)
	GetTrip(ctx context.Context, id int64) (Trip, error)
	GetUser(ctx context.Context, id int64) (User, error)
	GetUserByPhone(ctx context.Context, phone string) (User, error)
	ListDrivers(ctx context.Context, arg ListDriversParams) ([]Driver, error)
	ListEngagements(ctx context.Context, arg ListEngagementsParams) ([]Engagement, error)
	ListTrips(ctx context.Context, arg ListTripsParams) ([]Trip, error)
	UpdateDriverLoginCode(ctx context.Context, arg UpdateDriverLoginCodeParams) (Driver, error)
	UpdateEngagementLatLng(ctx context.Context, arg UpdateEngagementLatLngParams) (Engagement, error)
	UpdateEngagementStatus(ctx context.Context, arg UpdateEngagementStatusParams) (Engagement, error)
	UpdateEngagementTrip(ctx context.Context, arg UpdateEngagementTripParams) (Engagement, error)
	UpdateStartTrip(ctx context.Context, arg UpdateStartTripParams) (Trip, error)
	UpdateTripFare(ctx context.Context, arg UpdateTripFareParams) (Trip, error)
	UpdateUserLoginCode(ctx context.Context, arg UpdateUserLoginCodeParams) (User, error)
}

var _ Querier = (*Queries)(nil)
