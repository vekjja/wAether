package cmd

import (
	"fmt"
	"os"

	"github.com/seemywingz/go-toolbox"
	openWeather "github.com/seemywingz/openWeatherGO"
	"github.com/spf13/cobra"
)

var verbose bool
var location string

var lat, long float64

var rootCmd = &cobra.Command{
	Use:   "wAether",
	Short: "CLI Weather Information 🌤️",
	Long: `
  wAether 🌤️ in your cli.
  All Weather data comes from OpenWeather API.
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := openWeather.Init()
		toolbox.EoE(err, "Error Initializing OpenWeatherMap: ")

		if location != "" {
			geoData, err := openWeather.GetGeoData(location)
			toolbox.EoE(err, "Error getting GeoData: ")
			lat = geoData[0].Lat
			long = geoData[0].Lon
		} else {
			ipInfo, err := toolbox.GetIPInfo()
			toolbox.EoE(err, "Error getting Location from IP: ")
			lat = ipInfo.Latitude
			long = ipInfo.Longitude
		}

		if verbose {
			fmt.Println("Location: ", location)
			fmt.Println("Latitude: ", lat)
			fmt.Println("Longitude: ", long)
		}

		openWeather.Now()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&location, "location", "l", "", "location to get weather information")
}
