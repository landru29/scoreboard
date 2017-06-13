package parameters

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/landru29/scoreboard/src/server/database"
)

// DefineRoutes define all the games routes
func DefineRoutes(router *gin.Engine) (playerGroup *gin.RouterGroup) {
	playerGroup = router.Group("/parameters")
	{
		/**
		 * List all games
		 * GET /games
		 */
		playerGroup.GET("/", func(c *gin.Context) {

			parameter, err := GetParameter()
			if err != nil {
				if err.Error() == "Not found" {
					c.JSON(http.StatusOK, database.EmptyObj{})
					return
				}
				if database.CheckError(c, err, "No parameter found") != nil {
					return
				}
			}

			c.JSON(http.StatusOK, parameter)
		})

		/**
		 * Create a game
		 * POST /games
		 */
		playerGroup.POST("/", func(c *gin.Context) {
			parameter := ParameterInput{}
			err := c.BindJSON(&parameter)
			if database.CheckError(c, err, "Bad JSON format") != nil {
				return
			}
			parameter.Created = time.Now().Format("2006-01-02 15:04:05")

			newParameter, err := createParameter(parameter)
			if database.CheckError(c, err, "Parameter not updated") != nil {
				return
			}

			c.JSON(http.StatusOK, newParameter)

		})

	}
	return
}
