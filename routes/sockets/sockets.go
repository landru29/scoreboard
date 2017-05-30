package sockets

import (
	"github.com/gin-gonic/gin"
	"github.com/landru29/scoreboard/ws"
)

// Hub is the current Hub
var Hub *ws.Hub

// DefineRoutes define all the teams routes
func DefineRoutes(router *gin.Engine) (socketGroup *gin.RouterGroup) {
	socketGroup = router.Group("/ws")
	{
		/**
		 * List all teams
		 * GET /teams
		 */
		socketGroup.GET("/", func(c *gin.Context) {
			ws.ServeWs(Hub, c.Writer, c.Request)
		})

	}
	return
}
