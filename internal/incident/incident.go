package incident

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

const (
	incidentsDir = ".incidents"
	dotGitIgnore = ".gitignore"
)

// Incident represents the state of an incident.
type Incident struct {
	ID    string
	Title string
}

// Declare initializes a new incident folder and its state.
func Declare(title string) (*Incident, error) {
	// Generate ID: INC-YYYYMMDD-NanoID (4 chars)
	suffix, err := gonanoid.New(4)
	if err != nil {
		return nil, fmt.Errorf("failed to generate ID suffix: %w", err)
	}

	id := fmt.Sprintf("INC-%s-%s", time.Now().Format("20060102"), suffix)

	// Path to the incident folder
	path := filepath.Join(incidentsDir, id)

	// Create .incidents folder if it doesn't exist
	if err := os.MkdirAll(incidentsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create incidents directory: %w", err)
	}

	// Create .gitignore inside .incidents/ if it doesn't exist
	gitignorePath := filepath.Join(incidentsDir, dotGitIgnore)
	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		if err := os.WriteFile(gitignorePath, []byte("*\n"), 0644); err != nil {
			return nil, fmt.Errorf("failed to create .gitignore: %w", err)
		}
	}

	// Create current incident directory
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, fmt.Errorf("failed to create incident directory: %w", err)
	}

	return &Incident{
		ID:    id,
		Title: title,
	}, nil
}
