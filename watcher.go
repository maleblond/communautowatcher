package communautowatcher

import (
	"fmt"
	"time"
)

type Car struct {
	Distance     float64
	Latitude     float64
	Longitude    float64
	LocationName string
}

type CarQuery struct {
	FromLatitude  string
	FromLongitude string
	CityID        string
	MaxDistance   float64
	StartDate     time.Time
	EndDate       time.Time
}

type CityID string

const (
	Montreal   CityID = "59"
	Quebec            = "90"
	Sherbrooke        = "89"
	// Inspected requests from https://www.reservauto.net/Scripts/Client/Mobile/Default.asp?BranchID=1#
)

type Watcher interface {
	GetQueries() []CarQuery
	OnCarAvailable(query CarQuery, cars []Car)
}

type WatcherOptions struct {
	Interval time.Duration
	Watcher  Watcher
}

func StartWatcher(options WatcherOptions) {
	checkForAvailabilities(options)

	for range time.Tick(options.Interval) {
		checkForAvailabilities(options)
	}
}

func checkForAvailabilities(options WatcherOptions) {
	watcher := options.Watcher

	queries := options.Watcher.GetQueries()

	for _, query := range queries {
		cars, err := GetAvailableCars(query)

		if err != nil {
			fmt.Printf("Could not retrieve available cars: %v\n", err)
		}

		if len(cars) > 0 {
			watcher.OnCarAvailable(query, cars)
		}
	}
}
