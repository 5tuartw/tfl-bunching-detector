package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type splashTimeoutMsg struct{}

func (m *Model) updateSplashScreen(msg tea.Msg) (tea.Model, tea.Cmd) {
	var spinnerCmd tea.Cmd
	m.Spinner, spinnerCmd = m.Spinner.Update(msg)

	timeoutCmd := tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return splashTimeoutMsg{}
	})

	return m, tea.Batch(spinnerCmd, timeoutCmd)
}

func (m *Model) viewSplashScreen() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FAFAFA"))
	subtitleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	title := titleStyle.Render("TFL Bus Bunching Detector v1.0.0")
	subtitle := subtitleStyle.Render("A command-line tool to analyse and detect bus bunching.")
	loading := m.Spinner.View() + " Initialising..."

	content := lipgloss.JoinVertical(
		lipgloss.Top,
		title,
		"",
		subtitle,
		"",
		loading,
	)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), true).
		BorderForeground(lipgloss.Color("#4A90E2")).
		Padding(2, 4).
		Width(80)

	ui := boxStyle.Render(content)
	return lipgloss.NewStyle().
		MarginLeft(m.TermWidth - lipgloss.Width(ui)/2).
		Render(ui)
}
