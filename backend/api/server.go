package api

import (
	"log"
	db "uber-replica/db/sqlc"
	"uber-replica/token"
	"uber-replica/util"

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
		log.Fatal("cannot create token maker: ", err)
	}

	server := &Server{config: config, store: store, tokenMaker: tokenMaker}
	router := gin.Default()
	cors_config := cors.DefaultConfig()
	cors_config.AllowOrigins = []string{"http://localhost:8080"}
	cors_config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Access-Control-Allow-Headers"}
	// config.AllowAllOrigins = true

	router.Use(cors.New(cors_config))

	router.POST("/api/login-phone", server.loginPhone)
	router.POST("/api/login-phone/verify", server.verifyLoginPhone)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.GET("/api/auth", server.authUser)
	authRoutes.GET("/api/distance/:departure/:destination", server.getDistance)
	authRoutes.POST("/api/trip/bike", server.createTripBike)

	router.POST("/api/driver/login-phone", server.driverLoginPhone)
	router.POST("/api/driver/login-phone/verify", server.verifyDriverLoginPhone)
	authDriverRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authDriverRoutes.GET("/api/driver/auth", server.authDriver)
	authDriverRoutes.POST("/api/driver/update-engagement", server.driverUpdateEngagement)

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
