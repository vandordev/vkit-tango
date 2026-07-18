package river

import riverqueue "github.com/riverqueue/river"

// RegisterPeriodicJobs returns the deterministic schedule set installed on
// every worker replica. Deadline-sensitive work must enqueue a reconciliation
// job that scans due rows idempotently; a periodic tick is never the source of
// truth. The domain-neutral baseline intentionally has no schedules.
func RegisterPeriodicJobs() []*riverqueue.PeriodicJob {
	return []*riverqueue.PeriodicJob{}
}
