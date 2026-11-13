package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	owm "github.com/vekjja/go-openweather"
	"github.com/vekjja/go-toolbox"
)

var verbose int

var rootCmd = &cobra.Command{
	Use:   "wAether",
	Short: "CLI Weather Information 🌤️",
	Long: `
  wAether 🌤️ in your cli.
  All Weather data comes from OpenWeather API.
`,
	Run: func(cmd *cobra.Command, args []string) {

		if viper.GetString("location") == "" {
			toolbox.EoE(fmt.Errorf("No Location set please run `waether config` or `waether --location \"[Provide a Location]\"`"))
		}

		// Get the geo data
		geoData, err := owm.GetGeoData(viper.GetString("location"), viper.GetString("api_key"))
		if err != nil || len(geoData) == 0 {
			toolbox.EoE(fmt.Errorf("Error Getting Geographical Date for Location: " + viper.GetString("location")))
		}

		// Get weather data
		weatherData, err := owm.Get(geoData[0].Lat, geoData[0].Lon, viper.GetString("unit"), viper.GetString("api_key"))
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
			}
			fmt.Println()
		}
		locationLabel := "📍"
		if verbose > 0 {
			locationLabel = "📍 Location"
		}
		fmt.Printf("%s: %s, %s, %s: %s - %s %v\n", locationLabel, geoData[0].Name, geoData[0].State, geoData[0].Country, weatherData.Current.Weather[0].Main, weatherData.Current.Weather[0].Description, owm.Icon(weatherData.Current.Weather[0].Icon))
		summaryLabel := "ℹ️"
		if verbose > 0 {
			summaryLabel = "ℹ️ Summary"
		}
		fmt.Println(summaryLabel+":", weatherData.Daily[0].Summary)
		timeLabel := "⌚️"
		if verbose > 0 {
			timeLabel = "⌚️ Time"
		}
		fmt.Println(timeLabel+":", toolbox.TimeUTC(weatherData.Current.Dt, weatherData.TimezoneOffset, weatherData.Timezone, ""), owm.MoonPhaseIcon(weatherData.Daily[0].MoonPhase))
		temperatureLabel := "🌡️"
		if verbose > 0 {
			temperatureLabel = "🌡️ Temperature"
		}
		fmt.Println(temperatureLabel+":", weatherData.Current.Temp, owm.UnitSymbol(viper.GetString("unit")))
		feelsLikeLabel := "♥️"
		if verbose > 0 {
			feelsLikeLabel = "♥️ Feels Like"
		}
		fmt.Println(feelsLikeLabel+":", weatherData.Current.FeelsLike, owm.UnitSymbol(viper.GetString("unit")))
		if weatherData.Current.WindSpeed > 0 || verbose > 1 {
			windSpeedLabel := "💨"
			if verbose > 0 {
				windSpeedLabel = "💨 Wind Speed"
			}
			fmt.Println(windSpeedLabel+":", weatherData.Current.WindSpeed, owm.UnitSpeed(viper.GetString("unit")))
		}
		if weatherData.Current.Rain.OneH > 0 || verbose > 1 {
			rainLabel := "🌧️"
			if verbose > 0 {
				rainLabel = "🌧️ Rain"
			}
			fmt.Println(rainLabel+":", weatherData.Current.Rain.OneH, "mm")
		}
		if weatherData.Current.Snow.OneH > 0 || verbose > 1 {
			snowLabel := "🌨️"
			if verbose > 0 {
				snowLabel = "🌨️ Snow"
			}
			fmt.Println(snowLabel+":", weatherData.Current.Snow.OneH, "mm")
		}
		humidityLabel := "💧"
		if verbose > 0 {
			humidityLabel = "💧 Humidity"
		}
		fmt.Println(humidityLabel+":", weatherData.Current.Humidity, "%")
		if weatherData.Current.Visibility < 10000 || verbose > 1 {
			visibilityLabel := "👓"
			if verbose > 0 {
				visibilityLabel = "👓 Visibility"
			}
			fmt.Println(visibilityLabel+":", weatherData.Current.Visibility, "m")
		}
		timeFormat := "3:04 PM"
		if verbose > 0 {
			timeFormat = "3:04:05 PM MST"
		}
		sunriseLabel := "🌅"
		if verbose > 0 {
			sunriseLabel = "🌅 Sunrise"
		}
		fmt.Println(sunriseLabel+":", toolbox.TimeUTC(weatherData.Current.Sunrise, weatherData.TimezoneOffset, weatherData.Timezone, timeFormat))
		sunsetLabel := "🌇"
		if verbose > 0 {
			sunsetLabel = "🌇 Sunset"
		}
		fmt.Println(sunsetLabel+":", toolbox.TimeUTC(weatherData.Current.Sunset, weatherData.TimezoneOffset, weatherData.Timezone, timeFormat))
		if weatherData.Current.Uvi > 3 || verbose > 0 {
			uviLabel := "🔆"
			if verbose > 0 {
				uviLabel = "🔆 UV Index"
			}
			fmt.Println(uviLabel+":", weatherData.Current.Uvi)
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
