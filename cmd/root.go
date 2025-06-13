package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	owm "github.com/vekjja/go-openweather"
	"github.com/vekjja/go-toolbox"
)

var verbose int
var lat, lon float64

var rootCmd = &cobra.Command{
	Use:   "wAether",
	Short: "CLI Weather Information 🌤️",
	Long: `
  wAether 🌤️ in your cli.
  All Weather data comes from OpenWeather API.
`,
	Run: func(cmd *cobra.Command, args []string) {

		// Get the geo data
		geoData, err := toolbox.GetGeoData(viper.GetString("location"))
		toolbox.EoE(err, "Error getting GeoData: ")
		location := geoData[0].DisplayName

		lat, err = strconv.ParseFloat(geoData[0].Lat, 64)
		toolbox.EoE(err, "Error converting Latitude: ")

		lon, err = strconv.ParseFloat(geoData[0].Lon, 64)
		toolbox.EoE(err, "Error converting Longitude: ")

		// Get weather data
		weatherData, err := owm.Get(lat, lon, viper.GetString("unit"), viper.GetString("api_key"))
		toolbox.EoE(err, "Error getting Weather Data: ")

		// If vvverbose, show raw JSON data
		if verbose > 2 {
			fmt.Println("Weather Data:")
			j, err := json.Marshal(weatherData)
			toolbox.EoE(err, "Error marshalling JSON: ")
			toolbox.PrettyJson(j)
		}

		// Print weather info
		fmt.Println()
		if len(weatherData.Alerts) > 0 {
			fmt.Println("🚨", weatherData.Alerts[0].Event, "🚨")
			if verbose > 0 {
				fmt.Println(weatherData.Alerts[0].Description)
				fmt.Println()
			}
		}
		fmt.Printf("📍: %s: %s - %s %v\n", location, weatherData.Current.Weather[0].Main, weatherData.Current.Weather[0].Description, owm.Icon(weatherData.Current.Weather[0].Icon))
		fmt.Println("ℹ️ :", weatherData.Daily[0].Summary)
		fmt.Println("⌚️:", toolbox.TimeUTC(weatherData.Current.Dt, weatherData.TimezoneOffset, weatherData.Timezone, ""), owm.MoonPhaseIcon(weatherData.Daily[0].MoonPhase))
		fmt.Printf("🌡️ : %.2f %s", weatherData.Current.Temp, owm.UnitSymbol(viper.GetString("unit")))
		if weatherData.Current.FeelsLike != weatherData.Current.Temp || verbose > 0 {
			fmt.Printf(" (feels like %.2f %s)", weatherData.Current.FeelsLike, owm.UnitSymbol(viper.GetString("unit")))
		}
		fmt.Println()
		if weatherData.Current.WindSpeed > 0 || verbose > 1 {
			fmt.Printf("💨: %.2f %s\n", weatherData.Current.WindSpeed, owm.UnitSpeed(viper.GetString("unit")))
		}
		if weatherData.Current.Rain.OneH > 0 || verbose > 1 {
			fmt.Printf("🌧️ : %.2f mm\n", weatherData.Current.Rain.OneH)
		}
		if weatherData.Current.Snow.OneH > 0 || verbose > 1 {
			fmt.Printf("🌨️ : %.2f mm\n", weatherData.Current.Snow.OneH)
		}
		fmt.Printf("💧: %d%%\n", weatherData.Current.Humidity)
		if weatherData.Current.Visibility < 10000 || verbose > 1 {
			fmt.Printf("👓: %d m\n", weatherData.Current.Visibility)
		}
		timeFormat := "3:04 PM"
		if verbose > 0 {
			timeFormat = "3:04:05 PM MST"
		}
		if verbose > 0 {
			fmt.Printf("🌅: %s\n", toolbox.TimeUTC(weatherData.Current.Sunrise, weatherData.TimezoneOffset, weatherData.Timezone, timeFormat))
			fmt.Printf("🌇: %s\n", toolbox.TimeUTC(weatherData.Current.Sunset, weatherData.TimezoneOffset, weatherData.Timezone, timeFormat))
			if weatherData.Current.Uvi > 3 || verbose > 0 {
				fmt.Printf("🔆: %.2f\n", weatherData.Current.Uvi)
			}
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	initConfig("waether", "config")
	rootCmd.PersistentFlags().CountVarP(&verbose, "verbose", "v", "Increase verbosity (-v, -vv, -vvv)")

	rootCmd.PersistentFlags().StringP("location", "l", viper.GetString("location"), "location to get weather data for")
	viper.BindPFlag("location", rootCmd.PersistentFlags().Lookup("location"))

	rootCmd.PersistentFlags().StringP("unit", "u", "metric", "unit of measure")
	rootCmd.PersistentFlags().SetAnnotation("unit", cobra.BashCompOneRequiredFlag, owm.ValidUnits)
	viper.BindPFlag("unit", rootCmd.PersistentFlags().Lookup("unit"))
}
