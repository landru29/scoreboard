package teams

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/landru29/scoreboard/database"
)

// Team define to model of a team
type Team struct {
	ID      int64  `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Color   string `json:"color,omitempty"`
	Logo    string `json:"logo,omitempty"`
	Created string `json:"created,omitempty"`
}

func teamToGin(team Team) gin.H {
	return gin.H{
		"id":      team.ID,
		"name":    team.Name,
		"color":   team.Color,
		"logo":    team.Logo,
		"created": team.Created,
	}
}

func checkError(c *gin.Context, err error, message string) error {
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": message,
			"error":   err.Error(),
		})
	}
	return err
}

// DefineRoutes define all the teams routes
func DefineRoutes(router *gin.Engine) (teamGroup *gin.RouterGroup, identifiedTeamGroup *gin.RouterGroup) {
	teamGroup = router.Group("/teams")
	{
		teamGroup.GET("/", func(c *gin.Context) {
			var teams []Team
			rows, err := database.Database.Query("SELECT id, name, color, logo, created FROM team")
			if checkError(c, err, "Could not read the database (team)") != nil {
				return
			}

			for rows.Next() {
				team := Team{}
				err = rows.Scan(&team.ID, &team.Name, &team.Color, &team.Logo, &team.Created)
				if checkError(c, err, "Could not fetch data from the database (team)") != nil {
					return
				}
				teams = append(teams, team)
			}

			rows.Close() //good habit to close

			c.JSON(http.StatusOK, teams)
		})

		identifiedTeamGroup = teamGroup.Group("/:teamId")
		{
			identifiedTeamGroup.GET("/", func(c *gin.Context) {
				var teams []Team

				id, err := strconv.Atoi(c.Param("teamId"))
				if checkError(c, err, "Bad format of ID") != nil {
					return
				}

				rows, err := database.Database.Query("SELECT id, name, color, logo, created FROM team WHERE id=?", id)

				if checkError(c, err, "Could not read the database (team)") != nil {
					return
				}

				for rows.Next() {
					team := Team{}
					err = rows.Scan(&team.ID, &team.Name, &team.Color, &team.Logo, &team.Created)
					if checkError(c, err, "Could not fetch data from the database (team)") != nil {
						return
					}
					teams = append(teams, team)
				}

				rows.Close() //good habit to close

				c.JSON(http.StatusOK, teams)

			})

		}

		teamGroup.POST("/", func(c *gin.Context) {
			team := Team{
				Name:  "",
				Color: "",
				Logo:  "",
			}
			err := c.BindJSON(&team)
			if checkError(c, err, "Bad JSON format") != nil {
				return
			}

			team.Created = time.Now().Format("2006-01-02 15:04:05")
			stmt, err := database.Database.Prepare("INSERT INTO team(name, color, logo, created) values(?,?,?,?)")
			if checkError(c, err, "Mal-formed database query") != nil {
				return
			}

			res, err := stmt.Exec(team.Name, team.Color, team.Logo, team.Created)
			if checkError(c, err, "Could not create the team in the database") != nil {
				return
			}

			team.ID, err = res.LastInsertId()
			if checkError(c, err, "Could not check the team in the database") != nil {
				return
			}

			c.JSON(http.StatusOK, team)

		})

		teamGroup.DELETE("/:teamId", func(c *gin.Context) {
			id := c.Param("teamId")

			stmt, err := database.Database.Prepare("delete from team where id=?")
			if checkError(c, err, "Mal-formed database query") != nil {
				return
			}

			res, err := stmt.Exec(id)
			if checkError(c, err, "Could not delete the team from the database") != nil {
				return
			}

			affect, err := res.RowsAffected()
			if checkError(c, err, "Could not check the deletion") != nil {
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"id": affect,
			})

		})

		teamGroup.PUT("/:teamId", func(c *gin.Context) {
			id := c.Param("teamId")

			team := Team{
				Name:  "",
				Color: "",
				Logo:  "",
			}
			err := c.BindJSON(&team)
			if checkError(c, err, "Bad JSON format") != nil {
				return
			}

			stmt, err := database.Database.Prepare("update team set name=? color=? logo=? where id=?")
			if checkError(c, err, "Mal-formed database query") != nil {
				return
			}

			res, err := stmt.Exec(team.Name, team.Color, team.Logo, id)
			if checkError(c, err, "Could not update the team") != nil {
				return
			}

			affect, err := res.RowsAffected()
			if checkError(c, err, "Could not chack the update") != nil {
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"id": affect,
			})

		})

		/*teamGroup.PATCH("/:teamId", func(c *gin.Context) {
			id := c.Param("teamId")

			team := Team{}
			err := c.BindJSON(&team)
			if checkError(c, err, "Bad JSON format") != nil {
				return
			}

		})*/

	}
	return
}
