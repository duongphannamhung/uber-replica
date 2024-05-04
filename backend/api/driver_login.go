package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
	db "uber-replica/db/sqlc"

	"github.com/gin-gonic/gin"
)

type DriverLoginPhoneRequest struct {
	Phone    string `json:"phone" binding:"required"`
	LoginOTP string `json:"login_code" binding:"required"`
}

type DriverResponse struct {
	ID        int64     `json:"id" binding:"required"`
	Phone     string    `json:"phone" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
}

type DriverLoginPhoneResponse struct {
	AccessToken string         `json:"access_token"`
	Driver      DriverResponse `json:"driver"`
}

func newDriverResponse(driver db.Driver) DriverResponse {
	return DriverResponse{
		ID:        driver.ID,
		Phone:     driver.Phone,
		CreatedAt: driver.CreatedAt,
	}
}

// type OTP struct {
// 	Otp   string `json:"otp" binding:"required"`
// 	Token string `json:"token" binding:"required"`
// }

func (server *Server) driverLoginPhone(ctx *gin.Context) {
	val, err := ctx.GetRawData()
	if err != nil {
		log.Print("Error getting raw data: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var login_request LoginPhoneRequest
	err = json.Unmarshal(val, &login_request)
	if err != nil {
		log.Print("Error unmarshalling phone number: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	driver, err := server.store.GetDriverByPhone(ctx, login_request.Phone)
	if err != nil {
		if err == sql.ErrNoRows {
			// User does not exist, create a new one
			// You need to provide the necessary fields for a new user
			driver, err = server.store.CreateDriver(ctx, login_request.Phone)
			if err != nil {
				log.Print("Error creating driver: ", err)
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		} else {
			log.Print("Error getting driver by phone: ", err)
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	// TODO: temporary otp, replace later with twilio
	otp := "828954"
	server.store.UpdateDriverLoginCode(ctx, db.UpdateDriverLoginCodeParams{ID: driver.ID, LoginCode: sql.NullString{String: otp, Valid: true}})

	// TODO: return session token

	ctx.JSON(http.StatusOK, nil)
	return
}

func (server *Server) verifyDriverLoginPhone(ctx *gin.Context) {
	val, err := ctx.GetRawData()
	if err != nil {
		log.Print("Error getting raw data: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var login_request LoginPhoneRequest
	err = json.Unmarshal(val, &login_request)
	if err != nil {
		log.Print("Error unmarshalling phone number: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	user, err := server.store.GetDriverByPhone(ctx, login_request.Phone)
	if err != nil {
		log.Print("Error getting user by phone: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if user.LoginCode.String != login_request.LoginOTP {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
		return
	}
	accessToken, err := server.tokenMaker.CreateToken(user.Phone, server.config.AccessTokenDurationDriver)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	resp := DriverLoginPhoneResponse{
		AccessToken: accessToken,
		Driver:      newDriverResponse(user),
	}

	// TODO: add bearer token
	ctx.JSON(http.StatusOK, resp)
	return
}

func (server *Server) authDriver(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "Ok")
}

type RegisterRequest struct {
	DriverID     int32  `json:"driver_id" binding:"required"`
	Name         string `json:"name" binding:"required"`
	VehicleType  int32  `json:"vehicle_type" binding:"required"`
	VehicleModel string `json:"vehicle_model" binding:"required"`
	VehicleLabel string `json:"vehicle_label" binding:"required"`
	VehicleColor string `json:"vehicle_color" binding:"required"`
	VehiclePlate string `json:"vehicle_plate" binding:"required"`
}

func (server *Server) driverRegister(ctx *gin.Context) {
	val, err := ctx.GetRawData()
	if err != nil {
		log.Print("Error getting raw data: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var register_request RegisterRequest
	err = json.Unmarshal(val, &register_request)
	if err != nil {
		log.Print("Error unmarshalling register request: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	_, err = server.store.CreateEngagement(ctx, db.CreateEngagementParams{
		DriverID:     register_request.DriverID,
		Name:         register_request.Name,
		VehicleID:    register_request.VehicleType,
		Label:        register_request.VehicleLabel,
		Model:        register_request.VehicleModel,
		Color:        register_request.VehicleColor,
		LicensePlate: register_request.VehiclePlate,
	})
	if err != nil {
		log.Print("Error creating engagement: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
	return
}
