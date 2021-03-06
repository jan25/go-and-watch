package main

import (
	"encoding/json"
	"log"
	"time"
)

type MonitorSetting struct {
	Interval time.Duration // in seconds
}

// Configurable
var monitorSettings MonitorSetting = MonitorSetting{Interval: 3}

// Receiver channel from parser
// To be called from main.go
func processAndMonitor() {

	monitorCh := make(chan AggregatedStats)
	go func() {
		monitorTicker := time.NewTicker(monitorSettings.Interval * time.Second)
		for {
			select {
			case <-monitorTicker.C:
				go computeAggregateStats(monitorSettings.Interval, monitorCh)
			default:

			}
		}
	}()

	for aggregatedStats := range monitorCh {
		monitorEndpointStats(aggregatedStats.EndPointStats)
		monitorStatusCodeStats(aggregatedStats.RequestStatusStats)
	}

}

func monitorEndpointStats(endPointStat []EndPointStat) {
	maxHits := int(0)
	totalHits := int(0)
	maxHitEndpoint := ""

	for _, element := range endPointStat {
		if element.Hits > maxHits {
			maxHits = element.Hits
			maxHitEndpoint = element.EndPoint
		}
		totalHits += element.Hits
	}
	log.Printf("Total hits %d", totalHits)
	log.Printf("Maximum hit endpoint %s hits: %d", maxHitEndpoint, maxHits)
}

func monitorStatusCodeStats(requestStatusStats []RequestStatusStat) {
	log.Printf("Request statuscode stats over last %d secs: %v", monitorSettings.Interval, requestStatusStats)

	statusCount := make(map[int]int)
	for _, element := range requestStatusStats {
		statusCount[element.Status] += 1
	}

	log.Println("Satuscode hit distribution ....")
	// Pretty print status code distribution in json format
	statusJsonMarshall, err := json.MarshalIndent(statusCount, "", " ")

	if err != nil {
		log.Fatalf("Failed to masrshall status code %v", err)
	}
	log.Print(string(statusJsonMarshall))

}
