package cmd

import (
	"fmt"
	"os"
	"os/user"

	"goincidentcli/internal/incident"
	"goincidentcli/internal/metrics"
	"goincidentcli/internal/timeline"

	"github.com/spf13/cobra"
)

var (
	logMessage    string
	logIncidentID string
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Add or view timeline entries for an incident",
	Long: `Add a timestamped entry to an incident's timeline, or view the full timeline.

When --message is provided, a new entry is appended to the timeline of the most
recent incident (or the one specified with --incident). If Prometheus is configured,
the current value of the monitored metric is automatically captured and attached.

Without --message, the command lists all existing timeline entries.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if logMessage != "" {
			return runLogAdd()
		}
		return runLogList()
	},
}

var logListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all timeline entries for an incident",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runLogList()
	},
}

func init() {
	logCmd.PersistentFlags().StringVarP(&logIncidentID, "incident", "i", "", "Incident ID (defaults to most recent)")
	logCmd.Flags().StringVarP(&logMessage, "message", "m", "", "Message to add to the timeline")

	logCmd.AddCommand(logListCmd)
	RootCmd.AddCommand(logCmd)
}

func resolveIncidentDir() (string, string, error) {
	if logIncidentID != "" {
		return logIncidentID, incident.Dir(logIncidentID), nil
	}
	inc, err := incident.FindMostRecent()
	if err != nil {
		return "", "", fmt.Errorf("could not find an active incident (use --incident to specify one): %w", err)
	}
	return inc.ID, incident.Dir(inc.ID), nil
}

func currentAuthor() string {
	u, err := user.Current()
	if err == nil {
		name := u.Name
		if name == "" {
			name = u.Username
		}
		if name != "" {
			hostname, _ := os.Hostname()
			return fmt.Sprintf("%s@%s", name, hostname)
		}
	}
	hostname, _ := os.Hostname()
	return "unknown@" + hostname
}

func capturePrometheusMetrics() map[string]string {
	if appCfg.PrometheusURL == "" {
		return nil
	}
	query := appCfg.PrometheusQuery
	if query == "" {
		query = `sum(rate(http_requests_total{status=~"5.."}[5m]))`
	}

	client := metrics.NewPrometheusClient(appCfg.PrometheusURL)
	value, err := client.QueryCurrent(query)
	if err != nil {
		fmt.Printf("Warning: failed to capture Prometheus metric: %v\n", err)
		return nil
	}

	return map[string]string{query: value}
}

func runLogAdd() error {
	incID, incDir, err := resolveIncidentDir()
	if err != nil {
		return err
	}

	author := currentAuthor()
	captured := capturePrometheusMetrics()

	entry, err := timeline.AddEntry(incDir, logMessage, author, captured)
	if err != nil {
		return fmt.Errorf("failed to add timeline entry: %w", err)
	}

	fmt.Printf("Timeline entry added to %s\n", incID)
	fmt.Printf("  [%s] %s: %s\n", entry.Timestamp.Format("2006-01-02 15:04:05"), entry.Author, entry.Message)

	if len(entry.Metrics) > 0 {
		fmt.Println("  Metrics captured:")
		for k, v := range entry.Metrics {
			fmt.Printf("    %s = %s\n", k, v)
		}
	}

	return nil
}

func runLogList() error {
	incID, incDir, err := resolveIncidentDir()
	if err != nil {
		return err
	}

	entries, err := timeline.Load(incDir)
	if err != nil {
		return fmt.Errorf("failed to load timeline: %w", err)
	}

	if len(entries) == 0 {
		fmt.Printf("No timeline entries for incident %s.\n", incID)
		return nil
	}

	fmt.Printf("Timeline for %s (%d entries)\n", incID, len(entries))
	fmt.Println("────────────────────────────────────────────────────────────────")
	for _, e := range entries {
		fmt.Printf("[%s] %s\n  %s\n", e.Timestamp.Format("2006-01-02 15:04:05 -07:00"), e.Author, e.Message)
		for k, v := range e.Metrics {
			fmt.Printf("  metric: %s = %s\n", k, v)
		}
	}

	return nil
}
