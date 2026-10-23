import type { Angka, DetailPendidikan, Sekolah } from '$lib/api/tipe';
import type { KunciPesan } from '$lib/i18n';

export interface BarisGrafik {
	kunci: KunciPesan;
	angka: Angka;
}

const barisKursi = [
	['pen.graduates', 'calon_lulusan'],
	['pen.seats.negeri', 'kursi_negeri'],
	['pen.seats.spmb', 'kursi_spmb_bersama'],
	['pen.seats.madrasah', 'kursi_madrasah']
] as const satisfies readonly (readonly [KunciPesan, keyof DetailPendidikan])[];

export function susunBarisGrafik(detail: DetailPendidikan): BarisGrafik[] {
	return barisKursi.flatMap(([kunci, kolom]) => {
		const angka = detail[kolom];
		return angka ? [{ kunci, angka }] : [];
	});
}

export function persenLebar(nilai: number, maksimum: number): number {
	return maksimum > 0 ? Math.round((nilai / maksimum) * 100) : 0;
}

export function kunciJenisSekolah(sekolah: Sekolah): KunciPesan {
	if (sekolah.bentuk === 'MTs') return 'pen.school.madrasah';
	if (sekolah.jalur === 'spmb_bersama') return 'pen.school.spmb';
	return sekolah.status === 'negeri' ? 'pen.school.negeri' : 'pen.school.swasta';
}
