package main

import (
	"fmt"
	"log"
	"os"

	"github.com/5tuartw/tfl-bunching-detector/internal/analysis"
	"github.com/5tuartw/tfl-bunching-detector/internal/config"
	"github.com/5tuartw/tfl-bunching-detector/internal/display"
	"github.com/5tuartw/tfl-bunching-detector/internal/helpers"
	"github.com/5tuartw/tfl-bunching-detector/internal/models"
	"github.com/5tuartw/tfl-bunching-detector/internal/stops"
	"github.com/5tuartw/tfl-bunching-detector/internal/tflclient"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type appState int

const (
	stateMainMenu appState = iota
	stateStopSearchInput
	stateStopSearchResults
	stateThresholdInput
	stateRunningAnalysis
	stateQuitting
	// ...
)

type model struct {
	state     appState // current state
	allStops  []models.BusStop
	tflClient *tflclient.Client
	threshold int
	err       error

	//main menu data
	mainMenuChoices []string
	cursor          int

	//Stop search data
	stopSearchInput textinput.Model
	searchResults   []models.BusStop
	selectedStops   map[int]struct{}
	thresholdInput  textinput.Model
}

func initialModel(allStops []models.BusStop, client *tflclient.Client) *model {
	ti := textinput.New()
	ti.Placeholder = "e.g. Victoria Station"
	ti.CharLimit = 156
	ti.Width = 50

	return &model{
		state:           stateMainMenu,
		allStops:        allStops,
		tflClient:       client,
		threshold:       90,
		mainMenuChoices: []string{"Spot-check a Stop", "Spot-check a Line (coming soon)", "Start a Logging Session (coming soon)"},
		stopSearchInput: ti,
		selectedStops:   make(map[int]struct{}),
	}
}

func (m *model) Init() tea.Cmd {
	return textinput.Blink
}

// Main View Router
func (m *model) View() string {
	if m.state == stateQuitting {
		return ""
	}

	switch m.state {
	case stateStopSearchInput:
		return m.viewStopSearchInput()
	case stateStopSearchResults:
		return m.viewStopSearchResults()
	//case stateThresholdInput:
	//	return m.viewThresholdInput()
	default:
		return m.viewMainMenu()
	}
}

// Main Menu Router
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC || msg.String() == "q" {
			m.state = stateQuitting
			return m, tea.Quit
		}
	}

	switch m.state {
	case stateStopSearchInput:
		return m.updateStopSearchInput(msg)
	case stateStopSearchResults:
		return m.updateStopSearchResults(msg)
	//case stateThresholdInput:
	//	return m.updateThresholdInput(msg)
	// ... other states
	default: // stateMainMenu
		return m.updateMainMenu(msg)
	}
}

func (m *model) viewMainMenu() string {
	s := "What would you like to do?\n\n"
	for i, choice := range m.mainMenuChoices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}
	s += "\nUse up and down arrows and enter to select. (Press q or ctrl+c to quit)"
	return s
}

func (m *model) updateMainMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.mainMenuChoices)-1 {
				m.cursor++
			}

		case "enter":
			switch m.cursor {
			case 0: // spot-check a stop
				//transition to the next state
				m.state = stateStopSearchInput
				m.stopSearchInput.Focus()
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

func (m *model) viewStopSearchInput() string {
	return fmt.Sprintf(
		"Enter the name of a bus stop to search for:\n\n%s\n\n(Press Enter when done)",
		m.stopSearchInput.View(),
	)
}

func (m *model) updateStopSearchInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if key, ok := msg.(tea.KeyMsg); ok && key.Type == tea.KeyEnter {
		searchQuery := m.stopSearchInput.Value()
		m.searchResults = stops.SearchStops(searchQuery, m.allStops)
		m.state = stateStopSearchResults
		m.cursor = 0
		m.stopSearchInput.Blur()
		return m, nil
	}

	m.stopSearchInput, cmd = m.stopSearchInput.Update(msg)
	return m, cmd
}

func (m *model) viewStopSearchResults() string {
	s := "Found these stops. Use space to select, enter to confirm.\n\n"

	for i, stop := range m.searchResults {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		checked := " "
		if _, ok := m.selectedStops[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s (%s)\n ", cursor, checked, stop.StopName, helpers.HeadingToDirection(stop.Heading))
	}
	s += "\n(q to quit)\n"
	return s
}

func (m *model) updateStopSearchResults(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.searchResults)-1 {
				m.cursor++
			}
		case " ": // Toggle selection
			if _, ok := m.selectedStops[m.cursor]; ok {
				delete(m.selectedStops, m.cursor)
			} else {
				m.selectedStops[m.cursor] = struct{}{}
			}
		case "enter": // Confirm selection
			m.state = stateQuitting
			return m, tea.Quit
		}
	}
	return m, nil
}

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("ERROR: unable to load config: %v", err)
	}
	tflClient := tflclient.NewClient("https://api.tfl.gov.uk", cfg.TflKey)

	allStops, err := stops.LoadBusStops()
	if err != nil {
		log.Fatalf("Failed to load bus stops: %v", err)
	}

	p := tea.NewProgram(initialModel(allStops, tflClient))
	finalModel, err := p.Run()
	if err != nil {
		log.Fatalf("Error running program: %v", err)
	}

	m, ok := finalModel.(*model)
	if !ok {
		os.Exit(1)
	}

	if len(m.selectedStops) > 0 {
		fmt.Println("\n--- Analysing Selected Stop ---")
		var chosenStops []models.BusStop
		for index := range m.selectedStops {
			chosenStops = append(chosenStops, m.searchResults[index])
		}

		for _, stop := range chosenStops {
			fmt.Printf("\nFetching data for stop: %s\n", stop.StopName)
			arrivalInfo, err := stops.GetStopArrivalInfo(tflClient, stop.NaptanId)
			if err != nil {
				log.Fatalf("EROR: could not get arrival information for stop %s (%s): %v", stop.StopName, stop.NaptanId, err)
			}
			bunchingEvents := analysis.AnalyseArrivals(arrivalInfo, "", m.threshold)
			stopName := helpers.GetStopName(stop)
			display.PrintBunchingData(stopName, m.threshold, bunchingEvents)
		}
	}

}

/*
// creates initial model
func initialModel() model {
	return model{
		choices:  []string{"Spot-check a Stop", "Spot-check a Line", "Start a Logging Session"},
		selected: make(map[int]struct{}),
	}
}

// Init is the first command run when program starts
func (m model) Init() tea.Cmd {
	return nil
}

// Update to handle incoming events, e.g. key presses
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// logic for up/down arrow keys goes here

	// logic for quitting with 'q' or 'ctrl-c'
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case " ":
			_, ok := m.selected[m.cursor]
			if ok {
				// if selected, deselect
				delete(m.selected, m.cursor)
			} else {
				// otherwise select and add to map
				m.selected[m.cursor] = struct{}{}
			}

		case "enter":
			return m, tea.Quit
		}
	}
	return m, nil
}

// View renders the UI as a string
func (m model) View() string {
	s := "What would you like to do?\n\n"

	for i, choice := range m.choices {
		cursor := " " // a space for when cursor is not on this item
		if m.cursor == i {
			cursor = ">" // cursor shows which item is selected
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	s += "\nUse space to select, enter to confirm. Press q to quit.\n"
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error starting terminal user interface: %v", err)
		os.Exit(1)
	}
}
*/
