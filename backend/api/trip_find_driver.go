package api

import (
	"database/sql"
	"net/http"
	"strconv"
	db "uber-replica/db/sqlc"

	"github.com/gin-gonic/gin"
)

func getUserGeofence(lat float64, lng float64) int32 {
	_ = lat * lng
	return 1
}

type noDriverFoundResponse struct {
	FindDone bool `json:"find_done" binding:"required"`
}

type findDriverDoneResponse struct {
	FindDone        bool    `json:"find_done" binding:"required"`
	DriverId        int64   `json:"driver_id" binding:"required"`
	EngagementId    int64   `json:"engagement_id" binding:"required"`
	DriverLatitude  float64 `json:"driver_lat" binding:"required"`
	DriverLongitude float64 `json:"driver_lng" binding:"required"`
}

func (server *Server) tripFindDriver(ctx *gin.Context) {
	_trip_id := ctx.Query("trip_id")
	_vehicle_type := ctx.Query("vehicle_type")

	if _trip_id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing trip_id parameter"})
		return
	}

	trip_id, err := strconv.ParseInt(_trip_id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	vehicle_type, err := strconv.ParseInt(_vehicle_type, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	curr_trip, err := server.store.GetTrip(ctx, trip_id)
	geo_id := getUserGeofence(curr_trip.DepartureLatitude, curr_trip.DepartureLongitude)

	engagement, err := server.store.GetActiveEngagementInGeoWithVehicle(ctx, db.GetActiveEngagementInGeoWithVehicleParams{
		GeofenceID: sql.NullInt32{Int32: int32(geo_id), Valid: true},
		VehicleID:  int32(vehicle_type),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusOK, noDriverFoundResponse{FindDone: false})
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	_, err = server.store.UpdateEngagementStatus(
		ctx,
		db.UpdateEngagementStatusParams{
			DriverID: engagement.DriverID,
			Status:   sql.NullInt32{Int32: 3, Valid: true},
		})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	_, err = server.store.UpdateStartTrip(
		ctx,
		db.UpdateStartTripParams{
			ID:                      trip_id,
			DriverID:                sql.NullInt32{Int32: engagement.DriverID, Valid: true},
			ServiceType:             int32(vehicle_type),
			DriverLocationLatitude:  sql.NullFloat64{Float64: engagement.Latitude.Float64, Valid: true},
			DriverLocationLongitude: sql.NullFloat64{Float64: engagement.Longitude.Float64, Valid: true},
		},
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = server.store.UpdateEngagementTrip(
		ctx,
		db.UpdateEngagementTripParams{
			DriverID: engagement.DriverID,
			InTrip:   sql.NullInt32{Int32: int32(trip_id), Valid: true},
		},
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, findDriverDoneResponse{
		FindDone:        true,
		DriverId:        int64(engagement.DriverID),
		EngagementId:    engagement.ID,
		DriverLatitude:  engagement.Latitude.Float64,
		DriverLongitude: engagement.Longitude.Float64,
	})
	return
}

func (server *Server) getDriverInfo(ctx *gin.Context) {
	_trip_id := ctx.Param("tripId")

	if _trip_id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing trip_id parameter"})
		return
	}

	trip_id, err := strconv.ParseInt(_trip_id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	trip, err := server.store.GetTrip(ctx, trip_id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	driver, err := server.store.GetDriverInfo(ctx, db.GetDriverInfoParams{
		DriverID:  trip.DriverID.Int32,
		VehicleID: trip.ServiceType,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, driver)
	return
}
