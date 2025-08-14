package tui

import (
	"github.com/5tuartw/tfl-bunching-detector/internal/models"
	"github.com/5tuartw/tfl-bunching-detector/internal/tflclient"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

type AppState int

const (
	StateSplashScreen AppState = iota
	StateMainMenu
	StateStopSearchInput
	StateStopSearchResults
	StateThresholdInput
	StateRunningAnalysis
	StateQuitting
	// ...
)

type Model struct {
	State      AppState // current state
	AllStops   []models.BusStop
	TflClient  *tflclient.Client
	Threshold  int
	Err        error
	TermWidth  int
	TermHeight int

	//main menu data
	MainMenuChoices []string
	Cursor          int
	CursorStyle     lipgloss.Style

	//Stop search data
	StopSearchInput textinput.Model
	SearchResults   []models.BusStop
	SelectedStops   map[int]struct{}
	ThresholdInput  textinput.Model

	Spinner spinner.Model

	CurrentPage  int
	ItemsPerPage int
	TotalPages   int
}

func NewModel(allStops []models.BusStop, client *tflclient.Client) *Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	ti := textinput.New()
	ti.Placeholder = "e.g. Victoria Station"
	ti.CharLimit = 156
	ti.Width = 50

	return &Model{
		State:           StateSplashScreen,
		Spinner:         s,
		AllStops:        allStops,
		TflClient:       client,
		Threshold:       90,
		MainMenuChoices: []string{"Spot-check a Stop", "Spot-check a Line (coming soon)", "Start a Logging Session (coming soon)"},
		StopSearchInput: ti,
		SelectedStops:   make(map[int]struct{}),
		ItemsPerPage:    10, // Show 10 items per page
		CursorStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("#4A90E2")).Bold(true),
	}
}
