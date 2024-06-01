package api

import (
	"fmt"
	"math"
	db "uber-replica/db/sqlc"

	"github.com/gin-gonic/gin"
)

type CalculateNewUsersRequest struct {
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
}

type CalculateNewUsersResponse struct {
	IsIncrease         bool    `json:"is_increase"`
	NewUsers           int     `json:"new_users"`
	Difference         int     `json:"difference"`
	PercentageIncrease float64 `json:"percentage_increase"`
}

func (server *Server) calculateNewUsers(ctx *gin.Context) {
	var request CalculateNewUsersRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	count_new_users, err := server.store.CountNewUsers(ctx, db.CountNewUsersParams{
		StartDate: request.StartDate,
		EndDate:   request.EndDate,
	})

	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var new_number int64 = count_new_users.CountUsersInPeriod
	var old_number int64 = count_new_users.CountUsersPreviousPeriod

	fmt.Println(new_number, old_number)

	var is_increase bool = new_number >= old_number
	var difference int = int(math.Abs(float64(new_number - old_number)))
	var percentage_increase float64 = math.Abs(float64(new_number-old_number)) / float64(old_number) * 100

	response := CalculateNewUsersResponse{
		IsIncrease:         is_increase,
		NewUsers:           int(new_number),
		Difference:         difference,
		PercentageIncrease: percentage_increase,
	}

	ctx.JSON(200, response)
}
