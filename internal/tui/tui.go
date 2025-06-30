package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

// Main Menu Router
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.TermWidth = msg.Width
		return m, nil

	case splashTimeoutMsg:
		m.State = StateMainMenu
		m.Cursor = 0
		return m, nil

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC || msg.String() == "q" {
			m.State = StateQuitting
			return m, tea.Quit
		}
	}

	switch m.State {
	case StateSplashScreen:
		return m.updateSplashScreen(msg)
	case StateStopSearchInput:
		return m.updateStopSearchInput(msg)
	case StateStopSearchResults:
		return m.updateStopSearchResults(msg)
	//case stateThresholdInput:
	//	return m.updateThresholdInput(msg)
	// ... other states
	default: // stateMainMenu
		return m.updateMainMenu(msg)
	}
}

// Main View Router
func (m *Model) View() string {
	if m.Err != nil {
		return fmt.Sprintf("An error occured: %v", m.Err)
	}
	if m.State == StateQuitting {
		return ""
	}

	switch m.State {
	case StateSplashScreen:
		return m.viewSplashScreen()
	case StateStopSearchInput:
		return m.viewStopSearchInput()
	case StateStopSearchResults:
		return m.viewStopSearchResults()
	//case stateThresholdInput:
	//	return m.viewThresholdInput()
	default:
		return m.viewMainMenu()
	}
}
