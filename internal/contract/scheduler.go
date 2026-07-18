package contract

import riverqueue "github.com/riverqueue/river"

type SchedulerRegistrar interface {
	RegisterPeriodicJobs() []*riverqueue.PeriodicJob
}
