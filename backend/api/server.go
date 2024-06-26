package api

import (
	"log"
	db "uber-replica/db/sqlc"
	"uber-replica/token"
	"uber-replica/util"
	"uber-replica/ws"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		log.Print("cannot create token maker: ", err)
	}

	server := &Server{config: config, store: store, tokenMaker: tokenMaker}
	router := gin.Default()
	cors_config := cors.DefaultConfig()
	cors_config.AllowOrigins = []string{"http://localhost:8080", "http://localhost:3000"}
	cors_config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Access-Control-Allow-Headers"}
	// config.AllowAllOrigins = true

	router.Use(cors.New(cors_config))

	router.POST("/api/login-phone", server.loginPhone)
	router.POST("/api/login-phone/verify", server.verifyLoginPhone)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.GET("/api/auth", server.authUser)
	authRoutes.GET("/api/distance/:departure/:destination", server.getDistance)
	authRoutes.GET("/api/trip/:tripId", server.getTripInfo)
	authRoutes.POST("/api/create-trip", server.createTrip)
	authRoutes.GET("/api/trip/find-driver", server.tripFindDriver)
	authRoutes.GET("/api/trip/get-driver-info/:tripId", server.getDriverInfo)

	router.POST("/api/driver/login-phone", server.driverLoginPhone)
	router.POST("/api/driver/login-phone/verify", server.verifyDriverLoginPhone)
	router.POST("/api/driver/register", server.driverRegister)
	authDriverRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authDriverRoutes.GET("/api/driver/auth", server.authDriver)
	authDriverRoutes.GET("/api/driver/check-engagement", server.checkEngagement)
	authDriverRoutes.GET("/api/driver/current-status", server.currentDriverStatus)
	authDriverRoutes.POST("/api/driver/update-trip-fare", server.updateTripFare)
	authDriverRoutes.POST("/api/driver/update-engagement", server.driverUpdateEngagement)
	authDriverRoutes.POST("/api/driver/finish-engagement", server.finishEngagement)

	router.GET("/api/trip/get-list-trip", server.getListTrip)

	router.POST("/api/dashboard/calculate-new-users", server.calculateNewUsers)
	router.POST("/api/dashboard/calculate-total-revenue", server.calculateTotalRevenue)
	router.POST("/api/dashboard/calculate-total-trips", server.calculateTotalTrip)
	router.GET("/api/dashboard/calculate-revenue-year", server.calculateRevenueYear)
	router.POST("/api/create-bizops-trip", server.createBizopsTrip)
	router.GET("/api/delete-trip/:tripId", server.deleteTrip)
	router.GET("/api/bizops/find-driver", server.tripFindDriver)

	hub := ws.NewHub()
	wsHandler := ws.NewHandler(hub)

	router.POST("/ws/create-room", wsHandler.CreateRoom)
	router.GET("/ws/room-info/:roomId", wsHandler.GetRoomInfo)
	// TODO: add update driver_id to handle driver cancel
	router.GET("/ws/join-room/:roomId", wsHandler.JoinRoom)
	router.GET("/ws/get-clients/:roomId", wsHandler.GetClients)
	go hub.Run()

	server.router = router
	return server, nil
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
