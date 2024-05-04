package api

import (
	"database/sql"
	"net/http"
	"strconv"
	db "uber-replica/db/sqlc"
	"uber-replica/token"

	"github.com/gin-gonic/gin"
)

type CreateTripRequest struct {
	UserId           string   `json:"user_id" binding:"required"`
	UserPhone        string   `json:"user_phone" binding:"required"`
	Vehicle          int32    `json:"vehicle" binding:"required"`
	DeparturePoint   GeoPoint `json:"departure_point" binding:"required"`
	DepartureName    string   `json:"departure_name" binding:"required"`
	DestinationPoint GeoPoint `json:"destination_point" binding:"required"`
	DestinationName  string   `json:"destination_name" binding:"required"`
}

type GeoPoint struct {
	Latitude  float64 `json:"lat" binding:"required"`
	Longitude float64 `json:"lng" binding:"required"`
}

type TripCreateResp struct {
	TripId int64 `json:"trip_id" binding:"required"`
}

func (server *Server) createTrip(ctx *gin.Context) {
	var request CreateTripRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.PhoneNumber != request.UserPhone {
		ctx.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}

	user_id, err := strconv.ParseInt(request.UserId, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	curr_trip, err := server.store.CreateTrip(ctx, db.CreateTripParams{
		UserID:               user_id,
		ServiceType:          request.Vehicle,
		DepartureLatitude:    request.DeparturePoint.Latitude,
		DepartureLongitude:   request.DeparturePoint.Longitude,
		DepartureName:        request.DepartureName,
		DestinationLatitude:  request.DestinationPoint.Latitude,
		DestinationLongitude: request.DestinationPoint.Longitude,
		DestinationName:      request.DestinationName,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, TripCreateResp{TripId: curr_trip.ID})
	return
}

type getTripInfoResponse struct {
	TripId                  int64           `json:"trip_id"`
	UserID                  int64           `json:"user_id"`
	DriverID                int32           `json:"driver_id"`
	ServiceType             int32           `json:"service_type"`
	IsStarted               bool            `json:"is_started"`
	DepartureName           string          `json:"departure_name"`
	DepartureLatitude       float64         `json:"departure_latitude"`
	DepartureLongitude      float64         `json:"departure_longitude"`
	DestinationLatitude     float64         `json:"destination_latitude"`
	DestinationLongitude    float64         `json:"destination_longitude"`
	DestinationName         string          `json:"destination_name"`
	DriverLocationLatitude  sql.NullFloat64 `json:"driver_location_latitude"`
	DriverLocationLongitude sql.NullFloat64 `json:"driver_location_longitude"`
	Fare                    int32           `json:"fare"`
	TripCreatedAt           string          `json:"trip_created_at"`
}

func (server *Server) getTripInfo(ctx *gin.Context) {
	tripId, err := strconv.ParseInt(ctx.Param("tripId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	trip, err := server.store.GetTrip(ctx, tripId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, getTripInfoResponse{
		TripId:                  trip.ID,
		UserID:                  trip.UserID,
		DriverID:                trip.DriverID.Int32,
		ServiceType:             trip.ServiceType,
		IsStarted:               trip.IsStarted,
		DepartureName:           trip.DepartureName,
		DepartureLatitude:       trip.DepartureLatitude,
		DepartureLongitude:      trip.DepartureLongitude,
		DestinationLatitude:     trip.DestinationLatitude,
		DestinationLongitude:    trip.DestinationLongitude,
		DestinationName:         trip.DestinationName,
		DriverLocationLatitude:  trip.DriverLocationLatitude,
		DriverLocationLongitude: trip.DriverLocationLongitude,
		Fare:                    trip.Fare.Int32,
		TripCreatedAt:           trip.CreatedAt.String(),
	})
	return
}
