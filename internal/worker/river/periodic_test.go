package river

import "testing"

func TestRegisterPeriodicJobsUsesEmptyBaselineSchedule(t *testing.T) {
	if jobs := RegisterPeriodicJobs(); len(jobs) != 0 {
		t.Fatalf("RegisterPeriodicJobs() = %d jobs, want 0", len(jobs))
	}
}
