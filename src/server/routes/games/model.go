package games

import (
	"errors"

	"github.com/landru29/scoreboard/src/server/database"
	"github.com/landru29/scoreboard/src/server/routes/teams"
)

// Game define to model of a game
type Game struct {
	ID              int64      `json:"id,omitempty"`
	Start           string     `json:"start,omitempty"`
	End             string     `json:"end,omitempty"`
	Name            string     `json:"name,omitempty"`
	Period          int64      `json:"period,omitempty"`
	Jam             int64      `json:"jam,omitempty"`
	ScoreA          int64      `json:"scoreA,omitempty"`
	ScoreB          int64      `json:"scoreB,omitempty"`
	TeamTimeOutA    int64      `json:"teamTimeOutA,omitempty"`
	TeamTimeOutB    int64      `json:"teamTimeOutB,omitempty"`
	OfficialReviewA int64      `json:"officialReviewA,omitempty"`
	OfficialReviewB int64      `json:"officialReviewB,omitempty"`
	TeamA           teams.Team `json:"teamA,omitempty"`
	TeamB           teams.Team `json:"teamB,omitempty"`
	Created         string     `json:"created,omitempty"`
}

// GameInput define to model of a game for input
type GameInput struct {
	ID              int64  `json:"id,omitempty"`
	Start           string `json:"start,omitempty"`
	End             string `json:"end,omitempty"`
	Name            string `json:"name,omitempty"`
	Period          int64  `json:"period,omitempty"`
	Jam             int64  `json:"jam,omitempty"`
	ScoreA          int64  `json:"scoreA,omitempty"`
	ScoreB          int64  `json:"scoreB,omitempty"`
	TeamTimeOutA    int64  `json:"teamTimeOutA,omitempty"`
	TeamTimeOutB    int64  `json:"teamTimeOutB,omitempty"`
	OfficialReviewA int64  `json:"officialReviewA,omitempty"`
	OfficialReviewB int64  `json:"officialReviewB,omitempty"`
	TeamA           int64  `json:"teamA,omitempty"`
	TeamB           int64  `json:"teamB,omitempty"`
	Created         string `json:"created,omitempty"`
}

func checkTeams(idA int64, idB int64) (err error) {
	_, err = teams.GetTeamByID(idA)
	if err != nil {
		return
	}

	_, err = teams.GetTeamByID(idB)
	if err != nil {
		return
	}

	return
}

//GetGameByID read a game
func GetGameByID(id int64) (game Game, err error) {
	counter := 0
	rows, err := database.Database.Query(`
		SELECT 
		    g.id AS id,
			g.name AS name,
			g.start AS start,
			g.end AS end,
			g.period AS period,
			g.jam AS jam,
			g.scoreA AS scoreA,
			g.scoreB AS scoreB,
			g.teamTimeOutA AS teamTimeOutA,
			g.teamTimeOutB AS teamTimeOutB,
			g.officialReviewA AS officialReviewA,
			g.officialReviewB AS officialReviewB,
			g.created AS created,
			a.id AS teamA_id,
			a.name AS teamA_name,
			a.color AS teamA_color,
			a.color_code AS teamA_color_code,
			a.logo AS teamA_logo,
			a.created AS teamA_created,
			b.id AS teamB_id,
			b.name AS teamB_name,
			b.color AS teamB_color,
			b.color_code AS teamB_color_code,
			b.logo AS teamB_logo,
			b.created AS teamB_created
		FROM game AS g
		INNER JOIN team AS a ON g.teamA = a.id
		INNER JOIN team AS b ON  g.teamB = b.id
		WHERE g.id=?
	`, id)

	if err != nil {
		return
	}

	for rows.Next() {
		err = rows.Scan(
			&game.ID,
			&game.Name,
			&game.Start,
			&game.End,
			&game.Period,
			&game.Jam,
			&game.ScoreA,
			&game.ScoreB,
			&game.TeamTimeOutA,
			&game.TeamTimeOutB,
			&game.OfficialReviewA,
			&game.OfficialReviewB,
			&game.Created,
			&game.TeamA.ID,
			&game.TeamA.Name,
			&game.TeamA.Color,
			&game.TeamA.ColorCode,
			&game.TeamA.Logo,
			&game.TeamA.Created,
			&game.TeamB.ID,
			&game.TeamB.Name,
			&game.TeamB.Color,
			&game.TeamB.ColorCode,
			&game.TeamB.Logo,
			&game.TeamB.Created,
		)
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

//GetSimpleGameByID read a game
func GetSimpleGameByID(id int64) (game GameInput, err error) {
	counter := 0
	rows, err := database.Database.Query(`
		SELECT
			id,
			name,
			teamA,
			teamB,
			start,
			end,
			period,
			jam,
			scoreA,
			scoreB,
			teamTimeOutA,
			teamTimeOutB,
			officialReviewA,
			officialReviewB,
			created
		FROM game
		WHERE id=?
	`, id)

	if err != nil {
		return
	}

	for rows.Next() {
		err = rows.Scan(
			&game.ID,
			&game.Name,
			&game.Start,
			&game.End,
			&game.Period,
			&game.Jam,
			&game.ScoreA,
			&game.ScoreB,
			&game.TeamTimeOutA,
			&game.TeamTimeOutB,
			&game.OfficialReviewA,
			&game.OfficialReviewB,
			&game.TeamA,
			&game.TeamB,
			&game.Created,
		)
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
