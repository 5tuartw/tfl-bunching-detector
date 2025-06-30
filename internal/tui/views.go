package tui

import (
	"fmt"

	"github.com/5tuartw/tfl-bunching-detector/internal/helpers"
)

func (m *Model) viewMainMenu() string {
	s := "What would you like to do?\n\n"
	for i, choice := range m.MainMenuChoices {
		Cursor := " "
		if m.Cursor == i {
			Cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", Cursor, choice)
	}
	s += "\nUse up and down arrows and enter to select. (Press q or ctrl+c to quit)"
	return s
}

func (m *Model) viewStopSearchInput() string {
	return fmt.Sprintf(
		"Enter the name of a bus stop to search for:\n\n%s\n\n(Press Enter when done)",
		m.StopSearchInput.View(),
	)
}

func (m *Model) viewStopSearchResults() string {
	s := "Found these stops. Use space to select, enter to confirm.\n\n"

	for i, stop := range m.SearchResults {
		Cursor := " "
		if m.Cursor == i {
			Cursor = ">"
		}
		checked := " "
		if _, ok := m.SelectedStops[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s (%s)\n ", Cursor, checked, stop.StopName, helpers.HeadingToDirection(stop.Heading))
	}
	s += "\n(q to quit)\n"
	return s
}
