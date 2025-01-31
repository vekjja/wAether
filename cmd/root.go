package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/seemywingz/go-toolbox"
	openWeather "github.com/seemywingz/openWeatherGO"
	"github.com/spf13/cobra"
)

var verbose bool
var location, unit string

var lat, long float64

var rootCmd = &cobra.Command{
	Use:   "wAether",
	Short: "CLI Weather Information 🌤️",
	Long: `
  wAether 🌤️ in your cli.
  All Weather data comes from OpenWeather API.
`,
	Run: func(cmd *cobra.Command, args []string) {

		if location != "" {
			geoData, err := toolbox.GetGeoData(location)
			toolbox.EoE(err, "Error getting GeoData: ")
			location = geoData[0].DisplayName
			lat, err = strconv.ParseFloat(geoData[0].Lat, 64)
			toolbox.EoE(err, "Error converting Latitude: ")
			long, err = strconv.ParseFloat(geoData[0].Lon, 64)
			toolbox.EoE(err, "Error converting Longitude: ")
		} else {
			ipData, err := toolbox.GetIPData()
			toolbox.EoE(err, "Error getting Location from IP: ")
			lat = ipData.Latitude
			long = ipData.Longitude
		}

		fmt.Println()
		fmt.Println("Location:", location)

		if verbose {
			fmt.Println("Latitude:", lat)
			fmt.Println("Longitude:", long)
		}

		weatherData, err := openWeather.Get(lat, long, unit)
		toolbox.EoE(err, "Error getting Weather Data: ")

		// Print the weather data
		fmt.Printf("Temperature: %.2f %s\n", weatherData.Current.Temp, openWeather.GetUnitSymbol(unit))
		fmt.Printf("Feels Like: %.2f %s\n", weatherData.Current.FeelsLike, openWeather.GetUnitSymbol(unit))
		fmt.Println("Pressure:", weatherData.Current.Pressure, "hPa")
		fmt.Println("Humidity:", weatherData.Current.Humidity, "%")
		fmt.Printf("Dew Point: %.2f %s\n", weatherData.Current.DewPoint, openWeather.GetUnitSymbol(unit))
		fmt.Println("UV Index:", weatherData.Current.Uvi)
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
	rootCmd.PersistentFlags().StringVarP(&unit, "unit", "u", "metric", "unit of measurement (metric, imperial, standard)")

}
