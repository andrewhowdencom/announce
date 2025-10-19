package recurring

import (
	"fmt"
	"os"

	"github.com/andrewhowdencom/announce/internal/datastore"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// ListCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all recurring announcements",
	Long:  `List all recurring announcements.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := datastore.NewStore()
		if err != nil {
			return fmt.Errorf("failed to create datastore: %w", err)
		}
		defer store.Close()

		announcements, err := store.ListAnnouncements()
		if err != nil {
			return fmt.Errorf("failed to list announcements: %w", err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.Header([]string{"ID", "Channel ID", "Cron", "Scheduled At", "Status"})

		for _, a := range announcements {
			if a.Recurring {
				table.Append([]string{a.ID, a.ChannelID, a.Cron, a.ScheduledAt.String(), string(a.Status)})
			}
		}

		table.Render()

		return nil
	},
}