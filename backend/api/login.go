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

type LoginPhoneRequest struct {
	Phone    string `json:"phone" binding:"required"`
	LoginOTP string `json:"login_code" binding:"required"`
}

type userResponse struct {
	ID        int64     `json:"id" binding:"required"`
	Phone     string    `json:"phone" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginPhoneResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		ID:        user.ID,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
	}
}

// type OTP struct .value{
// 	Otp   string `json:"otp" binding:"required"`
// 	Token string `json:"token" binding:"required"`
// }

func (server *Server) loginPhone(ctx *gin.Context) {
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

	user, err := server.store.GetUserByPhone(ctx, login_request.Phone)
	if err != nil {
		if err == sql.ErrNoRows {
			// User does not exist, create a new one
			// You need to provide the necessary fields for a new user
			user, err = server.store.CreateUser(ctx, login_request.Phone)
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

	// TODO: temporary otp, replace later with twilio
	otp := "828954"
	server.store.UpdateUserLoginCode(ctx, db.UpdateUserLoginCodeParams{ID: user.ID, LoginCode: sql.NullString{String: otp, Valid: true}})

	// TODO: return session token

	ctx.JSON(http.StatusOK, nil)
	return
}

func (server *Server) verifyLoginPhone(ctx *gin.Context) {
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

	user, err := server.store.GetUserByPhone(ctx, login_request.Phone)
	if err != nil {
		log.Print("Error getting user by phone: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if user.LoginCode.String != login_request.LoginOTP {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
		return
	}
	accessToken, err := server.tokenMaker.CreateToken(user.Phone, server.config.AccessTokenDurationUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	resp := LoginPhoneResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
	}
	// TODO: add bearer token
	ctx.JSON(http.StatusOK, resp)
	return
}

func (server *Server) authUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "Ok")
}
