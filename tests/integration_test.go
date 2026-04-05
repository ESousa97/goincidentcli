package tests

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"testing"
	"time"

	"goincidentcli/cmd"
	"goincidentcli/internal/incident"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIncidentDeclare_Functional(t *testing.T) {
	// 1. Setup temporary directories for Home and Project Root
	tempHome := t.TempDir()
	tempProject := t.TempDir()

	// Change working directory to tempProject for the duration of the test
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	err = os.Chdir(tempProject)
	require.NoError(t, err)
	defer func() { _ = os.Chdir(originalWd) }()

	// 2. Mocking Home directory via environment variables
	homeEnv := "HOME"
	if runtime.GOOS == "windows" {
		homeEnv = "USERPROFILE"
	}
	originalHome := os.Getenv(homeEnv)
	os.Setenv(homeEnv, tempHome)
	defer os.Setenv(homeEnv, originalHome)

	// 3. Execution: Run 'incident declare --title "Test Incident"'
	// We call RootCmd.Execute() directly to avoid os.Exit(1) context
	cmd.RootCmd.SetArgs([]string{"declare", "--title", "Test Incident"})
	err = cmd.RootCmd.Execute()
	require.NoError(t, err)

	// --- VALIDATIONS ---

	// 4. Validate ~/.incident.yaml creation (Config Template)
	cfgPath := filepath.Join(tempHome, ".incident.yaml")
	assert.FileExists(t, cfgPath, "Config file should be created in home directory")
	
	// Read and check content of config
	cfgContent, err := os.ReadFile(cfgPath)
	require.NoError(t, err)
	assert.Contains(t, string(cfgContent), "api_token", "Config should contain api_token key")
	assert.Contains(t, string(cfgContent), "base_url", "Config should contain base_url key")

	// 5. Validate .incidents/.gitignore creation
	gitignorePath := filepath.Join(tempProject, ".incidents", ".gitignore")
	assert.FileExists(t, gitignorePath, ".incidents/.gitignore should be created")
	
	gitIgnoreContent, err := os.ReadFile(gitignorePath)
	require.NoError(t, err)
	assert.Equal(t, "*\n", string(gitIgnoreContent), ".gitignore content should be '*'")

	// 6. Validate Incident Folder and metadata.json
	entries, err := os.ReadDir(filepath.Join(tempProject, ".incidents"))
	require.NoError(t, err)

	var incidentID string
	found := false
	// Regular Expression for INC-YYYYMMDD-NanoID (4 chars)
	// Example: INC-20240405-abcd
	idRegex := regexp.MustCompile(`^INC-\d{8}-[a-zA-Z0-9]{4}$`)

	for _, entry := range entries {
		if entry.IsDir() && idRegex.MatchString(entry.Name()) {
			incidentID = entry.Name()
			found = true
			break
		}
	}

	assert.True(t, found, "An incident directory with valid ID format should be created")

	// 7. Validate metadata.json content
	metadataPath := filepath.Join(tempProject, ".incidents", incidentID, "metadata.json")
	assert.FileExists(t, metadataPath, "metadata.json should exist inside incident folder")

	metadataBytes, err := os.ReadFile(metadataPath)
	require.NoError(t, err)

	var inc incident.Incident
	err = json.Unmarshal(metadataBytes, &inc)
	require.NoError(t, err)

	assert.Equal(t, incidentID, inc.ID)
	assert.Equal(t, "Test Incident", inc.Title)
	assert.WithinDuration(t, time.Now(), inc.CreatedAt, 2*time.Second)
}
