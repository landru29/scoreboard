package games

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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

// DefineRoutes define all the games routes
func DefineRoutes(router *gin.Engine) (playerGroup *gin.RouterGroup) {
	playerGroup = router.Group("/games")
	{
		/**
		 * List all games
		 * GET /games
		 */
		playerGroup.GET("/", func(c *gin.Context) {

			var games []Game
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
			`)

			if err != nil {
				return
			}

			for rows.Next() {
				game := Game{}
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
				if database.CheckError(c, err, "Could not fetch data from the database (game)") != nil {
					return
				}
				games = append(games, game)
			}

			rows.Close()

			c.JSON(http.StatusOK, games)
		})

		/**
		 * Detail of a game
		 * GET /games/:gameId
		 */
		playerGroup.GET("/:gameId", func(c *gin.Context) {

			id, err := strconv.Atoi(c.Param("gameId"))
			if database.CheckError(c, err, "Bad format of ID") != nil {
				return
			}

			game, err := GetGameByID(int64(id))
			if database.CheckError(c, err, "Player not found") != nil {
				return
			}

			c.JSON(http.StatusOK, game)
		})

		/**
		 * Create a game
		 * POST /games
		 */
		playerGroup.POST("/", func(c *gin.Context) {

			game := GameInput{}
			err := c.BindJSON(&game)
			if database.CheckError(c, err, "Bad JSON format") != nil {
				return
			}
			game.Created = time.Now().Format("2006-01-02 15:04:05")

			if database.CheckError(c, checkTeams(game.TeamA, game.TeamB), "Unknown team") != nil {
				return
			}

			stmt, err := database.Database.Prepare("INSERT INTO game(name, teamA, teamB, created) values(?,?,?,?)")
			if database.CheckError(c, err, "Mal-formed database query") != nil {
				return
			}

			res, err := stmt.Exec(game.Name, game.TeamA, game.TeamB, game.Created)
			if database.CheckError(c, err, "Could not create the game in the database") != nil {
				return
			}

			game.ID, err = res.LastInsertId()
			if database.CheckError(c, err, "Could not check the game in the database") != nil {
				return
			}

			c.JSON(http.StatusOK, game)

		})

		/**
		 * Delete a game
		 * DELETE /game/:gameId
		 */
		playerGroup.DELETE("/:gameId", func(c *gin.Context) {
			id := c.Param("gameId")

			stmt, err := database.Database.Prepare("delete from game where id=?")
			if database.CheckError(c, err, "Mal-formed database query") != nil {
				return
			}

			res, err := stmt.Exec(id)
			if database.CheckError(c, err, "Could not delete the game from the database") != nil {
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
		 * Update a game
		 * PUT /games/:gameId
		 */
		playerGroup.PUT("/:gameId", func(c *gin.Context) {

			id, err := strconv.Atoi(c.Param("gameId"))
			if database.CheckError(c, err, "Bad format of ID") != nil {
				return
			}

			game, err := GetSimpleGameByID(int64(id))
			if database.CheckError(c, err, "Game not found") != nil {
				return
			}

			if database.CheckError(c, checkTeams(game.TeamA, game.TeamB), "Unknown team") != nil {
				return
			}

			err = c.BindJSON(&game)
			if database.CheckError(c, err, "Bad JSON format") != nil {
				return
			}

			stmt, err := database.Database.Prepare("update game set name=?, teamA=?, teamB=? where id=?")
			if database.CheckError(c, err, "Mal-formed database query") != nil {
				return
			}

			res, err := stmt.Exec(game.Name, game.TeamA, game.TeamB, id)
			if database.CheckError(c, err, "Could not update the game") != nil {
				return
			}

			_, err = res.RowsAffected()
			if database.CheckError(c, err, "Could not check the update") != nil {
				return
			}

			newGame, err := GetGameByID(int64(id))
			if database.CheckError(c, err, "Could not read the game") != nil {
				return
			}

			c.JSON(http.StatusOK, newGame)

		})
	}
	return
}