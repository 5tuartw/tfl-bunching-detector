package stops

import (
	"log"
	"time"

	"github.com/5tuartw/tfl-bunching-detector/internal/data"
)

func checkDataFreshness() {
	
	file, err := data.DataFS.Open("bus-stops.csv")
	if err != nil {
		log.Fatalf("FATAL: could not open embedded bus stop data: %v", err)
	}
	defer file.Close()

	stats, err := file.Stat()
	if err != nil {
		log.Fatalf("FATAL: could not get stats for embedded bus stop data: %v", err)
	}

	modTime := stats.ModTime()

	if time.Since(modTime).Hours() > 30*24 {
		log.Printf("WARN: the bus stop data is over 30 days old (last modified on %s).", modTime.Format("2006-01-02"))
	} else {
		log.Printf("Bus stop data is up to data (from %s).", modTime.Format("2006-01-02"))
	}
}
