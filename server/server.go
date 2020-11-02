/*

 */

package server

import (
	"fmt"
	"github.com/bitspawngg/tournament-bracket-manager/socketFunc"
	"log"
	"os"

	"github.com/bitspawngg/tournament-bracket-manager/authentication"

	"github.com/bitspawngg/tournament-bracket-manager/controllers"
	"github.com/bitspawngg/tournament-bracket-manager/models"
	"github.com/bitspawngg/tournament-bracket-manager/services"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"github.com/joho/godotenv"
)

type MatchServer struct {
	Router       *gin.Engine
	SocketServer *socketio.Server
}

func CreateServer() *MatchServer {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	/*
	  configure Logger
	*/
	logger := authentication.ConfigureLogger()

	/*
	  configure Database
	*/
	dbType, exists := os.LookupEnv("NEW_DB_TYPE")
	if !exists {
		logger.Fatal("missing DB_TYPE environment variable")
	}

	dbPath, exists := os.LookupEnv("NEW_DB_PATH")
	if !exists {
		logger.Fatal("missing DB_PATH environment variable")
	}
	db := models.NewDB(dbType, dbPath)
	if err := db.Connect(); err != nil {
		logger.Fatal("db connection failed")
	}

	/*
	 Initialize Services
	*/

	/*
	 Initialize Controllers
	*/

	socketServer, err := socketio.NewServer(nil)

	if err != nil {
		logger.Fatal(err)
	}

	ms := services.NewMatchService(logger, db)
	ts := authentication.NewTokenService(logger)

	matchController := controllers.NewMatchController(logger, ms, socketServer)
	tokenController := authentication.NewTokenController(logger, ts)

	/*
	 Initialize gin
	*/
	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(CORSMiddleware())

	socketServer.OnConnect("/", socketFunc.HandleSocketConn)
	socketServer.OnEvent("/", "ping", func(s socketio.Conn) string {
		fmt.Println(s.Context().(string))
		return "23333"
	})
	socketServer.OnError("/", socketFunc.HandleSocketError)
	socketServer.OnDisconnect("/", socketFunc.HandleSocketDisconn)

	// go socket_server.Serve()
	// defer socket_server.Close()

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

	r.GET("/socket.io/", gin.WrapH(socketServer))
	r.POST("/socket.io/", gin.WrapH(socketServer))
	/*
	 Start HTTP Server
	*/
	// initialize server
	// addr := fmt.Sprintf("%s:%d", "0.0.0.0", 8080)
	// server := makeServer(addr, r)

	// handle graceful shutdown
	// go handleGracefulShutdown(socket_server)

	return &MatchServer{Router: r, SocketServer: socketServer}
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
