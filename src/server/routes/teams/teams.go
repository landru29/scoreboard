package teams

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/landru29/scoreboard/src/server/database"
	"github.com/spf13/viper"
)

// DefineRoutes define all the teams routes
func DefineRoutes(router *gin.Engine) (teamGroup *gin.RouterGroup, identifiedTeamGroup *gin.RouterGroup) {
	teamGroup = router.Group("/teams")
	{
		/**
		 * List all teams
		 * GET /teams
		 */
		teamGroup.GET("/", func(c *gin.Context) {
			var teams []Team
			rows, err := database.Database.Query("SELECT id, name, color, color_code, logo, created FROM team")
			if database.CheckError(c, err, "Could not read the database (team)") != nil {
				return
			}

			for rows.Next() {
				team := Team{}
				err = rows.Scan(&team.ID, &team.Name, &team.Color, &team.ColorCode, &team.Logo, &team.Created)
				if database.CheckError(c, err, "Could not fetch data from the database (team)") != nil {
					return
				}
				teams = append(teams, team)
			}

			rows.Close()

			c.JSON(http.StatusOK, teams)
		})

		identifiedTeamGroup = teamGroup.Group("/:teamId")
		{
			/**
			 * Detail of a team
			 * GET /teams/:teamId
			 */
			identifiedTeamGroup.GET("/", func(c *gin.Context) {

				id, err := strconv.Atoi(c.Param("teamId"))
				if database.CheckError(c, err, "Bad format of ID") != nil {
					return
				}

				team, err := GetTeamByID(int64(id))
				if database.CheckError(c, err, "Team not found") != nil {
					return
				}

				c.JSON(http.StatusOK, team)

			})

			/**
			 * Upload a logo
			 * POST /teams/:teamId/logo
			 */
			identifiedTeamGroup.POST("/logo", func(c *gin.Context) {
				id, err := strconv.Atoi(c.Param("teamId"))
				if database.CheckError(c, err, "Bad format of ID") != nil {
					return
				}

				team, err := GetTeamByID(int64(id))
				if database.CheckError(c, err, "Could not read the database (team)") != nil {
					return
				}

				file, header, err := c.Request.FormFile("logo")
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"message": "An error occured",
						"error":   err.Error(),
					})
					return
				}
				if file == nil || header == nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"message": "You must specify a file",
						"error":   "No file",
					})
					return
				}

				team.Logo = "/logo/" + header.Filename

				if _, err := os.Stat(viper.GetString("logo_dir")); os.IsNotExist(err) {
					os.Mkdir(viper.GetString("logo_dir"), os.ModePerm)
				}

				filename := viper.GetString("logo_dir") + "/" + header.Filename
				fmt.Printf("Creating file %s\n", filename)

				out, err := os.Create(filename)
				if err != nil {
					log.Fatal(err)
				}
				defer out.Close()
				_, err = io.Copy(out, file)
				if err != nil {
					log.Fatal(err)
				}

				stmt, err := database.Database.Prepare("update team set logo=? where id=?")
				if database.CheckError(c, err, "Mal-formed database query") != nil {
					return
				}

				res, err := stmt.Exec(team.Logo, id)
				if database.CheckError(c, err, "Could not update the team") != nil {
					return
				}

				_, err = res.RowsAffected()
				if database.CheckError(c, err, "Could not check the upload") != nil {
					return
				}

				c.JSON(http.StatusOK, team)

			})

		}

		/**
		 * Create a team
		 * POST /teams
		 */
		teamGroup.POST("/", func(c *gin.Context) {
			team := Team{
				Name:      "",
				Color:     "",
				ColorCode: "",
				Logo:      "",
			}
			err := c.BindJSON(&team)
			if database.CheckError(c, err, "Bad JSON format") != nil {
				return
			}

			team.Created = time.Now().Format("2006-01-02 15:04:05")

			stmt, err := database.Database.Prepare("INSERT INTO team(name, color, color_code, logo, created) values(?,?,?,?,?)")
			if database.CheckError(c, err, "Mal-formed database query") != nil {
				return
			}

			res, err := stmt.Exec(team.Name, team.Color, team.ColorCode, team.Logo, team.Created)
			if database.CheckError(c, err, "Could not create the team in the database") != nil {
				return
			}

			team.ID, err = res.LastInsertId()
			if database.CheckError(c, err, "Could not check the team in the database") != nil {
				return
			}

			c.JSON(http.StatusOK, team)

		})

		/**
		 * delete a team
		 * DELETE /teams
		 */
		teamGroup.DELETE("/:teamId", func(c *gin.Context) {
			id := c.Param("teamId")

			teamDelete, err := database.Database.Prepare("delete from team where id=?")
			if database.CheckError(c, err, "Mal-formed database query") != nil {
				return
			}

			resTeam, err := teamDelete.Exec(id)
			if database.CheckError(c, err, "Could not delete the team from the database") != nil {
				return
			}

			affectTeam, err := resTeam.RowsAffected()
			if database.CheckError(c, err, "Could not check the deletion") != nil {
				return
			}

			playersDelete, err := database.Database.Prepare("delete from player where team=?")
			if database.CheckError(c, err, "Mal-formed database query") != nil {
				return
			}

			resPlayers, err := playersDelete.Exec(id)
			if database.CheckError(c, err, "Could not delete the players from the database") != nil {
				return
			}

			_, err = resPlayers.RowsAffected()
			if database.CheckError(c, err, "Could not check the deletion") != nil {
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"id": affectTeam,
			})

		})

		/**
		 * Update a team
		 * PUT /teams
		 */
		teamGroup.PUT("/:teamId", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("teamId"))
			if database.CheckError(c, err, "Bad format of ID") != nil {
				return
			}

			team, err := GetTeamByID(int64(id))
			if database.CheckError(c, err, "Team not found") != nil {
				return
			}

			err = c.BindJSON(&team)
			if database.CheckError(c, err, "Bad JSON format") != nil {
				return
			}

			stmt, err := database.Database.Prepare("update team set name=?, color=?, color_code=? where id=?")
			if database.CheckError(c, err, "Mal-formed database query") != nil {
				return
			}

			res, err := stmt.Exec(team.Name, team.Color, team.ColorCode, id)
			if database.CheckError(c, err, "Could not update the team") != nil {
				return
			}

			_, err = res.RowsAffected()
			if database.CheckError(c, err, "Could not check the update") != nil {
				return
			}

			c.JSON(http.StatusOK, team)

		})

	}
	return
}
