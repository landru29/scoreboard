package players

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// DefineRoutes define all the teams routes
func DefineRoutes(router *gin.RouterGroup) (playerGroup *gin.RouterGroup) {
	playerGroup = router.Group("/players")
	{
		playerGroup.GET("/", func(c *gin.Context) {
			name := c.Param("teamId")
			c.JSON(http.StatusOK, gin.H{
				"message": "Players",
				"team":    name,
			})
		})
	}
	return
}
