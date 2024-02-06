package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type OSRMResponse struct {
	Code      string         `json:"code" binding:"required"`
	Routes    []OSRMRoute    `json:"routes" binding:"required"`
	Waypoints []OSRMWaypoint `json:"waypoints" binding:"required"`
}

type OSRMRoute struct {
	Geometry   string         `json:"geometry" binding:"required"`
	Legs       []OSRMRouteLeg `json:"legs"`
	Distance   float64        `json:"distance" binding:"required"`
	Duration   float64        `json:"duration" binding:"required"`
	WeightName string         `json:"weight_name" binding:"required"`
	Weight     float64        `json:"weight" binding:"required"`
}

type OSRMRouteLeg struct {
	Step     []struct{} `json:"steps"`
	Distance float64    `json:"distance"`
	Duration float64    `json:"duration"`
	Summary  string     `json:"summary"`
	Weight   float64    `json:"weight"`
}

type OSRMWaypoint struct {
	Hint     string    `json:"hint" binding:"required"`
	Distance float64   `json:"distance" binding:"required"`
	Name     string    `json:"name"`
	Location []float64 `json:"location" binding:"required"`
}

type DistanceDuration struct {
	Distance     float64 `json:"distance" binding:"required"`
	DistanceText string  `json:"distance_text" binding:"required"`
	Duration     float64 `json:"duration" binding:"required"`
	DurationText string  `json:"duration_text" binding:"required"`
}

func (server *Server) getDistance(ctx *gin.Context) {
	departure := convertGooglePointToOSRMPoint(ctx.Param("departure"))
	destination := convertGooglePointToOSRMPoint(ctx.Param("destination"))

	// TODO: change localhost to .env
	response, err := http.Get("http://localhost:8282/route/v1/driving/" + departure + ";" + destination)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var osrmRespObj OSRMResponse
	json.Unmarshal(responseData, &osrmRespObj)

	if len(osrmRespObj.Routes) == 0 {
		ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("No routes found")))
		return
	}

	distanceObj := DistanceDuration{
		Distance:     osrmRespObj.Routes[0].Distance,
		DistanceText: fmt.Sprintf("%.2f", osrmRespObj.Routes[0].Distance/1000) + " km",
		Duration:     osrmRespObj.Routes[0].Duration,
		DurationText: fmt.Sprintf("%.2f", osrmRespObj.Routes[0].Duration/180) + " mins",
	}

	ctx.JSON(http.StatusOK, distanceObj)
}

func convertGooglePointToOSRMPoint(point string) string {
	point = strings.Replace(point, "(", "", -1)
	point = strings.Replace(point, ")", "", -1)
	point = strings.Replace(point, "%20", "", -1)
	point = strings.Replace(point, " ", "", -1)
	lat := strings.Split(point, ",")[0]
	long := strings.Split(point, ",")[1]
	return long + "," + lat
}
