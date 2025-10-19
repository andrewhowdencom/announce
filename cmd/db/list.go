package db

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/andrewhowdencom/announce/internal/datastore"
	"github.com/spf13/cobra"
)

// ListCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all scheduled announcements",
	Long:  `List all scheduled announcements.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create a new datastore
		store, err := datastore.NewStore()
		if err != nil {
			return fmt.Errorf("failed to create datastore: %w", err)
		}
		defer store.Close()

		// List all announcements
		announcements, err := store.ListAnnouncements()
		if err != nil {
			return fmt.Errorf("failed to list announcements: %w", err)
		}

		// Print the announcements in a table
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tChannel\tStatus\tScheduled At\tContent")
		for _, a := range announcements {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%q\n", a.ID, a.ChannelID, a.Status, a.ScheduledAt.Format("2006-01-02 15:04:05"), a.Content)
		}
		w.Flush()

		return nil
	},
}
