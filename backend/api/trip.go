package api

import (
	"database/sql"
	"log"
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
	Fare             int32    `json:"fare"`
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
		Fare:                 sql.NullInt32{Int32: request.Fare, Valid: true},
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

type ResponseListTrip struct {
	Page       int32     `json:"page"`
	PerPage    int32     `json:"per_page"`
	Total      int64     `json:"total"`
	TotalPages int64     `json:"total_pages"`
	Data       []db.Trip `json:"data"`
}

func (server *Server) getListTrip(ctx *gin.Context) {
	_limit := ctx.DefaultQuery("limit", "10")
	_offset := ctx.DefaultQuery("offset", "0")

	limit, err := strconv.ParseInt(_limit, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	offset, err := strconv.ParseInt(_offset, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	trips, err := server.store.ListTrips(ctx, db.ListTripsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// Get total count of trips
	total, err := server.store.CountAllTrips(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// Calculate total pages
	totalPages := total / limit
	if total%limit != 0 {
		totalPages++
	}

	response := ResponseListTrip{
		Page:       int32(offset/limit) + 1,
		PerPage:    int32(limit),
		Total:      total,
		TotalPages: totalPages,
		Data:       trips,
	}

	ctx.JSON(http.StatusOK, response)
	return
}

type CreateTripBizopsRequest struct {
	UserName         string   `json:"user_name" binding:"required"`
	UserPhone        string   `json:"user_phone" binding:"required"`
	Vehicle          int32    `json:"vehicle" binding:"required"`
	DeparturePoint   GeoPoint `json:"departure_point" binding:"required"`
	DepartureName    string   `json:"departure_name" binding:"required"`
	DestinationPoint GeoPoint `json:"destination_point" binding:"required"`
	DestinationName  string   `json:"destination_name" binding:"required"`
	Fare             int32    `json:"fare"`
}

func (server *Server) createBizopsTrip(ctx *gin.Context) {
	var request CreateTripBizopsRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := server.store.GetUserByPhone(ctx, request.UserPhone)
	if err != nil {
		if err == sql.ErrNoRows {
			// User does not exist, create a new one
			// You need to provide the necessary fields for a new user
			user, err = server.store.CreateUserWithName(ctx, db.CreateUserWithNameParams{
				Phone: request.UserPhone,
				Name:  sql.NullString{String: request.UserName, Valid: true},
			})
			if err != nil {
				log.Print("Error creating user: ", err)
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		} else {
			log.Print("Error getting user by phone: ", err)
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	curr_trip, err := server.store.CreateTrip(ctx, db.CreateTripParams{
		UserID:               user.ID,
		ServiceType:          request.Vehicle,
		DepartureLatitude:    request.DeparturePoint.Latitude,
		DepartureLongitude:   request.DeparturePoint.Longitude,
		DepartureName:        request.DepartureName,
		DestinationLatitude:  request.DestinationPoint.Latitude,
		DestinationLongitude: request.DestinationPoint.Longitude,
		DestinationName:      request.DestinationName,
		Fare:                 sql.NullInt32{Int32: request.Fare, Valid: true},
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, TripCreateResp{TripId: curr_trip.ID})
	return
}

func (server *Server) deleteTrip(ctx *gin.Context) {
	tripId, err := strconv.ParseInt(ctx.Param("tripId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = server.store.DeleteTrip(ctx, tripId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
