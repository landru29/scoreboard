package main

import (
	"fmt"

	"github.com/landru29/scoreboard/src/server/database"
	"github.com/landru29/scoreboard/src/server/routes"
	"github.com/landru29/scoreboard/src/server/routes/sockets"
	"github.com/landru29/scoreboard/src/server/ws"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var mainCommand = &cobra.Command{
	Use:   "scoreboard",
	Short: "Roller Derby Scoreboard",
	Long:  "Scoreboard for Roller Derby",
	Run: func(cmd *cobra.Command, args []string) {
		viper.SetEnvPrefix("rd")
		viper.AutomaticEnv()
		viper.SetConfigType("json")
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
		err := viper.ReadInConfig()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Printf("Client files: %s\n", viper.GetString("client_dir"))

		// Application startup here
		_, err = database.OpenDatabase()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		err = database.InitDatabase()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		router := routes.DefineRoutes()
		router.Run(":" + viper.GetString("api_port"))
	},
}

func init() {
	flags := mainCommand.Flags()

	flags.String("api-host", "your-api-host", "API host")
	flags.String("api-port", "3000", "API port")
	flags.String("api-protocol", "http", "API protocol")
	flags.String("sqlite-file", "./database.db", "Database")
	flags.String("logo-dir", "./logos", "Logo folder")
	flags.String("client-dir", "./assets", "Client folder")

	viper.BindPFlag("api_host", flags.Lookup("api-host"))
	viper.BindPFlag("api_port", flags.Lookup("api-port"))
	viper.BindPFlag("api_protocol", flags.Lookup("api-protocol"))
	viper.BindPFlag("sqlite_file", flags.Lookup("sqlite-file"))
	viper.BindPFlag("logo_dir", flags.Lookup("logo-dir"))
	viper.BindPFlag("client_dir", flags.Lookup("client-dir"))

	sockets.Hub = ws.NewHub()
	go sockets.Hub.Run()
}

func main() {
	mainCommand.Execute()
}
