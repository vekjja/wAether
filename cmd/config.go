package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	owm "github.com/seemywingz/go-openweather"
	"github.com/seemywingz/go-toolbox" // still using for EoE, etc. if you like
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configName,
	configDir,
	configFile string
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure Custom Defaults",
	Long: `
Configure custom defaults for wAether.`,
	Run: func(cmd *cobra.Command, args []string) {
		configure()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

// initConfig is called early on (from elsewhere) to set up Viper and read config
func initConfig(appName, configName string) {

	configDir = filepath.Join(toolbox.HomeDir(), ".config", appName)
	configFile = filepath.Join(configDir, configName)

	viper.SetConfigType("yaml")
	viper.SetConfigName(configName)
	viper.AddConfigPath(configDir)
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println()
			fmt.Println("It appears you don't have a config file...")
			configure()
			saveConfig()
		}
	}
}

func printCurrentConfig() {
	fmt.Println()
	fmt.Println("Current config values:")
	for _, key := range viper.AllKeys() {
		fmt.Printf("%s: %v\n", key, viper.Get(key))
	}
}

// configAPIKey prompts for the OpenWeather API key if not set, or asks user if they want to replace it
func configAPIKey() {
	if viper.GetString("api_key") == "" {
		fmt.Println()
		fmt.Println("Enter Your OpenWeather API Key")
		fmt.Println("Don't Have One, Get One for Free At: https://openweathermap.org/home/sign_up")
		viper.Set("api_key", toolbox.PromptInput("API Key", viper.GetString("api_key")))
	} else {
		fmt.Println()
		fmt.Println("Replace Your OpenWeather API Key?")
		replace := toolbox.PromptConfirm("Want to Replace Your OpenWeather API Key?")
		if replace {
			viper.Set("api_key", toolbox.PromptInput("API Key", viper.GetString("api_key")))
		}
	}
}

// confirmSave asks if user wants to save current config to file
func confirmSave() {
	printCurrentConfig()
	fmt.Println("\nSave Config?")
	if toolbox.PromptConfirm("\nWant to save these parameters?") {
		saveConfig()
		os.Exit(0)
	}
}

func saveConfig() {
	fmt.Println("\nSaving Config As:", configFile)
	os.MkdirAll(configDir, 0744)
	os.Remove(configFile)
	_, err := os.Create(configFile)
	toolbox.EoE(err, "Error Creating Config File")
	toolbox.EoE(viper.WriteConfig(), "Error Writing Config")
	printCurrentConfig()
}

// configure is the main flow that updates the config interactively
func configure() {
	fmt.Println("\nLet's make some choices:")

	fmt.Println("\nChoose Your Measurement Unit:")
	chosenUnit := toolbox.PromptSelect("Unit", owm.ValidUnits)
	viper.Set("unit", chosenUnit)

	fmt.Println()
	fmt.Println("Enter Your Location:")
	location := toolbox.PromptInput("Location", viper.GetString("location"))
	viper.Set("location", location)

	configAPIKey()
	confirmSave()

	fmt.Println("\nOkay, Great!")
}
