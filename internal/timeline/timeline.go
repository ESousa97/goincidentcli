// Package timeline handles the creation and loading of chronological event records for an incident.
package timeline

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

const timelineFile = "timeline.json"

// Entry represents a single event in an incident's timeline.
type Entry struct {
	ID        string            `json:"id"`
	Timestamp time.Time         `json:"timestamp"`
	Author    string            `json:"author"`
	Message   string            `json:"message"`
	Metrics   map[string]string `json:"metrics,omitempty"`
}

func filePath(incidentDir string) string {
	return filepath.Join(incidentDir, timelineFile)
}

// Load reads all timeline entries from an incident directory.
func Load(incidentDir string) ([]Entry, error) {
	data, err := os.ReadFile(filePath(incidentDir))
	if os.IsNotExist(err) {
		return []Entry{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read timeline: %w", err)
	}
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("failed to parse timeline: %w", err)
	}
	return entries, nil
}

func save(incidentDir string, entries []Entry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal timeline: %w", err)
	}
	return os.WriteFile(filePath(incidentDir), data, 0644)
}

// AddEntry appends a new event to the incident timeline and persists it.
func AddEntry(incidentDir string, message string, author string, metrics map[string]string) (*Entry, error) {
	entries, err := Load(incidentDir)
	if err != nil {
		return nil, err
	}

	id, err := gonanoid.New(8)
	if err != nil {
		return nil, fmt.Errorf("failed to generate entry ID: %w", err)
	}

	entry := Entry{
		ID:        id,
		Timestamp: time.Now(),
		Author:    author,
		Message:   message,
		Metrics:   metrics,
	}

	entries = append(entries, entry)

	if err := save(incidentDir, entries); err != nil {
		return nil, err
	}

	return &entry, nil
}
