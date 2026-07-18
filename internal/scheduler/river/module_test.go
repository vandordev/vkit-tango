package river_test

import (
	"testing"

	"github.com/vandordev/vkit-tango/internal/contract"
	scheduler "github.com/vandordev/vkit-tango/internal/scheduler/river"
)

func TestBaselineRegistrarUsesEmptySchedule(t *testing.T) {
	registrar := scheduler.NewBaselineRegistrar()
	var _ contract.SchedulerRegistrar = registrar
	if jobs := registrar.RegisterPeriodicJobs(); len(jobs) != 0 {
		t.Fatalf("RegisterPeriodicJobs() = %d jobs, want 0", len(jobs))
	}
}
