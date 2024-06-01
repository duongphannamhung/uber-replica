package api

import (
	"math"
	db "uber-replica/db/sqlc"

	"github.com/gin-gonic/gin"
)

type DashboardRequest struct {
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
}

type CalculateResponse struct {
	IsIncrease         bool    `json:"is_increase"`
	Value              int     `json:"value"`
	Difference         int     `json:"difference"`
	PercentageIncrease float64 `json:"percentage_increase"`
}

func (server *Server) calculateNewUsers(ctx *gin.Context) {
	var request DashboardRequest
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

	var is_increase bool = new_number >= old_number
	var difference int = int(math.Abs(float64(new_number - old_number)))
	var percentage_increase float64 = math.Abs(float64(new_number-old_number)) / float64(old_number) * 100

	response := CalculateResponse{
		IsIncrease:         is_increase,
		Value:              int(new_number),
		Difference:         difference,
		PercentageIncrease: percentage_increase,
	}

	ctx.JSON(200, response)
}

func (server *Server) calculateTotalRevenue(ctx *gin.Context) {
	var request DashboardRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	number, err := server.store.TotalRevenue(ctx, db.TotalRevenueParams{
		StartDate: request.StartDate,
		EndDate:   request.EndDate,
	})

	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var new_number int64 = number.SumRevenueInPeriod
	var old_number int64 = number.SumRevenuePreviousPeriod

	var is_increase bool = new_number >= old_number
	var difference int = int(math.Abs(float64(new_number - old_number)))
	var percentage_increase float64 = math.Abs(float64(new_number-old_number)) / float64(old_number) * 100

	response := CalculateResponse{
		IsIncrease:         is_increase,
		Value:              int(new_number),
		Difference:         difference,
		PercentageIncrease: percentage_increase,
	}

	ctx.JSON(200, response)
}

func (server *Server) calculateTotalTrip(ctx *gin.Context) {
	var request DashboardRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	number, err := server.store.TotalTrip(ctx, db.TotalTripParams{
		StartDate: request.StartDate,
		EndDate:   request.EndDate,
	})

	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var new_number int64 = number.CountTripInPeriod
	var old_number int64 = number.CountTripPreviousPeriod

	var is_increase bool = new_number >= old_number
	var difference int = int(math.Abs(float64(new_number - old_number)))
	var percentage_increase float64 = math.Abs(float64(new_number-old_number)) / float64(old_number) * 100

	response := CalculateResponse{
		IsIncrease:         is_increase,
		Value:              int(new_number),
		Difference:         difference,
		PercentageIncrease: percentage_increase,
	}

	ctx.JSON(200, response)
}

type RevenueYearResponse struct {
	RowNumber []int64  `json:"row_number"`
	Month     []string `json:"month"`
	Revenue   []int64  `json:"revenue"`
}

func (server *Server) calculateRevenueYear(ctx *gin.Context) {
	number, err := server.store.RevenueYear(ctx)

	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	row_number := []int64{}
	month := []string{}
	revenue := []int64{}

	for i := 0; i < len(number); i++ {
		row_number = append(row_number, number[i].RowNum)
		month = append(month, number[i].Month)
		revenue = append(revenue, number[i].SumRevenue)
	}

	response := RevenueYearResponse{
		RowNumber: row_number,
		Month:     month,
		Revenue:   revenue,
	}

	ctx.JSON(200, response)
}
