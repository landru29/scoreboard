package teams

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// DefineRoutes define all the teams routes
func DefineRoutes(router *gin.Engine) *gin.RouterGroup {
	group := router.Group("/teams")
	{
		group.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Teams",
			})
		})
	}
	return group
}
