package cmd

import (
	"fmt"
	"strings"

	"github.com/andrewhowdencom/ruf/internal/sourcer"
	"github.com/andrewhowdencom/ruf/internal/validator"
	"github.com/spf13/cobra"
)

// debugValidateCmd represents the debug validate command
var debugValidateCmd = &cobra.Command{
	Use:   "validate [uri]",
	Short: "Validate a calls file.",
	Long:  `Validate a calls file.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		uri := args[0]

		fetcher := sourcer.NewCompositeFetcher()
		fetcher.AddFetcher("http", sourcer.NewHTTPFetcher())
		fetcher.AddFetcher("https", sourcer.NewHTTPFetcher())
		fetcher.AddFetcher("file", sourcer.NewFileFetcher())
		// Not including git fetcher for now, as it requires more configuration

		parser := sourcer.NewYAMLParser()
		s := sourcer.NewSourcer(fetcher, parser)

		calls, _, err := s.Source(uri)
		if err != nil {
			return err
		}

		errs := validator.Validate(calls)
		if len(errs) > 0 {
			var errStrings []string
			for _, err := range errs {
				errStrings = append(errStrings, err.Error())
			}
			return fmt.Errorf("validation failed:\n%s", strings.Join(errStrings, "\n"))
		}

		fmt.Fprintln(cmd.OutOrStdout(), "OK")
		return nil
	},
}

func init() {
	debugCmd.AddCommand(debugValidateCmd)
}
