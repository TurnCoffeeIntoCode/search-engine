package utils

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

func StartCronJobs() {
	c := cron.New()
	c.AddFunc("0 * * * *", runEngine) // Run every hour
	c.AddFunc("15 * * * *", runIndex) // Run every hour at 15 minutes past
	c.Start()
	cronCount := len(c.Entries())
	fmt.Printf("setup %d cron jobs \n", cronCount)
}

func runEngine() {
	fmt.Println("Running engine")
}
func runIndex() {
	fmt.Println("Running index")
}
