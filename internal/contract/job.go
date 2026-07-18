package contract

import riverqueue "github.com/riverqueue/river"

type WorkerRegistrar interface {
	RegisterWorkers(*riverqueue.Workers)
}
