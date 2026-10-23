package usecase

import (
	"errors"
	"testing"
)

func TestIzinkanJob(t *testing.T) {
	kasus := []struct {
		mode, job string
		harap     error
	}{
		{ModeDemo, JobLPSEDaftar, ErrJobDimatikanModeDemo},
		{ModeDemo, JobDapodik, ErrJobDimatikanModeDemo},
		{ModeDemo, JobRingkasUmpan, ErrJobDimatikanModeDemo},
		{ModeDemo, JobIndikator, nil},
		{ModeDemo, JobSkenarioDemo, nil},
		{"live", JobLPSEDaftar, nil},
		{"", JobWilayah, nil},
	}
	for _, k := range kasus {
		t.Run(k.mode+"/"+k.job, func(t *testing.T) {
			if err := IzinkanJob(k.mode, k.job); !errors.Is(err, k.harap) {
				t.Errorf("err = %v, harap %v", err, k.harap)
			}
		})
	}
}
