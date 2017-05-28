package database

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

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

//InitDatabase initialize the database
func InitDatabase() (err error) {

	teamTable, err := Database.Prepare(`CREATE TABLE IF NOT EXISTS "team" (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"name" VARCHAR(64) NULL,
        "color" VARCHAR(64) NULL,
		"logo" VARCHAR(55) NULL,
        "created" DATETIME NULL
    );`)
	if err != nil {
		return
	}
	_, err = teamTable.Exec()
	if err != nil {
		return
	}

	playerTable, err := Database.Prepare(`CREATE TABLE IF NOT EXISTS "player" (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"team" INTEGER,
		"name" VARCHAR(64) NULL,
        "number" VARCHAR(64) NULL,
        "created" DATETIME NULL
    );`)
	if err != nil {
		return
	}
	_, err = playerTable.Exec()
	if err != nil {
		return
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
