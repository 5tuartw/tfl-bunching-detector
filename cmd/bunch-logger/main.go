package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/5tuartw/tfl-bunching-detector/internal/analysis"
	"github.com/5tuartw/tfl-bunching-detector/internal/config"
	"github.com/5tuartw/tfl-bunching-detector/internal/lines"
	"github.com/5tuartw/tfl-bunching-detector/internal/logger"
	"github.com/5tuartw/tfl-bunching-detector/internal/models"
	"github.com/5tuartw/tfl-bunching-detector/internal/tflclient"
)

func main() {

	linesFlag := flag.String("lines", "", "comma-separated list of bus lines to track")
	interval := flag.Int("interval", 5, "interval between checks, default=5")
	bunchingThreshold := flag.Int("threshold", 90, "threshold for bunched buses in seconds")
	selectRoutesFlag := flag.Bool("select-routes", false, "set 'true' to select specific routes on a line")

	flag.Parse()

	if *linesFlag == "" {
		log.Fatal("Please specify at least one bus line to track and log using the -lines flag.")
	}

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("ERROR: unable to load config: %v", err)
	}
	httpClient := tflclient.NewClient("https://api.tfl.gov.uk", cfg.TflKey)

	var allSelectedLineRoutes []models.LineInfo

	linesSlice := strings.Split(*linesFlag, ",")
	for _, line := range linesSlice {
		lineInfo, err := lines.GetLineInfo(httpClient, strings.Trim(line, " "))
		if err != nil {
			log.Printf("Could not get line information for %s: %v", *linesFlag, err)
		}
		if *selectRoutesFlag {
			var modifiedLineInfo models.LineInfo
			modifiedLineInfo.LineId = lineInfo.LineId
			selectedRoutes := lines.ChooseRoute(lineInfo)
			for _, route := range selectedRoutes {
				modifiedLineInfo.Routes = append(modifiedLineInfo.Routes, lineInfo.Routes[route])
			}
			allSelectedLineRoutes = append(allSelectedLineRoutes, modifiedLineInfo)
		} else {
			allSelectedLineRoutes = append(allSelectedLineRoutes, lineInfo)
		}
	}

	log.Println("Running initial analysis...")

	err = runAnalysis(httpClient, *bunchingThreshold, allSelectedLineRoutes)
	if err != nil {
		log.Fatalf("FATAL: unable to run initial analysis: %v", err)
	}

	ticker := time.NewTicker(time.Duration(*interval) * time.Minute)
	defer ticker.Stop()
	log.Printf("Initial analysis complete. Starting periodic logger (interval: %d mins)\n", *interval)

	for range ticker.C {
		log.Println("Running analysis...")
		err = runAnalysis(httpClient, *bunchingThreshold, allSelectedLineRoutes)
		if err != nil {
			log.Printf("ERROR: unable to run periodic analysis: %v", err)
		}
	}
}

func runAnalysis(client *tflclient.Client, threshold int, linesToAnalyse []models.LineInfo) error {
	for _, line := range linesToAnalyse {
		for _, route := range line.Routes {
			routeBunchingEvents := analysis.AnalyseRoute(*client, line.LineId, threshold, route)
			err := logger.LogBunchingEvents(routeBunchingEvents, line.LineId)
			if err != nil {
				return fmt.Errorf("error logging bunching event for route %s, line %s: %v", route.Name, line.LineId, err)
			}
		}
	}
	return nil
}
