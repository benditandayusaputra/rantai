package usecase

import (
	"errors"
	"fmt"
)

const ModeDemo = "demo"

var ErrJobDimatikanModeDemo = errors.New("job pengambilan data dimatikan dalam mode demo")

var jobPengambilan = map[string]bool{
	JobWilayah: true, JobKPUDPRD: true, JobLPSEDaftar: true, JobLPSEDetail: true, JobSIRUP: true,
	JobOpentender: true, JobUmpanDPRD: true, JobUmpanMedia: true, JobRingkasUmpan: true, JobDapodik: true,
}

func IzinkanJob(mode, job string) error {
	if mode == ModeDemo && jobPengambilan[job] {
		return fmt.Errorf("%s: %w", job, ErrJobDimatikanModeDemo)
	}
	return nil
}
