package routes

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/landru29/scoreboard/routes/players"
	"github.com/landru29/scoreboard/routes/teams"
)

var notFoundHTML []byte

// DefineRoutes define all the routes
func DefineRoutes() *gin.Engine {
	router := gin.Default()

	noRoute(router)

	router.Use(static.Serve("/scoreboard", static.LocalFile("./assets", true)))

	//router.Static("/scoreboard", "./assets")

	teams.DefineRoutes(router)
	players.DefineRoutes(router)

	return router
}

func init() {
	notFoundFile, err := os.Open("./assets/404.html")
	if err != nil {
		fmt.Println(err.Error())
	}

	notFoundHTML, err = ioutil.ReadAll(notFoundFile)
	if err != nil {
		fmt.Println(err.Error())
	}
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
