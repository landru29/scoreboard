package teams

import (
	"errors"

	"github.com/landru29/scoreboard/src/server/database"
)

// Team define to model of a team
type Team struct {
	ID        int64  `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Color     string `json:"color,omitempty"`
	ColorCode string `json:"color_code,omitempty"`
	Logo      string `json:"logo,omitempty"`
	Created   string `json:"created,omitempty"`
}

//GetTeamByID read a team
func GetTeamByID(id int64) (team Team, err error) {
	counter := 0
	rows, err := database.Database.Query("SELECT id, name, color, color_code, logo, created FROM team WHERE id=?", id)

	if err != nil {
		return
	}

	for rows.Next() {
		err = rows.Scan(&team.ID, &team.Name, &team.Color, &team.ColorCode, &team.Logo, &team.Created)
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
