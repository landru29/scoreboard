package database

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

// EmptyObj define an empty object
type EmptyObj struct {
}

//Database is the public connection to the database
var Database *sql.DB

//OpenDatabase create a nex connection to the database
func OpenDatabase() (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", viper.GetString("sqlite_file"))
	if err == nil {
		Database = db
	}
	return
}

func execSQL(query string) (err error) {
	table, err := Database.Prepare(query)
	if err != nil {
		return
	}
	_, err = table.Exec()
	if err != nil {
		return
	}
	return
}

//InitDatabase initialize the database
func InitDatabase() (err error) {

	if err := execSQL(`
		CREATE TABLE IF NOT EXISTS "team" (
			"id" INTEGER PRIMARY KEY AUTOINCREMENT,
			"name" VARCHAR(64) NULL,
			"color" VARCHAR(64) NULL,
			"color_code" VARCHAR(10) NULL,
			"logo" VARCHAR(55) NULL,
			"created" DATETIME NULL
		);
	`); err != nil {
		return err
	}

	if err := execSQL(`
		CREATE TABLE IF NOT EXISTS "player" (
			"id" INTEGER PRIMARY KEY AUTOINCREMENT,
			"team" INTEGER NOT NULL DEFAULT 0,
			"name" VARCHAR(64) NULL,
			"number" VARCHAR(64) NULL,
			"created" DATETIME NULL
		);
	`); err != nil {
		return err
	}

	if err := execSQL(`
		CREATE TABLE IF NOT EXISTS "game" (
			"id" INTEGER PRIMARY KEY AUTOINCREMENT,
			"start" DATETIME NOT NULL DEFAULT "",
			"end" DATETIME NOT NULL DEFAULT "",
			"period" INTEGER NOT NULL DEFAULT 0,
			"jam" INTEGER NOT NULL DEFAULT 0,
			"scoreA" INTEGER NOT NULL DEFAULT 0,
			"scoreB" INTEGER NOT NULL DEFAULT 0,
			"teamTimeOutA" INTEGER NOT NULL DEFAULT 3,
			"teamTimeOutB" INTEGER NOT NULL DEFAULT 3,
			"officialReviewA" INTEGER NOT NULL DEFAULT 1,
			"officialReviewB" INTEGER NOT NULL DEFAULT 1,
			"name" VARCHAR(64) NULL,
			"teamA" INTEGER NOT NULL DEFAULT 0,
			"teamB" INTEGER NOT NULL DEFAULT 0,
			"created" DATETIME NULL
		);
	`); err != nil {
		return err
	}

	if err := execSQL(`
		CREATE TABLE IF NOT EXISTS "parameter" (
			"id" INTEGER PRIMARY KEY AUTOINCREMENT,
			"game" INTEGER NOT NULL DEFAULT 0,
			"created" DATETIME NULL
		);
	`); err != nil {
		return err
	}

	return
}

//CheckError check errors from database
func CheckError(c *gin.Context, err error, message string) error {
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": message,
			"error":   err.Error(),
		})
	}
	return err
}
