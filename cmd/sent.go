package cmd

import (
	"github.com/andrewhowdencom/announce/cmd/sent"
	"github.com/spf13/cobra"
)

// sentCmd represents the sent command
var sentCmd = &cobra.Command{
	Use:   "sent",
	Short: "Interact with sent messages",
	Long:  `Interact with sent messages.`,
}

func init() {
	rootCmd.AddCommand(sentCmd)
	sentCmd.AddCommand(sent.ListCmd)
}