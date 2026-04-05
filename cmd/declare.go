package cmd

import (
	"fmt"
	"goincidentcli/internal/incident"

	"github.com/spf13/cobra"
)

var title string

// declareCmd represents the declare command
var declareCmd = &cobra.Command{
	Use:   "declare",
	Short: "Declare a new incident",
	Long:  `The declare command initializes a new incident state locally.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if title == "" {
			return fmt.Errorf("the --title flag is required")
		}

		inc, err := incident.Declare(title)
		if err != nil {
			return fmt.Errorf("failed to declare incident: %w", err)
		}

		fmt.Printf("Incident declared successfully!\nID: %s\nTitle: %s\nFolder created: .incidents/%s\n", inc.ID, inc.Title, inc.ID)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(declareCmd)

	declareCmd.Flags().StringVarP(&title, "title", "t", "", "Title of the incident (required)")
	_ = declareCmd.MarkFlagRequired("title")
}
