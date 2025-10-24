package cmd

import (
	"github.com/spf13/cobra"
)

// sentCmd represents the sent command
var sentCmd = &cobra.Command{
	Use:   "sent",
	Short: "Interact with sent calls.",
	Long:  `Interact with sent calls.`,
}

func init() {
	rootCmd.AddCommand(sentCmd)
}
