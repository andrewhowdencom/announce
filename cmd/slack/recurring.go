package slack

import (
	"github.com/spf13/cobra"
)

// RecurringCmd represents the recurring command
var RecurringCmd = &cobra.Command{
	Use:   "recurring",
	Short: "Manage recurring calls",
	Long:  `Manage recurring calls.`,
}
