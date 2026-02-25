package channel

import (
	"github.com/charmbracelet/lipgloss"
)

// Styles holds all the UI styles
var Styles = struct {
	Title         lipgloss.Style
	Subtitle      lipgloss.Style
	ChannelName   lipgloss.Style
	Description   lipgloss.Style
	Selected      lipgloss.Style
	Normal        lipgloss.Style
	Configured    lipgloss.Style
	NotConfigured lipgloss.Style
	Enabled       lipgloss.Style
	Disabled      lipgloss.Style
	Help          lipgloss.Style
	Error         lipgloss.Style
	Success       lipgloss.Style
	Warning       lipgloss.Style
	Header        lipgloss.Style
	Footer        lipgloss.Style
	FormLabel     lipgloss.Style
	FormInput     lipgloss.Style
	FormFocused   lipgloss.Style
	FormHelp      lipgloss.Style
	Box           lipgloss.Style
}{
	Title: lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		MarginLeft(2).
		MarginBottom(1),

	Subtitle: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		MarginLeft(2),

	ChannelName: lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")),

	Description: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")),

	Selected: lipgloss.NewStyle().
		Background(lipgloss.Color("#7D56F4")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true).
		Padding(0, 1),

	Normal: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#CCCCCC")).
		Padding(0, 1),

	Configured: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#04B575")),

	NotConfigured: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF6B6B")),

	Enabled: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#04B575")).
		Bold(true),

	Disabled: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")),

	Help: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Italic(true),

	Error: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF6B6B")).
		Bold(true),

	Success: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#04B575")).
		Bold(true),

	Warning: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFB347")).
		Bold(true),

	Header: lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(1, 2).
		Width(80),

	Footer: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Padding(0, 2),

	FormLabel: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true).
		Width(25).
		Align(lipgloss.Right),

	FormInput: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#CCCCCC")),

	FormFocused: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true),

	FormHelp: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Italic(true).
		MarginTop(1),

	Box: lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(1, 2).
		Margin(1, 2),
}

// GetStatusIcon returns the appropriate icon for a channel status
func GetStatusIcon(configured, enabled bool) string {
	if !configured {
		return Styles.NotConfigured.Render("✗")
	}
	if enabled {
		return Styles.Enabled.Render("▶")
	}
	return Styles.Configured.Render("✓")
}

// GetStatusText returns the appropriate text for a channel status
func GetStatusText(configured, enabled bool) string {
	if !configured {
		return Styles.NotConfigured.Render("未配置")
	}
	if enabled {
		return Styles.Enabled.Render("已启用")
	}
	return Styles.Configured.Render("已配置")
}
