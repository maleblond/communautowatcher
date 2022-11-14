package communautowatcher

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

type Car struct {
	Distance     float64
	Latitude     float64
	Longitude    float64
	LocationName string
	IsPromo      bool
	CarBrand     string
	CarModel     string
	CarNo        int
	CarPlate     string
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
	OnFlexCarAvailable(cars []Car)
}

type WatcherOptions struct {
	Interval              time.Duration
	Watcher               Watcher
	IsEnableFetchStations bool
	IsEnableFetchFlexCars bool
}

func StartWatcher(ctx context.Context, options WatcherOptions) {
	ticker := time.NewTicker(options.Interval)
	defer ticker.Stop()

	err := checkForAvailabilities(ctx, options)
	if err != nil {
		log.Error("[checkForAvailabilities] Error: %s", err)
	}

	err = checkForFlexCars(ctx, options)
	if err != nil {
		log.Error("[checkForFlexCars] Error: %s", err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err = checkForAvailabilities(ctx, options); err != nil {
				log.Error("[checkForAvailabilities] Error: %s", err)
			}
			if err = checkForFlexCars(ctx, options); err != nil {
				log.Error("[checkForFlexCars] Error: %s", err)
			}
		}
	}
}

func checkForAvailabilities(ctx context.Context, options WatcherOptions) error {
	if !options.IsEnableFetchStations {
		return nil
	}

	watcher := options.Watcher

	queries := options.Watcher.GetQueries()

	for _, query := range queries {
		cars, err := GetAvailableCars(ctx, query)

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
	if !options.IsEnableFetchFlexCars {
		return nil
	}

	query := options.Watcher.GetFlexCarQuery()
	response, err := GetAvailableFlexCars(ctx, query)
	if err != nil {
		return fmt.Errorf("Could not retrieve flex cars: %v\n", err)
	}
	if len(response) > 0 {
		options.Watcher.OnFlexCarAvailable(response)
	}

	return nil
}
