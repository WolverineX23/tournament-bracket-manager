package controllers

import (
	"fmt"
	"github.com/bitspawngg/tournament-bracket-manager/services"
	"github.com/sirupsen/logrus"

	socketio "github.com/googollee/go-socket.io"
)

type SocketController struct {
	log *logrus.Entry
	ss  *services.SocketService
}

func NewSocketController(logger *logrus.Logger, service *services.SocketService) *SocketController {
	return &SocketController{
		log: logger.WithField("controller", "socket"),
		ss:  service,
	}
}

func (sc *SocketController) HandleConnect(s socketio.Conn) error {
	s.SetContext("")
	fmt.Println("New Connection:", s.ID())
	return nil
}

func (sc *SocketController) HandleDisconnect(s socketio.Conn, reason string) {
	fmt.Println("closed", reason)
}

func (sc *SocketController) HandleError(s socketio.Conn, err error) {
	fmt.Println("meet error:", err)
}

func (sc *SocketController) HandlePing(s socketio.Conn, msg string) string {
	fmt.Println("Receive ping: " + msg)
	s.Emit("pong", "pong")
	return "recv" + msg
}

func (sc *SocketController) HandleJoin(s socketio.Conn, msg string) string {
	fmt.Println("Receive Join:", msg)
	s.Join(msg)
	return s.ID() + " Joined " + msg
}

func (sc *SocketController) HandleLeave(s socketio.Conn, msg string) string {
	fmt.Println(s.ID(), " Left from ", msg)
	s.Leave(msg)
	return s.ID() + " Left " + msg
}
