package cmd

import (
	"fmt"

	"goincidentcli/internal/incident"
	"goincidentcli/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var statusIncidentID string

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Open the interactive incident dashboard",
	Long: `Opens a live TUI dashboard showing:
  • Elapsed time since the incident was declared
  • The 5 most recent timeline events
  • Health status of configured services

Controls: q / Esc — quit   r — force refresh`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var inc *incident.Incident
		var err error

		if statusIncidentID != "" {
			inc, err = incident.Get(statusIncidentID)
		} else {
			inc, err = incident.FindMostRecent()
		}
		if err != nil {
			return fmt.Errorf("could not find incident: %w", err)
		}

		var svcConfigs []tui.ServiceConfig
		for _, s := range appCfg.Services {
			svcConfigs = append(svcConfigs, tui.ServiceConfig{Name: s.Name, URL: s.URL})
		}

		m := tui.NewModel(inc, incident.Dir(inc.ID), svcConfigs)
		p := tea.NewProgram(m, tea.WithAltScreen())
		_, err = p.Run()
		return err
	},
}

func init() {
	statusCmd.Flags().StringVarP(&statusIncidentID, "incident", "i", "", "Incident ID (defaults to most recent)")
	RootCmd.AddCommand(statusCmd)
}
