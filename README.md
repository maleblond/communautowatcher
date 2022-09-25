# Communauto-watcher

Communauto is a car sharing service in Canada. It can sometimes be tricky to get a hold on a Communauto.

This package makes it easier to build an app that notifies you when a car gets available nearby for a given time frame:

```
package main

import (
	"sync"
	"time"

	"github.com/maleblond/communautowatcher"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done() // Should never be called, `StartWatcher` never finishes
		communautowatcher.StartWatcher(communautowatcher.WatcherOptions{
			Interval: time.Minute * 5,
			Watcher:  &Watcher{},
		})
	}()

	wg.Wait()
}

type Watcher struct{}

func (w *Watcher) GetQueries() []communautowatcher.CarQuery {
    // Could fetch your "queries" from a database or a config file
	startDate, _ := time.Parse("2006-01-02T15:04", "2022-10-01T11:00")
	endDate, _ := time.Parse("2006-01-02T15:04", "2022-10-01T11:30")

	return []communautowatcher.CarQuery{
		{
			StartDate:     startDate,
			EndDate:       endDate,
			FromLatitude:  "46.8046123",
			FromLongitude: "-71.2342123",
			CityID:        communautowatcher.Quebec,
		},
	}
}

func (w *Watcher) OnCarAvailable(query communautowatcher.CarQuery, cars []communautowatcher.Car) {
	// Notify by email, update the state of the query to stop notifying etc...
}
```