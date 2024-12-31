package server

import (
	"flag"
	"log"
	"net/http"
	"termvim/pkg/term"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	port     = "8080"
	host     = "localhost"
	portHost = host + ":" + port
)

var (
	addr     = flag.String("addr", portHost, "http service address")
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func handleTtySocket(tty string, ctx *gin.Context) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	connection, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Print("upgrade err:", err)
		return
	}
	term.HandleConnection(connection, tty)
}
