// Package tui contains the BubbleTea terminal UI components for goincidentcli.
package tui

import "github.com/charmbracelet/lipgloss"

// Severity palette
var (
	sev1Fg     = lipgloss.Color("#FF4444")
	sev1Bg     = lipgloss.Color("#2D0000")
	sev1Border = lipgloss.Color("#FF4444")

	sev2Fg     = lipgloss.Color("#FF8C00")
	sev2Bg     = lipgloss.Color("#2D1800")
	sev2Border = lipgloss.Color("#FF8C00")

	sev3Fg     = lipgloss.Color("#FFD700")
	sev3Bg     = lipgloss.Color("#2D2600")
	sev3Border = lipgloss.Color("#FFD700")
)

// General styles (package-level, reused across renders)
var (
	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#555555"))

	sectionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#AAAAAA"))

	timerValueStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00E5FF"))

	authorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7B9FFF"))

	messageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#DDDDDD"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#444444"))
)

// SeverityBadgeStyle returns the badge style for a severity level.
func SeverityBadgeStyle(sev string) lipgloss.Style {
	switch sev {
	case "SEV1":
		return lipgloss.NewStyle().Bold(true).Foreground(sev1Fg).Background(sev1Bg).Padding(0, 1)
	case "SEV2":
		return lipgloss.NewStyle().Bold(true).Foreground(sev2Fg).Background(sev2Bg).Padding(0, 1)
	case "SEV3":
		return lipgloss.NewStyle().Bold(true).Foreground(sev3Fg).Background(sev3Bg).Padding(0, 1)
	default:
		return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#888888")).Padding(0, 1)
	}
}

// SeverityBorderColor returns the box border color for a severity level.
func SeverityBorderColor(sev string) lipgloss.Color {
	switch sev {
	case "SEV1":
		return sev1Border
	case "SEV2":
		return sev2Border
	case "SEV3":
		return sev3Border
	default:
		return lipgloss.Color("#444444")
	}
}

// ServiceStatusStyle returns the style for a service status indicator.
func ServiceStatusStyle(status string) lipgloss.Style {
	switch status {
	case "HEALTHY":
		return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00C896"))
	case "DEGRADED":
		return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF8C00"))
	case "DOWN":
		return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF4444"))
	default:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
	}
}
