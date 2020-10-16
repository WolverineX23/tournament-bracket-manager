/*

 */

package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bitspawngg/tournament-bracket-manager/controllers"
	"github.com/bitspawngg/tournament-bracket-manager/models"
	"github.com/bitspawngg/tournament-bracket-manager/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func CreateServer() *http.Server {

	/*
	 configure Logger
	*/
	log := logrus.New()
	log.Out = os.Stdout
	log.Level = 4 // Info

	if err := godotenv.Load(); err != nil {
		log.Fatal(err.Error())
	}

	LOG_FILE_LOCATION, exists := os.LookupEnv("LOG_FILE_LOCATION")
	if !exists {
		log.Fatal("missing LOG_FILE_LOCATION environment variable")
	}
	logfile, err := os.OpenFile(LOG_FILE_LOCATION, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("failed to open file for log")
	} else {
		log.Out = logfile
		log.Formatter = &logrus.JSONFormatter{}
	}

	/*
	 configure Database
	*/
	db_type, exists := os.LookupEnv("DB_TYPE")
	if !exists {
		log.Fatal("missing DB_TYPE environment variable")
	}
	db_path, exists := os.LookupEnv("DB_PATH")
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
		Initialize gin
	*/
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(CORSMiddleware())

	// health check
	r.GET("/ping", matchController.HandlePing)
	r.GET("/login", matchController.HandleLogin)
	r.GET("/verifytoken", matchController.HandleVerify)
	r.GET("/refreshtoken", matchController.HandleRefreshToken)
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
