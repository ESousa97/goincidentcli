package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"goincidentcli/internal/config"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	appCfg  config.Config
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "incident",
	Short: "A simple incident management CLI",
	Long:  `incident is a CLI tool to manage incidents locally and via API.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.incident.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".incident" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".incident")

		cfgPath := filepath.Join(home, ".incident.yaml")

		// Check if config file exists
		if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
			// Create a default empty config file
			viper.Set("api_token", "")
			viper.Set("base_url", "")
			viper.Set("slack_token", "")
			viper.Set("prometheus_url", "")
			viper.Set("prometheus_query", "")
			if err := viper.SafeWriteConfigAs(cfgPath); err != nil {
				fmt.Printf("Error creating config file: %v\n", err)
			} else {
				fmt.Printf("Created template config file at %s\n", cfgPath)
			}
		}
	}

	// Load .env if it exists
	_ = godotenv.Load(".env")

	viper.SetEnvPrefix("INCIDENT")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	_ = viper.BindEnv("api_token")
	_ = viper.BindEnv("base_url")
	_ = viper.BindEnv("slack_token")
	_ = viper.BindEnv("prometheus_url")
	_ = viper.BindEnv("prometheus_query")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	_ = viper.ReadInConfig()

	// Unmarshal into typed struct (from config file and environment)
	if err := viper.Unmarshal(&appCfg); err != nil {
		fmt.Printf("Error unmarshaling config: %v\n", err)
	}
}
