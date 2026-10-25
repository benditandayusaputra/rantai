package domain

import "testing"

func TestSusunKredit(t *testing.T) {
	kasus := []struct {
		nama, pembuat, lisensi, harap string
	}{
		{"pembuat dan lisensi", "Gunawan Kartapranata", "CC BY-SA 3.0", "Gunawan Kartapranata · CC BY-SA 3.0"},
		{"tanpa pembuat", "  ", "CC0", "CC0"},
		{"tanpa lisensi", "Anonim", "", "Anonim"},
		{"kosong", "", "", ""},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			if hasil := SusunKredit(k.pembuat, k.lisensi); hasil != k.harap {
				t.Errorf("SusunKredit(%q, %q) = %q, harap %q", k.pembuat, k.lisensi, hasil, k.harap)
			}
		})
	}
}
