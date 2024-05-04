package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	db "uber-replica/db/sqlc"
	"uber-replica/token"

	"github.com/gin-gonic/gin"
)

type updateEngagementActiveRequest struct {
	DriverId    string  `json:"driver_id" binding:"required"`
	DriverPhone string  `json:"driver_phone" binding:"required"`
	Status      int32   `json:"status" binding:"required"`
	Latitude    float64 `json:"lat" binding:"required"`
	Longitude   float64 `json:"lng" binding:"required"`
	GeoId       int     `json:"geo_id" binding:"required"`
}

func (server *Server) driverUpdateEngagement(ctx *gin.Context) {
	var request updateEngagementActiveRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.PhoneNumber != request.DriverPhone {
		ctx.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}

	_driver_id, err := strconv.ParseInt(request.DriverId, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	driver_id := int32(_driver_id)

	params := db.UpdateEngagementLatLngParams{
		DriverID:   driver_id,
		Latitude:   sql.NullFloat64{Float64: (request.Latitude), Valid: true},
		Longitude:  sql.NullFloat64{Float64: (request.Longitude), Valid: true},
		GeofenceID: sql.NullInt32{Int32: 1, Valid: true}, // TODO: change this later
	}

	_, err = server.store.UpdateEngagementLatLng(ctx, params)
	if err != nil {
		log.Print("Error on update engagement: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = server.store.UpdateEngagementStatus(
		ctx,
		db.UpdateEngagementStatusParams{
			DriverID: params.DriverID,
			Status:   sql.NullInt32{Int32: request.Status, Valid: true},
		},
	)
	if err != nil {
		log.Print("Error on update engagement: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "Ok")
	return
}

type DriverStatusResponse struct {
	Status int32 `json:"status" binding:"required"`
	TripId int32 `json:"trip_id"`
}

func (server *Server) currentDriverStatus(ctx *gin.Context) {
	_driver_id := ctx.Query("driver_id")
	if _driver_id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing driver_id parameter"})
		return
	}

	__driver_id, err := strconv.ParseInt(_driver_id, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	driver_id := int32(__driver_id)

	curr_engagement, err := server.store.GetEngagementDriver(ctx, driver_id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	resp := DriverStatusResponse{Status: curr_engagement.Status.Int32}
	if curr_engagement.Status.Int32 != 1 && curr_engagement.Status.Int32 != 2 {
		resp.TripId = curr_engagement.InTrip.Int32
	}

	ctx.JSON(http.StatusOK, resp)
	return
}

type UpdateTripFareRequest struct {
	TripId int64 `json:"trip_id" binding:"required"`
	Fare   int32 `json:"fare" binding:"required"`
}

type UpdateTripFareResponse struct {
	TripId int64 `json:"trip_id" binding:"required"`
	Fare   int32 `json:"fare" binding:"required"`
}

func (server *Server) updateTripFare(ctx *gin.Context) {
	val, err := ctx.GetRawData()
	if err != nil {
		log.Print("Error getting raw data: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var update_fare_req UpdateTripFareRequest
	err = json.Unmarshal(val, &update_fare_req)
	if err != nil {
		log.Print("Error unmarshalling update trip fare req: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	_, err = server.store.UpdateTripFare(ctx, db.UpdateTripFareParams{
		ID:   update_fare_req.TripId,
		Fare: sql.NullInt32{Int32: update_fare_req.Fare, Valid: true},
	})
	if err != nil {
		log.Print("Error updating trip fare: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, UpdateTripFareResponse{
		TripId: update_fare_req.TripId,
		Fare:   update_fare_req.Fare,
	})
	return
}

type DoneEngagementRequest struct {
	DriverPhone string `json:"driver_phone" binding:"required"`
	DriverId    string `json:"driver_id" binding:"required"`
	TripId      string `json:"trip_id" binding:"required"`
}

type DoneEngagementResponse struct {
	DriverId        string `json:"driver_id" binding:"required"`
	TripId          string `json:"trip_id" binding:"required"`
	Fare            int32  `json:"fare" binding:"required"`
	DepartureName   string `json:"departure_name" binding:"required"`
	DestinationName string `json:"destination_name" binding:"required"`
	TripCreatedAt   string `json:"trip_created_at" binding:"required"`
}

func (server *Server) finishEngagement(ctx *gin.Context) {
	var request DoneEngagementRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.PhoneNumber != request.DriverPhone {
		ctx.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}

	driver_id, err := strconv.ParseInt(request.DriverId, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	_, err = server.store.UpdateEngagementStatus(
		ctx, db.UpdateEngagementStatusParams{
			DriverID: int32(driver_id),
			Status:   sql.NullInt32{Int32: 1, Valid: true},
		})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	trip_id, err := strconv.ParseInt(request.TripId, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	_, err = server.store.UpdateEngagementTrip(
		ctx,
		db.UpdateEngagementTripParams{
			DriverID: int32(driver_id),
			InTrip:   sql.NullInt32{},
		},
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	trip, err := server.store.GetTrip(ctx, trip_id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, DoneEngagementResponse{
		DriverId:        request.DriverId,
		TripId:          request.TripId,
		Fare:            trip.Fare.Int32,
		DepartureName:   trip.DepartureName,
		DestinationName: trip.DestinationName,
		TripCreatedAt:   trip.CreatedAt.String(),
	})
}

func (server *Server) checkEngagement(ctx *gin.Context) {
	_driver_id := ctx.Query("driver_id")
	if _driver_id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing driver_id parameter"})
		return
	}

	driver_id, err := strconv.ParseInt(_driver_id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	_, err = server.store.GetEngagementDriver(ctx, int32(driver_id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusOK, false)
			return
		}

		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, true)
	return
}
