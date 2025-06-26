package stops

import (
	"log"
	"time"

	"github.com/5tuartw/tfl-bunching-detector/internal/data"
)

func checkDataFreshness(info data.FileInfo) {
	if !info.IsOS {
		log.Println("Using embedded bus stop data; freshness check not available.")
		return
	}

	if time.Since(info.ModTime).Hours() > 30*24 {
		log.Printf("WARN: the bus stop data is over 30 days old (last modified on %s).", info.ModTime.Format("2006-01-02"))
	} else {
		log.Printf("Bus stop data is up to date (from %s).", info.ModTime.Format("2006-01-02"))
	}
}
