package cmd

import (
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate configuration files to the latest version.",
	Long:  `Migrate configuration files to the latest version.`,
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
