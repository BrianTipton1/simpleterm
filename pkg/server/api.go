package server

import (
	"net/http"
	"termvim/pkg/term"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// lol this is dumb
func (*Api) Null(api *Api) {
}

type Api struct {
	GetTtyDto     *GetTtyDto
	SetTtySizeDto *SetTtySizeDto
}

func NewApi() *Api {
	return &Api{
		GetTtyDto:     &GetTtyDto{},
		SetTtySizeDto: &SetTtySizeDto{},
	}
}

var router = gin.Default()

type GetTtyDto struct {
	Tty string `json:"tty"`
}

type SetTtySizeDto struct {
	Tty string `json:"tty"`
	Col uint16 `json:"col"`
	Row uint16 `json:"row"`
}

func getTty(ctx *gin.Context) {
	path := term.StartTty()
	name := term.ServerTtyPathToName(path)
	router.GET("/"+name, func(ctx *gin.Context) {
		handleTtySocket(path, ctx)
	})
	ctx.JSON(http.StatusOK, GetTtyDto{Tty: name})
}

func handleTtyResize(ctx *gin.Context) {
	var setTtyReq SetTtySizeDto

	if err := ctx.ShouldBindJSON(&setTtyReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	term.ResizeTty(term.ClientTtyNameToPath(setTtyReq.Tty), setTtyReq.Col, setTtyReq.Row)
}

func StartServer() {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))
	router.GET("/tty", getTty)
	router.POST("/resize", handleTtyResize)
	router.Run()
}
