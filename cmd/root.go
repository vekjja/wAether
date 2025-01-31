package cmd

import (
	"encoding/json"
	"fmt"
	_ "image/png" // Enable PNG decoding
	"os"
	"strconv"
	"time"

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
		fmt.Printf("📍: %s: %s - %s %v\n", location, weatherData.Current.Weather[0].Main, weatherData.Current.Weather[0].Description, openWeather.GetIconEmoji(weatherData.Current.Weather[0].Icon))
		if verbose {
			fmt.Println("Latitude:", lat)
			fmt.Println("Longitude:", long)
		}

		// Print current time in the location's timezone
		// Go snippet for local time conversion
		utcSeconds := int64(weatherData.Current.Dt) // UTC-based seconds
		offsetSecs := int(weatherData.TimezoneOffset)

		// Create a *time.Location with the correct offset
		loc := time.FixedZone("LocalTime", offsetSecs)

		// Convert the UTC timestamp to “local time” by .In(loc)
		localTime := time.Unix(utcSeconds, 0).In(loc)

		fmt.Println("⌚️:", localTime.Format("Mon, Jan, 2 - 03:04 PM"), openWeather.GetTimeZoneTLA(weatherData.Timezone))
		fmt.Printf("🌡️ : %.2f %s feels like %.2f %s \n", weatherData.Current.Temp, openWeather.GetUnitSymbol(unit), weatherData.Current.FeelsLike, openWeather.GetUnitSymbol(unit))
		fmt.Printf("💨: %.2f %s\n", weatherData.Current.WindSpeed, openWeather.GetUnitSpeed(unit))
		fmt.Printf("💧: %d%%\n", weatherData.Current.Humidity)
		fmt.Printf("👓: %d m\n", weatherData.Current.Visibility)
		fmt.Printf("🌅: %s %s\n", time.Unix(int64(weatherData.Current.Sunrise), 0).In(loc).Format("3:04 PM"), openWeather.GetTimeZoneTLA(weatherData.Timezone))
		fmt.Printf("🌇: %s %s\n", time.Unix(int64(weatherData.Current.Sunset), 0).In(loc).Format("3:04 PM"), openWeather.GetTimeZoneTLA(weatherData.Timezone))
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
