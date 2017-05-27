package main

import (
	"fmt"

	"github.com/landru29/scoreboard/routes"
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
		}

		// Application startup here
		router := routes.DefineRoutes()
		router.Run(":" + viper.GetString("api_port"))
	},
}

func init() {
	flags := mainCommand.Flags()

	flags.String("api-host", "your-api-host", "API host")
	flags.String("api-port", "3000", "API port")
	flags.String("api-protocol", "http", "API protocol")

	viper.BindPFlag("api_host", flags.Lookup("api-host"))
	viper.BindPFlag("api_port", flags.Lookup("api-port"))
	viper.BindPFlag("api_protocol", flags.Lookup("api-protocol"))
}

func main() {
	mainCommand.Execute()
}
