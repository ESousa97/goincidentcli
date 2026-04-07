package tui

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"goincidentcli/internal/incident"
	"goincidentcli/internal/timeline"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	timelineReloadInterval = 10 * time.Second
	serviceCheckInterval   = 30 * time.Second
	maxTimelineEntries     = 5
)

// ServiceConfig holds the display name and health-check URL for a service.
type ServiceConfig struct {
	Name string
	URL  string
}

type serviceStatus struct {
	Name   string
	Status string // HEALTHY | DEGRADED | DOWN | UNKNOWN
}

// BubbleTea message types
type (
	tickMsg     time.Time
	timelineMsg []timeline.Entry
	servicesMsg []serviceStatus
)

// Model is the BubbleTea model for the incident status dashboard.
type Model struct {
	inc      *incident.Incident
	incDir   string
	entries  []timeline.Entry
	services []serviceStatus
	configs  []ServiceConfig
	width    int
	height   int
}

// NewModel creates a new dashboard model for the given incident.
func NewModel(inc *incident.Incident, incDir string, configs []ServiceConfig) Model {
	return Model{
		inc:     inc,
		incDir:  incDir,
		configs: configs,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(),
		cmdLoadTimeline(m.incDir),
		cmdCheckServices(m.configs),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "r":
			return m, tea.Batch(cmdLoadTimeline(m.incDir), cmdCheckServices(m.configs))
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tickMsg:
		return m, tickCmd()

	case timelineMsg:
		m.entries = []timeline.Entry(msg)
		return m, cmdScheduleTimelineReload(m.incDir)

	case servicesMsg:
		m.services = []serviceStatus(msg)
		return m, cmdScheduleServiceCheck(m.configs)
	}

	return m, nil
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading…"
	}

	// Inner content width (account for outer padding/borders)
	w := m.width - 2
	if w < 40 {
		w = 40
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		m.viewHeader(w),
		m.viewTimer(w),
		m.viewTimeline(w),
		m.viewServices(w),
		helpStyle.Render("  q quit • r refresh"),
	)
}

// ── Section renderers ────────────────────────────────────────────────────────

func (m Model) viewHeader(width int) string {
	sev := m.inc.Severity
	if sev == "" {
		sev = "N/A"
	}

	idStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFFFF"))
	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#CCCCCC"))

	badge := SeverityBadgeStyle(sev).Render(sev)
	top := lipgloss.JoinHorizontal(lipgloss.Top, idStyle.Render(m.inc.ID), "  ", badge)
	content := lipgloss.JoinVertical(lipgloss.Left,
		top,
		titleStyle.Render(m.inc.Title),
		dimStyle.Render("Declared: "+m.inc.CreatedAt.Format("2006-01-02 15:04:05 -07:00")),
	)

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(SeverityBorderColor(sev)).
		Width(width).
		Padding(0, 1).
		Render(content)
}

func (m Model) viewTimer(width int) string {
	elapsed := time.Since(m.inc.CreatedAt)
	h := int(elapsed.Hours())
	min := int(elapsed.Minutes()) % 60
	sec := int(elapsed.Seconds()) % 60

	content := lipgloss.JoinHorizontal(lipgloss.Center,
		sectionStyle.Render("⏱  Duration  "),
		timerValueStyle.Render(fmt.Sprintf("%02d:%02d:%02d", h, min, sec)),
	)

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#333333")).
		Width(width).
		Padding(0, 1).
		Render(content)
}

func (m Model) viewTimeline(width int) string {
	innerW := width - 6
	if innerW < 10 {
		innerW = 10
	}

	lines := []string{
		sectionStyle.Render("📋 Recent Events"),
		dimStyle.Render(strings.Repeat("─", innerW)),
	}

	entries := m.entries
	if len(entries) > maxTimelineEntries {
		entries = entries[len(entries)-maxTimelineEntries:]
	}

	if len(entries) == 0 {
		lines = append(lines, dimStyle.Render("No events recorded yet."))
	} else {
		// Most recent first
		for i := len(entries) - 1; i >= 0; i-- {
			e := entries[i]
			header := lipgloss.JoinHorizontal(lipgloss.Top,
				dimStyle.Render("["+e.Timestamp.Format("15:04:05")+"] "),
				authorStyle.Render(e.Author),
			)
			lines = append(lines, header)
			lines = append(lines, messageStyle.Render("  "+e.Message))
			for k, v := range e.Metrics {
				lines = append(lines, dimStyle.Render(fmt.Sprintf("  metric: %s = %s", k, v)))
			}
		}
	}

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#333366")).
		Width(width).
		Padding(0, 1).
		Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
}

func (m Model) viewServices(width int) string {
	innerW := width - 6
	if innerW < 10 {
		innerW = 10
	}

	lines := []string{
		sectionStyle.Render("🔍 Service Status"),
		dimStyle.Render(strings.Repeat("─", innerW)),
	}

	if len(m.services) == 0 {
		lines = append(lines,
			dimStyle.Render("No services configured."),
			dimStyle.Render("Add to ~/.incident.yaml:"),
			dimStyle.Render("  services:"),
			dimStyle.Render("    - name: My API"),
			dimStyle.Render("      url: http://api.example.com/health"),
		)
	} else {
		for _, svc := range m.services {
			st := ServiceStatusStyle(svc.Status)
			line := lipgloss.JoinHorizontal(lipgloss.Top,
				"  ",
				st.Render("●"),
				" ",
				fmt.Sprintf("%-24s", svc.Name),
				"  ",
				st.Render(svc.Status),
			)
			lines = append(lines, line)
		}
	}

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#333333")).
		Width(width).
		Padding(0, 1).
		Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
}

// ── BubbleTea commands ───────────────────────────────────────────────────────

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func cmdLoadTimeline(incDir string) tea.Cmd {
	return func() tea.Msg {
		entries, _ := timeline.Load(incDir)
		return timelineMsg(entries)
	}
}

func cmdScheduleTimelineReload(incDir string) tea.Cmd {
	return tea.Tick(timelineReloadInterval, func(_ time.Time) tea.Msg {
		return cmdLoadTimeline(incDir)()
	})
}

func cmdCheckServices(configs []ServiceConfig) tea.Cmd {
	return func() tea.Msg {
		client := &http.Client{Timeout: 5 * time.Second}
		statuses := make([]serviceStatus, len(configs))
		for i, cfg := range configs {
			s := serviceStatus{Name: cfg.Name}
			if cfg.URL == "" {
				s.Status = "UNKNOWN"
			} else {
				resp, err := client.Get(cfg.URL)
				if err != nil {
					s.Status = "DOWN"
				} else {
					resp.Body.Close()
					switch {
					case resp.StatusCode < 300:
						s.Status = "HEALTHY"
					case resp.StatusCode < 500:
						s.Status = "DEGRADED"
					default:
						s.Status = "DOWN"
					}
				}
			}
			statuses[i] = s
		}
		return servicesMsg(statuses)
	}
}

func cmdScheduleServiceCheck(configs []ServiceConfig) tea.Cmd {
	return tea.Tick(serviceCheckInterval, func(_ time.Time) tea.Msg {
		return cmdCheckServices(configs)()
	})
}
