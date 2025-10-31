package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/andrewhowdencom/ruf/internal/model"
	"github.com/andrewhowdencom/ruf/internal/templater"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var debugCallsCmd = &cobra.Command{
	Use:   "calls",
	Short: "List all scheduled calls from all sources.",
	Long:  `List all scheduled calls from all sources.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := buildSourcer()

		urls := viper.GetStringSlice("source.urls")
		var allCalls []*model.Call

		for _, url := range urls {
			calls, _, err := s.Source(url)
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Error sourcing from %s: %v\n", url, err)
				continue
			}
			allCalls = append(allCalls, calls...)
		}

		for _, call := range allCalls {
			subject, err := templater.Render(call.Subject)
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Error rendering subject for call %s: %v\n", call.ID, err)
			} else {
				call.Subject = subject
			}

			content, err := templater.Render(call.Content)
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Error rendering content for call %s: %v\n", call.ID, err)
			} else {
				call.Content = content
			}
		}

		output, err := json.MarshalIndent(allCalls, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal calls to JSON: %w", err)
		}

		fmt.Fprintln(cmd.OutOrStdout(), string(output))

		return nil
	},
}

func init() {
	debugCmd.AddCommand(debugCallsCmd)
}
