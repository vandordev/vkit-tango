package river

import (
	"database/sql"

	riverqueue "github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverdatabasesql"
)

func NewClient(database *sql.DB, workers *riverqueue.Workers, maxWorkers int, periodicJobs []*riverqueue.PeriodicJob) (*riverqueue.Client[*sql.Tx], error) {
	return riverqueue.NewClient(riverdatabasesql.New(database), &riverqueue.Config{
		Workers:      workers,
		PeriodicJobs: periodicJobs,
		Queues: map[string]riverqueue.QueueConfig{
			"default": {MaxWorkers: maxWorkers},
		},
	})
}
