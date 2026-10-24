package domain

import "time"

type Label string

const (
	LabelResmi    Label = "resmi"
	LabelProksi   Label = "proksi"
	LabelEstimasi Label = "estimasi"
	LabelSimulasi Label = "simulasi"
)

type SumberData string

const (
	SumberKemendagri SumberData = "kemendagri"
	SumberRantai     SumberData = "rantai"
	SumberDapodik    SumberData = "dapodik"
	SumberDisdik     SumberData = "disdik"
	SumberSimulasi   SumberData = "simulasi"
	SumberInstansi   SumberData = "instansi"
)

type Provenance struct {
	ID          string
	Sumber      SumberData
	URL         string
	DiambilPada time.Time
	Label       Label
	HashKonten  string
	Catatan     string
}

func (l Label) IsValid() bool {
	switch l {
	case LabelResmi, LabelProksi, LabelEstimasi, LabelSimulasi:
		return true
	default:
		return false
	}
}
