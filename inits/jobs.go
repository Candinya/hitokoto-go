package inits

import "hitokoto-go/jobs"

func Jobs() error {
	jobs.StartJobs()
	return nil
}
