package domain

import (
	"cmp"
	"slices"
)

const (
	BatasItemPerIsu       = 5
	BatasUmpanIsuPerWakil = 3
)

type IsuDisuarakan struct {
	Tag         string
	JumlahResmi int
	JumlahMedia int
	Item        []UmpanTampil
}

func (i IsuDisuarakan) JumlahSumber() int {
	return i.JumlahResmi + i.JumlahMedia
}

func KelompokkanIsu(umpanTerbaruDulu []UmpanTampil) []IsuDisuarakan {
	perTag := map[string]*IsuDisuarakan{}
	for _, u := range umpanTerbaruDulu {
		for _, tag := range u.TagIsu {
			if perTag[tag] == nil {
				perTag[tag] = &IsuDisuarakan{Tag: tag, Item: []UmpanTampil{}}
			}
			perTag[tag].tambah(u)
		}
	}
	hasil := make([]IsuDisuarakan, 0, len(perTag))
	for _, isu := range perTag {
		hasil = append(hasil, *isu)
	}
	slices.SortFunc(hasil, func(a, b IsuDisuarakan) int {
		return cmp.Or(cmp.Compare(b.JumlahSumber(), a.JumlahSumber()), cmp.Compare(a.Tag, b.Tag))
	})
	return hasil
}

func (i *IsuDisuarakan) tambah(u UmpanTampil) {
	if u.JenisSumber == SumberUmpanResmiDPRD {
		i.JumlahResmi++
	} else {
		i.JumlahMedia++
	}
	if len(i.Item) < BatasItemPerIsu {
		i.Item = append(i.Item, u)
	}
}
