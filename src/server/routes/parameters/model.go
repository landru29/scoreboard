package parameters

import (
	"errors"

	"github.com/landru29/scoreboard/src/server/database"
	"github.com/landru29/scoreboard/src/server/routes/games"
)

// Parameter define to model of a game
type Parameter struct {
	ID      int64      `json:"id,omitempty"`
	Game    games.Game `json:"game,omitempty"`
	Created string     `json:"created,omitempty"`
}

// ParameterInput define to model of a game for input
type ParameterInput struct {
	ID      int64  `json:"id,omitempty"`
	Game    int64  `json:"game,omitempty"`
	Created string `json:"created,omitempty"`
}

func checkGame(id int64) (err error) {
	_, err = games.GetGameByID(id)
	if err != nil {
		return
	}

	return
}

// GetParameter gets the game parameters
func GetParameter() (parameter Parameter, err error) {
	counter := 0
	firstID := int64(-1)
	rows, err := database.Database.Query(`
		SELECT 
		    p.id AS id,
			p.created AS created,
			g.id AS game_id,
			g.name AS game_name,
			g.start AS game_start,
			g.end AS game_end,
			g.period AS game_period,
			g.jam AS game_jam,
			g.scoreA AS game_scoreA,
			g.scoreB AS game_scoreB,
			g.teamTimeOutA AS game_teamTimeOutA,
			g.teamTimeOutB AS game_teamTimeOutB,
			g.officialReviewA AS game_officialReviewA,
			g.officialReviewB AS game_officialReviewB,
			g.created AS game_created,
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
		FROM parameter AS p
		INNER JOIN game AS g ON p.game = g.id
		INNER JOIN team AS a ON g.teamA = a.id
		INNER JOIN team AS b ON  g.teamB = b.id
	`)

	if err != nil {
		return
	}

	for rows.Next() {
		err = rows.Scan(
			&parameter.ID,
			&parameter.Created,
			&parameter.Game.ID,
			&parameter.Game.Name,
			&parameter.Game.Start,
			&parameter.Game.End,
			&parameter.Game.Period,
			&parameter.Game.Jam,
			&parameter.Game.ScoreA,
			&parameter.Game.ScoreB,
			&parameter.Game.TeamTimeOutA,
			&parameter.Game.TeamTimeOutB,
			&parameter.Game.OfficialReviewA,
			&parameter.Game.OfficialReviewB,
			&parameter.Game.Created,
			&parameter.Game.TeamA.ID,
			&parameter.Game.TeamA.Name,
			&parameter.Game.TeamA.Color,
			&parameter.Game.TeamA.ColorCode,
			&parameter.Game.TeamA.Logo,
			&parameter.Game.TeamA.Created,
			&parameter.Game.TeamB.ID,
			&parameter.Game.TeamB.Name,
			&parameter.Game.TeamB.Color,
			&parameter.Game.TeamB.ColorCode,
			&parameter.Game.TeamB.Logo,
			&parameter.Game.TeamB.Created,
		)
		if counter == 0 {
			firstID = parameter.ID
		}
		if err != nil {
			return
		}
		counter++
	}

	if counter == 0 {
		err = errors.New("Not found")
	}

	if counter > 1 {
		stmt, err := database.Database.Prepare("delete from parameter where id<>?")
		if err != nil {
			err = errors.New("Mal-formed database query")
			return parameter, err
		}

		_, err = stmt.Exec(firstID)
		if err != nil {
			err = errors.New("Could not clean parameters from the database")
			return parameter, err
		}

		return GetParameter()
	}

	rows.Close()

	return
}

func createParameter(parameterIn ParameterInput) (parameter Parameter, err error) {
	parameter, err = GetParameter()
	if err != nil {
		if err.Error() != "Not found" {
			return
		}

		// Create entry
		stmt, err := database.Database.Prepare("INSERT INTO parameter(game, created) values(?,?)")
		if err != nil {
			return parameter, err
		}

		_, err = stmt.Exec(parameterIn.Game, parameterIn.Created)
		if err != nil {
			return parameter, err
		}
	}

	// update entry
	stmt, err := database.Database.Prepare("update parameter set game=?, created=? where id=?")
	if err != nil {
		return
	}

	_, err = stmt.Exec(parameterIn.Game, parameterIn.Created, parameter.ID)
	if err != nil {
		return
	}

	return GetParameter()
}
