/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ruf",
	Short: "A tool to send calls to different platforms.",
	Long: `A tool to send calls to different platforms.

This application is a CLI tool to send calls to different platforms.
Currently, it supports Slack.`,
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $XDG_CONFIG_HOME/ruf/config.yaml)")
	rootCmd.PersistentFlags().String("email.from", "", "Email from address")
	rootCmd.PersistentFlags().String("email.host", "", "Email host")
	rootCmd.PersistentFlags().Int("email.port", 587, "Email port")
	rootCmd.PersistentFlags().String("email.username", "", "Email username")
	rootCmd.PersistentFlags().String("email.password", "", "Email password")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find xdg config path and set it for viper if found.
		configPath, err := xdg.ConfigFile("ruf/config.yaml")
		if err == nil {
			// Search config in the XDG config directory with name "config.yaml".
			viper.AddConfigPath(filepath.Dir(configPath))
			viper.SetConfigName(filepath.Base(configPath))
			viper.SetConfigType("yaml")
		}
	}

	viper.SetEnvPrefix("RUF")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	_ = viper.BindPFlag("email.from", rootCmd.PersistentFlags().Lookup("email.from"))
	_ = viper.BindPFlag("email.host", rootCmd.PersistentFlags().Lookup("email.host"))
	_ = viper.BindPFlag("email.port", rootCmd.PersistentFlags().Lookup("email.port"))
	_ = viper.BindPFlag("email.username", rootCmd.PersistentFlags().Lookup("email.username"))
	_ = viper.BindPFlag("email.password", rootCmd.PersistentFlags().Lookup("email.password"))
}
