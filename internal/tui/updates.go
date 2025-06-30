package tui

import (
	"github.com/5tuartw/tfl-bunching-detector/internal/stops"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) updateMainMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}

		case "down", "j":
			if m.Cursor < len(m.MainMenuChoices)-1 {
				m.Cursor++
			}

		case "enter":
			switch m.Cursor {
			case 0: // spot-check a stop
				//transition to the next state
				m.State = StateStopSearchInput
				m.StopSearchInput.Focus()
				return m, nil
			case 1: //spot-check a line
				// todo
			case 2: // start a logging sesssion
				//todo
			}
		}
	}
	return m, nil
}

func (m *Model) updateStopSearchInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if key, ok := msg.(tea.KeyMsg); ok && key.Type == tea.KeyEnter {
		searchQuery := m.StopSearchInput.Value()
		m.SearchResults = stops.SearchStops(searchQuery, m.AllStops)
		m.State = StateStopSearchResults
		m.Cursor = 0
		m.StopSearchInput.Blur()
		return m, nil
	}

	m.StopSearchInput, cmd = m.StopSearchInput.Update(msg)
	return m, cmd
}

func (m *Model) updateStopSearchResults(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			if m.Cursor < len(m.SearchResults)-1 {
				m.Cursor++
			}
		case " ": // Toggle selection
			if _, ok := m.SelectedStops[m.Cursor]; ok {
				delete(m.SelectedStops, m.Cursor)
			} else {
				m.SelectedStops[m.Cursor] = struct{}{}
			}
		case "enter": // Confirm selection
			m.State = StateQuitting
			return m, tea.Quit
		}
	}
	return m, nil
}
