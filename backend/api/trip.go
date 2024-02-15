package api

import (
	"net/http"
	"strconv"
	db "uber-replica/db/sqlc"
	"uber-replica/token"

	"github.com/gin-gonic/gin"
)

type TripBikeRequest struct {
	UserId           string   `json:"user_id" binding:"required"`
	UserPhone        string   `json:"user_phone" binding:"required"`
	Vehicle          string   `json:"vehicle" binding:"required"`
	OriginPoint      GeoPoint `json:"origin_point" binding:"required"`
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

func (server *Server) createTripBike(ctx *gin.Context) {
	const service_type = 1
	var request TripBikeRequest
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
		ServiceType:          service_type,
		OriginLatitude:       request.OriginPoint.Latitude,
		OriginLongitude:      request.OriginPoint.Longitude,
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
