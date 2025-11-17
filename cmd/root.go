package cmd

import (
	"os"

	"github.com/philiplambok/task-api/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "task-backend",
	Short: "A Go application with HTTP server capabilities",
	Long:  `task-backend is a CLI application that provides various commands to manage and run services.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.task-backend.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func newConfig() (internal.Config, error) {
	viper.SetConfigName("env")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return internal.Config{}, err
	}

	var cfg internal.Config
	if err = viper.Unmarshal(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func newDatabase(config internal.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config.Database.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
