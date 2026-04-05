package tests

import (
	"strings"
	"testing"

	"goincidentcli/internal/slack"

	"github.com/stretchr/testify/assert"
)

func TestSlugify(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple title",
			input:    "incident test",
			expected: "incident-test",
		},
		{
			name:     "Uppercase and special chars",
			input:    "INCIDENT #123: Database Failure!",
			expected: "incident-123-database-failure",
		},
		{
			name:     "Multiple spaces and hyphens",
			input:    "test---run  multiple",
			expected: "test-run-multiple",
		},
		{
			name:     "Trailing characters",
			input:    "incident-test-",
			expected: "incident-test",
		},
		{
			name:     "Long title",
			input:    strings.Repeat("a", 100),
			expected: strings.Repeat("a", 80),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := slack.Slugify(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
