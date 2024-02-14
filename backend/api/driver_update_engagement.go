package api

import (
	"net/http"
	"strconv"
	"uber-replica/token"

	"github.com/gin-gonic/gin"
)

type updateEngagementActiveRequest struct {
	DriverId    string  `json:"driver_id" binding:"required"`
	DriverPhone string  `json:"driver_phone" binding:"required"`
	Status      int     `json:"status" binding:"required"`
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

	_, err := strconv.ParseInt(request.DriverId, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, "Ok")
	return
}
