package players

import (
	"errors"
	"strconv"

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

//GetPlayerByID read a player
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
