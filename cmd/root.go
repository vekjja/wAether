package cmd

import (
	"encoding/json"
	"fmt"
	_ "image/png" // Enable PNG decoding
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
		// Validate unit
		if _, ok := openWeather.ValidUnits[unit]; !ok {
			toolbox.EoE(fmt.Errorf("Invalid unit: %s", unit), "Error: ")
		}

		// Get the geo data
		geoData, err := toolbox.GetGeoData(location)
		toolbox.EoE(err, "Error getting GeoData: ")
		location = geoData[0].DisplayName

		lat, err = strconv.ParseFloat(geoData[0].Lat, 64)
		toolbox.EoE(err, "Error converting Latitude: ")

		long, err = strconv.ParseFloat(geoData[0].Lon, 64)
		toolbox.EoE(err, "Error converting Longitude: ")

		// Get weather data
		weatherData, err := openWeather.Get(lat, long, unit)
		toolbox.EoE(err, "Error getting Weather Data: ")

		// If verbose, show raw JSON data
		if verbose {
			fmt.Println("Weather Data:")
			j, err := json.Marshal(weatherData)
			toolbox.EoE(err, "Error marshalling JSON: ")
			toolbox.PrettyJson(j)
		}

		// Print standard weather info
		fmt.Println()
		fmt.Println("📍:", location, openWeather.GetIconEmoji(weatherData.Current.Weather[0].Icon))
		if verbose {
			fmt.Println("Latitude:", lat)
			fmt.Println("Longitude:", long)
		}

		fmt.Println("⌚️:", toolbox.FormatTime(int64(weatherData.Current.Dt)))
		fmt.Printf("🌡️ : %.2f %s feels like %.2f %s \n", weatherData.Current.Temp, openWeather.GetUnitSymbol(unit), weatherData.Current.FeelsLike, openWeather.GetUnitSymbol(unit))
		fmt.Printf("💨: %.2f m/s\n", weatherData.Current.WindSpeed)
		fmt.Printf("💧: %d%%\n", weatherData.Current.Humidity)
		fmt.Printf("👓: %d m\n", weatherData.Current.Visibility)
		fmt.Printf("🌅: %s\n", toolbox.FormatTime(int64(weatherData.Current.Sunrise)))
		fmt.Printf("🌇: %s\n", toolbox.FormatTime(int64(weatherData.Current.Sunset)))
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

	rootCmd.MarkPersistentFlagRequired("location")
}
