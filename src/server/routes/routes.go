package routes

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/landru29/scoreboard/src/server/routes/games"
	"github.com/landru29/scoreboard/src/server/routes/parameters"
	"github.com/landru29/scoreboard/src/server/routes/players"
	"github.com/landru29/scoreboard/src/server/routes/sockets"
	"github.com/landru29/scoreboard/src/server/routes/teams"
	"github.com/spf13/viper"
)

var notFoundHTML []byte

func notFound() {
	notFoundFile, err := os.Open(viper.GetString("client_dir") + "/404.html")
	if err != nil {
		fmt.Println(err.Error())
	}

	notFoundHTML, err = ioutil.ReadAll(notFoundFile)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// DefineRoutes define all the routes
func DefineRoutes() *gin.Engine {
	router := gin.Default()

	notFound()
	noRoute(router)

	router.Use(static.Serve("/scoreboard", static.LocalFile(viper.GetString("client_dir"), true)))
	router.Use(static.Serve("/logo", static.LocalFile("./logos", true)))

	_, identifiedTeamGroup := teams.DefineRoutes(router)
	players.DefineRoutes(identifiedTeamGroup)
	sockets.DefineRoutes(router)
	games.DefineRoutes(router)
	parameters.DefineRoutes(router)

	return router
}

func noRoute(router *gin.Engine) {
	router.NoRoute(func(c *gin.Context) {
		if c.Request.URL.Path == "/" {
			c.Redirect(http.StatusMovedPermanently, "/scoreboard")
		} else {
			fmt.Fprintf(c.Writer, string(notFoundHTML))
		}

	})
}
