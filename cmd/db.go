/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/andrewhowdencom/announce/cmd/db"
	"github.com/spf13/cobra"
)

// dbCmd represents the db command
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Manage the datastore",
	Long:  `Manage the datastore.`,
}

func init() {
	rootCmd.AddCommand(dbCmd)
	dbCmd.AddCommand(db.ListCmd)
}