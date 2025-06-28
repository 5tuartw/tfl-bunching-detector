package logger

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/5tuartw/tfl-bunching-detector/internal/models"
)

func LogBunchingEvents(bunches []models.BunchingEvent, lineId string) error {

	csvFile, err := os.OpenFile("internal/data/logged_bunching_events.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not open/create log file: %v", err)
	}
	defer csvFile.Close()

	fileInfo, err := csvFile.Stat()
	if err != nil {
		return fmt.Errorf("could not get file info of csv file: %v", err)
	}

	csvWriter := csv.NewWriter(csvFile)

	if fileInfo.Size() == 0 {
		log.Println("Log file is new, writing headers...")
		headers := []string{"LineID", "StopID", "StopName", "EventTime", "Headway", "VehicleIDs"}
		err := csvWriter.Write(headers)
		if err != nil {
			return fmt.Errorf("could not write headers to csv file: %v", err)
		}
	}

	log.Printf("Writing %d bunching events to the log for line %s...", len(bunches), lineId)

	for _, event := range bunches {
		record := []string{
			event.LineId,
			event.NaptanId,
			event.StationName,
			event.EventTime.Format("2006-01-02 15:04:05"),
			fmt.Sprintf("%v", event.Headway),
			strings.Join(event.VehicleIds, "|"),
		}
		csvWriter.Write(record)
	}
	csvWriter.Flush()
	if csvWriter.Error() != nil {
		return fmt.Errorf("error writing entries to csv file: %v", csvWriter.Error())
	}
	return nil
}
