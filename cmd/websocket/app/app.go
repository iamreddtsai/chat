package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
)

type App struct {
	HttpHandler http.Handler
}

var (
	_upgrader *websocket.Upgrader
)

func New() *App {
	g := gin.Default()
	v1Group := g.Group("/v1")
	// v1Group.Use(authMiddleware)
	v1Group.GET("/ws", ConnectWS)

	_upgrader = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// todo nats
	nc, _ := nats.Connect(nats.DefaultURL)
	nc.JetStream()

	return &App{
		HttpHandler: g,
	}
}

// @Summary conn websocket
// @Description conn websocket
// @Tags chat
// @Accept json
// @Produce json
// @Success 200 {object} BaseResponse "ok"
// @Failure 400 {object} BaseResponse "bad request"
// @Failure 401 {object} BaseResponse "unauthorized"
// @Failure 403 {object} BaseResponse "forbidden"
// @Failure 500 {object} BaseResponse "server error"
// @Router /v1/ws [get]
func ConnectWS(c *gin.Context) {
	if c.GetHeader("Sec-Websocket-Version") == "" &&
		c.Request.Header.Values("Sec-Websocket-Extensions") == nil &&
		c.GetHeader("Sec-Websocket-Key") == "" {
		c.Status(http.StatusOK)

		return
	}

	conn, err := _upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		// todo log
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	defer conn.Close()

}
