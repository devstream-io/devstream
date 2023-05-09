package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/devstream-io/devstream/internal/log"
	"github.com/devstream-io/devstream/internal/option"
)

var cfgFile string

// isDebug is a flag to enable debug level log
var isDebug bool

// OutputFormat is the output format for the command. One of: json|yaml|raw
// Default value is "raw"
var OutputFormat string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dtm",
	Short: "dtm is a tool to manage variaties of development platforms",
	Long:  `dtm is a tool to manage variaties of development platforms.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initLog()
	},
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
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.devstream.yaml)")
	rootCmd.PersistentFlags().StringVarP(&OutputFormat, "output", "o", "raw", "Output format. One of: json|yaml|raw")
	rootCmd.PersistentFlags().BoolVarP(&isDebug, "debug", "", false, "debug level log")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".devstream" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".devstream")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	if OutputFormat != "raw" {
		option.Silence = true
	}
}

func initLog() {
	// if OutputFormat is not "raw", set log level to PanicLevel to disable log
	if OutputFormat != "raw" {
		logrus.SetLevel(logrus.PanicLevel)
	} else if isDebug {
		logrus.SetLevel(logrus.DebugLevel)
		log.Infof("Log level is: %s.", logrus.GetLevel())
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}
