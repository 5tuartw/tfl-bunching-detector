package tui

import (
	"fmt"

	"github.com/5tuartw/tfl-bunching-detector/internal/helpers"
	"github.com/5tuartw/tfl-bunching-detector/internal/models"
)

func (m *Model) viewMainMenu() string {
	s := "What would you like to do?\n\n"
	for i, choice := range m.MainMenuChoices {
		cursor := " "
		if m.Cursor == i {
			cursor = m.CursorStyle.Render(">")
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
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
	if len(m.SearchResults) == 0 {
		return "No stops found. Press Enter to go back.\n"
	}

	s := fmt.Sprintf("Found these stops. Use space to select, enter to confirm. (%d selected)\n\n", len(m.SelectedStops))

	// Get current page items
	currentPageItems := m.getCurrentPageItems()

	// Display items for current page
	for i, stop := range currentPageItems {
		cursor := " "
		if m.Cursor == i {
			cursor = m.CursorStyle.Render(">")
		}

		// Calculate the actual index in the full SearchResults slice
		actualIndex := m.CurrentPage*m.ItemsPerPage + i
		checked := " "
		if _, ok := m.SelectedStops[actualIndex]; ok {
			checked = m.CursorStyle.Render("x")
		}

		s += fmt.Sprintf("%s [%s] %s (%s)\n", cursor, checked, stop.StopName, helpers.HeadingToDirection(stop.Heading))
	}

	// Show pagination info
	if m.TotalPages > 1 {
		s += fmt.Sprintf("\nPage %d of %d", m.CurrentPage+1, m.TotalPages)
		s += "\nUse left/right (h/l) arrows to navigate pages"
	}

	s += "\nSpace to select/deselect, Enter to confirm, q to quit"
	s += fmt.Sprintf("\n[DEBUG] Terminal height: %d", m.TermHeight)
	return s
}

func (m *Model) getCurrentPageItems() []models.BusStop {
	start := m.CurrentPage * m.ItemsPerPage
	end := start + m.ItemsPerPage
	if end > len(m.SearchResults) {
		end = len(m.SearchResults)
	}
	return m.SearchResults[start:end]
}
