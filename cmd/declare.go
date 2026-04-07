package cmd

import (
	"fmt"
	"goincidentcli/internal/incident"
	"goincidentcli/internal/slack"
	"goincidentcli/internal/timeline"

	"github.com/spf13/cobra"
)

var (
	title    string
	severity string
)

// declareCmd represents the declare command
var declareCmd = &cobra.Command{
	Use:   "declare",
	Short: "Declare a new incident",
	Long:  `The declare command initializes a new incident state locally and optionally creates a Slack channel.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if title == "" {
			return fmt.Errorf("the --title flag is required")
		}

		validSeverities := map[string]bool{"SEV1": true, "SEV2": true, "SEV3": true}
		if !validSeverities[severity] {
			return fmt.Errorf("invalid severity %q: must be SEV1, SEV2, or SEV3", severity)
		}

		inc, err := incident.Declare(title, severity)
		if err != nil {
			return fmt.Errorf("failed to declare incident: %w", err)
		}

		fmt.Printf("Incident declared locally!\nID: %s\nTitle: %s\nSeverity: %s\nFolder created: .incidents/%s\n", inc.ID, inc.Title, inc.Severity, inc.ID)

		// Initial timeline entry with optional Prometheus snapshot
		prometheusMetrics := capturePrometheusMetrics()
		initialMessage := fmt.Sprintf("Incident declared: %s", title)
		if _, err := timeline.AddEntry(incident.Dir(inc.ID), initialMessage, currentAuthor(), prometheusMetrics); err != nil {
			fmt.Printf("Warning: failed to create initial timeline entry: %v\n", err)
		} else if len(prometheusMetrics) > 0 {
			fmt.Println("Prometheus metrics captured at declaration time.")
		}

		// Slack integration (Fault-tolerant)
		if appCfg.SlackToken != "" {
			fmt.Println("Slack Integration: Initializing...")
			slackSvc := slack.NewClient(appCfg.SlackToken)

			// inc-{date}-{sanitized_title}
			channelName := fmt.Sprintf("inc-%s-%s", inc.CreatedAt.Format("2006-01-02"), title)

			channelID, err := slackSvc.CreateIncidentChannel(channelName)
			if err != nil {
				fmt.Printf("Warning: Initialized incident locally, but Slack integration failed: %v\n", err)
			} else {
				// Purpose
				purpose := fmt.Sprintf("Incident: %s (%s)", title, inc.ID)
				if err := slackSvc.SetChannelPurpose(channelID, purpose); err != nil {
					fmt.Printf("Warning: Failed to set Slack channel purpose: %v\n", err)
				}

				// Initial Message
				if err := slackSvc.PostInitialMessage(channelID, inc.ID, title); err != nil {
					fmt.Printf("Warning: Failed to post initial message to Slack: %v\n", err)
				}

				// Update local state with channel ID
				if err := incident.UpdateSlackChannel(inc.ID, channelID); err != nil {
					fmt.Printf("Warning: Failed to update local metadata with Slack channel ID: %v\n", err)
				} else {
					fmt.Printf("Slack channel created successfully: #%s\n", slack.Slugify(channelName))
				}
			}
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(declareCmd)

	declareCmd.Flags().StringVarP(&title, "title", "t", "", "Title of the incident (required)")
	declareCmd.Flags().StringVarP(&severity, "severity", "s", "SEV3", "Incident severity: SEV1, SEV2, or SEV3")
	_ = declareCmd.MarkFlagRequired("title")
}
