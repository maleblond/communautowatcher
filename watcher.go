package communautowatcher

import (
	"context"
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
	LanguageID    string
	BranchID      string
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
	GetFlexCarQuery() CarQuery
	OnCarAvailable(query CarQuery, cars []Car)
	OnFlexCarAvailable(cars []FlexCarAvailabilityResp)
}

type WatcherOptions struct {
	Interval        time.Duration
	Watcher         Watcher
	IsFetchStations bool
	IsFetchFlexCars bool
}

func StartWatcher(ctx context.Context, options WatcherOptions) error {
	ticker := time.NewTicker(options.Interval)
	defer ticker.Stop()

	err := checkForAvailabilities(options)
	if err != nil {
		return fmt.Errorf("[checkForAvailabilities] Error: %s", err)
	}

	err = checkForFlexCars(ctx, options)
	if err != nil {
		return fmt.Errorf("[checkForFlexCars] Error: %s", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err = checkForAvailabilities(options); err != nil {
				return fmt.Errorf("[checkForAvailabilities] Error: %s", err)
			}
			if err = checkForFlexCars(ctx, options); err != nil {
				return fmt.Errorf("[checkForFlexCars] Error: %s", err)
			}
		}
	}
}

func checkForAvailabilities(options WatcherOptions) error {
	if !options.IsFetchStations {
		return nil
	}

	watcher := options.Watcher

	queries := options.Watcher.GetQueries()

	for _, query := range queries {
		cars, err := GetAvailableCars(query)

		if err != nil {
			return fmt.Errorf("Could not retrieve available cars: %v\n", err)
		}

		if len(cars) > 0 {
			watcher.OnCarAvailable(query, cars)
		}
	}

	return nil
}

func checkForFlexCars(ctx context.Context, options WatcherOptions) error {
	if !options.IsFetchFlexCars {
		return nil
	}

	query := options.Watcher.GetFlexCarQuery()
	response, err := GetAvailableFlexCars(ctx, query)
	if err != nil {
		return fmt.Errorf("Could not retrieve flex cars: %v\n", err)
	}
	if len(response.Data.Vehicles) > 0 {
		options.Watcher.OnFlexCarAvailable(response.Data.Vehicles)
	}

	return nil
}
