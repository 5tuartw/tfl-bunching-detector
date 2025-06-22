package stops

import (
	"fmt"
	"log"
	"time"

	"github.com/5tuartw/tfl-bunching-detector/internal/data"
)

func checkDataFreshness() error {

	reader, err := data.NewBusStopReader()
	if err != nil {
		return fmt.Errorf("using embedded data, freshness check not available")
	}
	defer reader.Close()

	info := reader.Info()
	if !info.IsOS {
		return fmt.Errorf("using embedded data, freshness check not available")
	}

	if time.Since(info.ModTime).Hours() > 30*24 {
		log.Printf("WARN: the bus stop data is over 30 days old (last modified on %s).", info.ModTime.Format("2006-01-02"))
	} else {
		log.Printf("Bus stop data is up to date (from %s).", info.ModTime.Format("2006-01-02"))
	}

	return nil
}
