package db

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/andrewhowdencom/ruf/internal/datastore"
	"github.com/spf13/cobra"
)

// ListCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all scheduled calls",
	Long:  `List all scheduled calls.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the status flag
		status, _ := cmd.Flags().GetString("status")

		// Create a new datastore
		store, err := datastore.NewStore()
		if err != nil {
			return fmt.Errorf("failed to create datastore: %w", err)
		}
		defer store.Close()

		// List all calls
		calls, err := store.ListCalls()
		if err != nil {
			return fmt.Errorf("failed to list calls: %w", err)
		}

		// Print the calls in a table
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tChannel\tStatus\tScheduled At\tCron\tRecurring\tContent")
		for _, a := range calls {
			if status == "" || a.Status == datastore.Status(status) {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%t\t%q\n", a.ID, a.ChannelID, a.Status, a.ScheduledAt.Format("2006-01-02 15:04:05"), a.Cron, a.Recurring, a.Content)
			}
		}
		w.Flush()

		return nil
	},
}

func init() {
	ListCmd.Flags().String("status", "", "Filter by status")
}