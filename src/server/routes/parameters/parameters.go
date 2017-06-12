package parameters

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

func getParameter() (parameter Parameter, err error) {
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

	if counter == 1 {
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

		return getParameter()
	}

	rows.Close()

	return
}

func createParameter(parameterIn ParameterInput) (parameter Parameter, err error) {
	paramater, err := getParameter()
	if err != nil {
		if err.Error() != "Not found" {
			return
		}

		// Create entry
		stmt, err := database.Database.Prepare("INSERT INTO parameter(game, created) values(?,?)")
		if err != nil {
			return paramater, err
		}

		_, err = stmt.Exec(parameterIn.Game, parameterIn.Created)
		if err != nil {
			return paramater, err
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

	return getParameter()
}

// DefineRoutes define all the games routes
func DefineRoutes(router *gin.Engine) (playerGroup *gin.RouterGroup) {
	playerGroup = router.Group("/parameters")
	{
		/**
		 * List all games
		 * GET /games
		 */
		playerGroup.GET("/", func(c *gin.Context) {

			parameter, err := getParameter()
			if database.CheckError(c, err, "No parameter found") != nil {
				return
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
