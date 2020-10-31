package socketFunc

import (
	"fmt"

	socketio "github.com/googollee/go-socket.io"
)

func HandleSocketConn(s socketio.Conn) error {
	s.SetContext("")
	fmt.Println("New Connection:", s.ID())
	s.Join("tournament")
	return nil
}

func HandleSocketDisconn(s socketio.Conn, msg string) {
	fmt.Println("closed", msg)
}

func HandleSocketError(s socketio.Conn, err error) {
	fmt.Println("meet error:", err)
}

func HandlePing(s socketio.Conn, msg string) string {
	fmt.Println("Receive ping: " + msg)
	s.Emit("pong", "{ data: 2333 }")
	return "recv" + msg
}
