package players

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/landru29/scoreboard/src/server/database"
	"github.com/landru29/scoreboard/src/server/routes/teams"
)

// Player define to model of a player
type Player struct {
	ID      int64  `json:"id,omitempty"`
	Team    int64  `json:"team,omitempty"`
	Name    string `json:"name,omitempty"`
	Number  string `json:"number,omitempty"`
	Created string `json:"created,omitempty"`
}

func getTeam(c *gin.Context) (team teams.Team, err error) {
	id, err := strconv.Atoi(c.Param("teamId"))
	if database.CheckError(c, err, "Bad format of ID") != nil {
		return
	}

	team, err = teams.GetTeamByID(int64(id))
	if database.CheckError(c, err, "Could not find the team") != nil {
		return
	}

	return
}

//GetPlayerByID read a team
func GetPlayerByID(id int64) (player Player, err error) {
	counter := 0
	rows, err := database.Database.Query("SELECT id, name, number, team, created FROM player WHERE id=?", id)

	if err != nil {
		return
	}

	for rows.Next() {
		err = rows.Scan(&player.ID, &player.Name, &player.Number, &player.Team, &player.Created)
		if err != nil {
			return
		}
		counter++
	}

	if counter != 1 {
		err = errors.New("Not found")
		return
	}

	rows.Close()
	return
}

// DefineRoutes define all the teams routes
func DefineRoutes(router *gin.RouterGroup) (playerGroup *gin.RouterGroup) {
	playerGroup = router.Group("/players")
	{
		/**
		 * List all players
		 * GET /teams/:teamId/players
		 */
		playerGroup.GET("/", func(c *gin.Context) {
			team, err := getTeam(c)
			if err != nil {
				return
			}
			var players []Player
			rows, err := database.Database.Query("SELECT id, name, number, team, created FROM player WHERE team=? ORDER BY number ASC", team.ID)
			if database.CheckError(c, err, "Could not read the database (player)") != nil {
				return
			}

			for rows.Next() {
				player := Player{}
				err = rows.Scan(&player.ID, &player.Name, &player.Number, &player.Team, &player.Created)
				if database.CheckError(c, err, "Could not fetch data from the database (teplayeram)") != nil {
					return
				}
				players = append(players, player)
			}

			rows.Close()

			c.JSON(http.StatusOK, players)
		})

		/**
		 * Detail of a player
		 * GET /teams/:teamId/players/:playerId
		 */
		playerGroup.GET("/:playerId", func(c *gin.Context) {
			team, err := getTeam(c)
			if err != nil {
				return
			}

			id, err := strconv.Atoi(c.Param("playerId"))
			if database.CheckError(c, err, "Bad format of ID") != nil {
				return
			}

			player, err := GetPlayerByID(int64(id))
			if database.CheckError(c, err, "Player not found") != nil {
				return
			}

			if player.Team != team.ID {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "Player not found in this team",
					"error":   "Not found",
				})
				return
			}

			c.JSON(http.StatusOK, player)
		})

		/**
		 * Create a player
		 * POST /teams/:teamId/players
		 */
		playerGroup.POST("/", func(c *gin.Context) {
			team, err := getTeam(c)
			if err != nil {
				return
			}

			player := Player{
				Name:   "",
				Number: "",
			}
			err = c.BindJSON(&player)
			if database.CheckError(c, err, "Bad JSON format") != nil {
				return
			}

			player.Team = team.ID
			player.Created = time.Now().Format("2006-01-02 15:04:05")

			stmt, err := database.Database.Prepare("INSERT INTO player(name, number, team, created) values(?,?,?,?)")
			if database.CheckError(c, err, "Mal-formed database query") != nil {
				return
			}

			res, err := stmt.Exec(player.Name, player.Number, player.Team, player.Created)
			if database.CheckError(c, err, "Could not create the player in the database") != nil {
				return
			}

			player.ID, err = res.LastInsertId()
			if database.CheckError(c, err, "Could not check the player in the database") != nil {
				return
			}

			c.JSON(http.StatusOK, player)
		})

		/**
		 * Delete a player
		 * DELETE /teams/:teamId/players/:playerId
		 */
		playerGroup.DELETE("/:playerId", func(c *gin.Context) {
			id := c.Param("teamId")

			team, err := getTeam(c)
			if err != nil {
				return
			}

			stmt, err := database.Database.Prepare("delete from player where id=? AND team=?")
			if database.CheckError(c, err, "Mal-formed database query") != nil {
				return
			}

			res, err := stmt.Exec(id, team.ID)
			if database.CheckError(c, err, "Could not delete the player from the database") != nil {
				return
			}

			affect, err := res.RowsAffected()
			if database.CheckError(c, err, "Could not check the deletion") != nil {
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"id": affect,
			})
		})

		/**
		 * Update a player
		 * PUT /teams/:teamId/players/:playerId
		 */
		playerGroup.PUT("/:playerId", func(c *gin.Context) {
			team, err := getTeam(c)
			if err != nil {
				return
			}

			id, err := strconv.Atoi(c.Param("playerId"))
			if database.CheckError(c, err, "Bad format of ID") != nil {
				return
			}

			player, err := GetPlayerByID(int64(id))
			if database.CheckError(c, err, "Player not found") != nil {
				return
			}

			if player.Team != team.ID {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "Player not found in this team",
					"error":   "Not found",
				})
				return
			}

			err = c.BindJSON(&player)
			if database.CheckError(c, err, "Bad JSON format") != nil {
				return
			}

			stmt, err := database.Database.Prepare("update player set name=?, number=? where id=? AND team=?")
			if database.CheckError(c, err, "Mal-formed database query") != nil {
				return
			}

			res, err := stmt.Exec(player.Name, player.Number, id, team.ID)
			if database.CheckError(c, err, "Could not update the player") != nil {
				return
			}

			_, err = res.RowsAffected()
			if database.CheckError(c, err, "Could not check the update") != nil {
				return
			}

			c.JSON(http.StatusOK, player)

		})
	}
	return
}
