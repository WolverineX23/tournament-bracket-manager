/*

 */

package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bitspawngg/tournament-bracket-manager/socketFunc"

	"github.com/bitspawngg/tournament-bracket-manager/authentication"

	"github.com/bitspawngg/tournament-bracket-manager/controllers"
	"github.com/bitspawngg/tournament-bracket-manager/models"
	"github.com/bitspawngg/tournament-bracket-manager/services"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
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

	socket_server, err := socketio.NewServer(nil)

	if err != nil {
		log.Fatal(err)
	}

	matchController := controllers.NewMatchController(log, ms, socket_server)

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

	socket_server.OnConnect("/", socketFunc.HandleSocketConn)

	socket_server.OnEvent("/", "ping", socketFunc.HandlePing)

	socket_server.OnError("/", socketFunc.HandleSocketError)

	socket_server.OnDisconnect("/", socketFunc.HandleSocketDisconn)

	// health check
	r.POST("/login", tokenController.HandleLogin)
	r.GET("/verifytoken", tokenController.HandleVerify)
	r.GET("/refreshtoken", tokenController.HandleRefreshToken)
	r.POST("/refreshtable", matchController.HandleRefreshTable)
	r.GET("/ping", matchController.HandlePing)
	r.POST("/matchschedule", matchController.HandleGetMatchSchedule)
	r.POST("/setresults", matchController.HandleSetMatchResultS)
	r.POST("/setresultc", matchController.HandleSetMatchResultC)
	r.POST("/getalltournamentid", matchController.HandleGetAlltournamentID)
	r.POST("/getrate", matchController.HandleGetRate)

	r.GET("/socket.io/*any", gin.WrapH(socket_server))
	r.POST("/socket.io/*any", gin.WrapH(socket_server))
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
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, X-Auth-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Request.Header.Del("Origin")

		c.Next()
	}
}
