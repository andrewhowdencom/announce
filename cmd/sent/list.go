package sent

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
	Short: "List all sent messages",
	Long:  `List all sent messages.`, // Corrected: Removed unnecessary backticks around the string literal
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := datastore.NewStore()
		if err != nil {
			return fmt.Errorf("failed to create datastore: %w", err)
		}
		defer store.Close()

		sentMessages, err := store.ListSentMessages()
		if err != nil {
			return fmt.Errorf("failed to list sent messages: %w", err)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tAnnouncement ID\tTimestamp\tStatus") // Corrected: \t is the correct escape for tab
		for _, sm := range sentMessages {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", sm.ID, sm.AnnouncementID, sm.Timestamp, sm.Status) // Corrected: \t for tab, \n for newline
		}
		w.Flush()

		return nil
	},
}