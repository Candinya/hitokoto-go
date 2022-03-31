package jobs

import (
	"log"
	"time"
)

func StartJobs() {
	// Init run
	RefreshCache()

	go func() {
		// Interval run
		t := time.NewTicker(1 * time.Minute) // Run per minute
		for {
			select {
			case <-t.C:
				jobs()
			}
		}
	}()

}

func jobs() {
	log.Println("Start cron job")

	// Collect request data
	CollectStats()

	// Refresh cache
	RefreshCache()

	log.Println("Cron job done")
}
