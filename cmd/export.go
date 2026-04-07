package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"goincidentcli/internal/incident"
	"goincidentcli/internal/timeline"

	"github.com/spf13/cobra"
)

var exportIncidentID string

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export an incident to a markdown post-mortem report",
	Long: `Generate a pre-filled markdown post-mortem report for an incident.
Reads all events from the incident's timeline and generates a 'post-mortem.md'
file in the incident's directory ready for review.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var inc *incident.Incident
		var err error

		if exportIncidentID != "" {
			inc, err = incident.Get(exportIncidentID)
		} else {
			inc, err = incident.FindMostRecent()
		}

		if err != nil {
			return fmt.Errorf("could not resolve incident for export: %w", err)
		}

		incDir := incident.Dir(inc.ID)
		entries, err := timeline.Load(incDir)
		if err != nil {
			return fmt.Errorf("failed to load timeline for incident %s: %w", inc.ID, err)
		}

		// Prepare data for the template
		type TimelineEntryData struct {
			Time    string
			Author  string
			Message string
			Metrics map[string]string
		}

		type PostMortemData struct {
			ID      string
			Title   string
			Date    string
			Entries []TimelineEntryData
		}

		data := PostMortemData{
			ID:    inc.ID,
			Title: inc.Title,
			Date:  inc.CreatedAt.Format("2006-01-02 15:04:05 -0700"),
		}

		for _, e := range entries {
			data.Entries = append(data.Entries, TimelineEntryData{
				Time:    e.Timestamp.Format("15:04:05"),
				Author:  e.Author,
				Message: e.Message,
				Metrics: e.Metrics,
			})
		}

		// Markdown template
		const tmplBytes = `# Post-Mortem Report: {{ .Title }}

**Incident ID:** {{ .ID }}
**Date:** {{ .Date }}

## Summary
*(Provide a brief overview of the incident, its cause, and the resolution)*

## Impact
*(Detail the impact on systems and users during the outage)*

## Timeline
{{- if .Entries }}
{{- range .Entries }}
* **{{ .Time }}** - {{ .Author }}: {{ .Message }}
{{- range $key, $val := .Metrics }}
  * ` + "`" + `metric: {{ $key }} = {{ $val }}` + "`" + `
{{- end }}
{{- end }}
{{- else }}
*(No timeline entries recorded)*
{{- end }}

## Corrective Actions
* [ ] *(Action item 1)*
* [ ] *(Action item 2)*
`

		t, err := template.New("postmortem").Parse(tmplBytes)
		if err != nil {
			return fmt.Errorf("failed to parse template: %w", err)
		}

		outputPath := filepath.Join(incDir, "post-mortem.md")
		file, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer file.Close()

		if err := t.Execute(file, data); err != nil {
			return fmt.Errorf("failed to execute template: %w", err)
		}

		fmt.Printf("Post-mortem report exported successfully to: %s\n", outputPath)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringVarP(&exportIncidentID, "incident", "i", "", "Incident ID to export (defaults to most recent)")
}
