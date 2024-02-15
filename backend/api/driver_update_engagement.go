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
		Latitude:   request.Latitude,
		Longitude:  request.Longitude,
		GeofenceID: 1, // TODO: change this later
	}

	_, err = server.store.UpdateEngagementLatLng(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			// User does not exist, create a new one
			// You need to provide the necessary fields for a new user
			_, err = server.store.CreateEngagement(ctx, db.CreateEngagementParams{
				DriverID:   params.DriverID,
				Status:     0,
				Latitude:   params.Latitude,
				Longitude:  params.Longitude,
				GeofenceID: params.GeofenceID,
			})
			if err != nil {
				log.Fatal("Error creating engagement: ", err)
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		} else {
			log.Fatal("Error on update engagement: ", err)
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	_, err = server.store.UpdateEngagementStatus(
		ctx,
		db.UpdateEngagementStatusParams{
			DriverID: params.DriverID,
			Status:   request.Status,
		},
	)
	if err != nil {
		log.Fatal("Error on update engagement: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "Ok")
	return
}

type DriverStatusResponse struct {
	Status int32 `json:"status" binding:"required"`
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

	ctx.JSON(http.StatusOK, DriverStatusResponse{Status: curr_engagement.Status})
	return
}
