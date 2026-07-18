package river

import (
	riverqueue "github.com/riverqueue/river"
	"github.com/vandordev/vkit-tango/internal/contract"
)

type BaselineRegistrar struct{}

var _ contract.SchedulerRegistrar = (*BaselineRegistrar)(nil)

func NewBaselineRegistrar() *BaselineRegistrar { return &BaselineRegistrar{} }

func (*BaselineRegistrar) RegisterPeriodicJobs() []*riverqueue.PeriodicJob {
	return []*riverqueue.PeriodicJob{}
}
