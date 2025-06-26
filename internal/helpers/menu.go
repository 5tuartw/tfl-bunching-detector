package helpers

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func GetMenuChoice(maxChoice int) ([]int, bool) {
	var selectedStops []int
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter one or more numbers from the list, separated by a space: ")

	if scanner.Scan() {
		userInput := scanner.Text()
		if userInput == "" {
			log.Println("Please enter at least one number.")
			return nil, false
		}
		selectedStopsStr := strings.Split(userInput, " ")
		for _, stopStr := range selectedStopsStr {
			stopInt, err := strconv.Atoi(stopStr)
			if err != nil {
				log.Printf("Error parsing choice: '%s'\n", stopStr)
				return nil, false
			}
			if stopInt > maxChoice || stopInt < 1 {
				log.Printf("Choice must be between 1-%d.\n", maxChoice)
				return nil, false
			}
			selectedStops = append(selectedStops, stopInt)
		}
	}

	return selectedStops, true
}
