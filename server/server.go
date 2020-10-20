/*

 */

package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bitspawngg/tournament-bracket-manager/authentication"

	"github.com/bitspawngg/tournament-bracket-manager/controllers"
	"github.com/bitspawngg/tournament-bracket-manager/models"
	"github.com/bitspawngg/tournament-bracket-manager/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func CreateServer() *http.Server {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	/*
	 configure Logger
	*/
	log := authentication.ConfigureLogger()

	/*
	 configure Database
	*/
	db_type, exists := os.LookupEnv("NEW_DB_TYPE")
	if !exists {
		log.Fatal("missing DB_TYPE environment variable")
	}

	db_path, exists := os.LookupEnv("NEW_DB_PATH")
	if !exists {
		log.Fatal("missing DB_PATH environment variable")
	}
	db := models.NewDB(db_type, db_path)
	if err := db.Connect(); err != nil {
		log.Fatal("db connection failed")
	}

	/*
		Initialize Services
	*/
	ms := services.NewMatchService(
		log,
		db,
	)

	/*
		Initialize Controllers
	*/
	matchController := controllers.NewMatchController(log, ms)

	/*
		Initialize TokenService
	*/
	ts := authentication.NewTokenService(log)

	/*
		Initialize TokenController
	*/
	tokenController := authentication.NewTokenController(log, ts)

	/*
		Initialize gin
	*/
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(CORSMiddleware())

	// health check
	r.GET("/login", tokenController.HandleLogin)
	r.GET("/verifytoken", tokenController.HandleVerify)
	r.GET("/refreshtoken", tokenController.HandleRefreshToken)
	r.GET("/ping", matchController.HandlePing)
	r.POST("/matchschedule", matchController.HandleGetMatchSchedule)
	r.POST("/setresults", matchController.HandleSetMatchResultS)
	r.POST("/setresultc", matchController.HandleSetMatchResultC)
	/*
		Start HTTP Server
	*/
	// initialize server
	addr := fmt.Sprintf("%s:%d", "0.0.0.0", 8080)
	server := makeServer(addr, r)

	// handle graceful shutdown
	go handleGracefulShutdown(server)

	return server
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, X-Auth-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
