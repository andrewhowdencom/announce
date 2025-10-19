package slack

import (
	"fmt"

	"github.com/andrewhowdencom/announce/internal/datastore"
	"github.com/spf13/cobra"
)

// DeleteCmd represents the delete command
var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an announcement from the datastore",
	Long: `Delete an announcement from the datastore.

This command deletes an announcement from the datastore given an announcement ID.

Example:
  announce slack delete --id 1`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the ID from the flags
		id, err := cmd.Flags().GetString("id")
		if err != nil {
			return err
		}

		// Create a new datastore
		store, err := datastore.NewStore()
		if err != nil {
			return fmt.Errorf("failed to create datastore: %w", err)
		}
		defer store.Close()

		// Delete the announcement from the datastore
		if err := store.DeleteAnnouncement(id); err != nil {
			return fmt.Errorf("failed to delete announcement: %w", err)
		}

		fmt.Println("Announcement deleted successfully")
		return nil
	},
}

func init() {
	DeleteCmd.Flags().String("id", "", "ID of the announcement to delete")
}
