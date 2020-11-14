/*

 */

package main

import (
	"github.com/bitspawngg/tournament-bracket-manager/server"
	socketio "github.com/googollee/go-socket.io"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	srv := server.CreateServer()

	// listen and serve
	go serve(srv.SocketServer)
	// defer srv.SocketServer.Close()
	// defer server.HandleGracefulShutdown(srv.SocketServer)
	err := http.ListenAndServe(":8080", srv.Router)
	if err == http.ErrServerClosed {
		log.Println("server shutting down gracefully...")
	} else {
		log.Println("unexpected server shutdown...")
		log.Println("ERR: ", err)
	}
}

func serve(server *socketio.Server) {
	if err := server.Serve(); err != nil {
		log.Println("Serving fatal:", err)
	}
}
