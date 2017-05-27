package players

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// DefineRoutes define all the teams routes
func DefineRoutes(router *gin.Engine) *gin.RouterGroup {
	group := router.Group("/teams/:team/players")
	{
		group.GET("/", func(c *gin.Context) {
			name := c.Param("team")
			c.JSON(http.StatusOK, gin.H{
				"message": "Players",
				"team":    name,
			})
		})
	}
	return group
}
