import type { Angka, BentukSekolah, DetailPendidikan, KursiBentuk, Sekolah } from '$lib/api/tipe';
import { tampilkanAngka } from '$lib/features/kecamatan/kecamatan';
import { isKunciPesan, type KunciPesan, type ParameterPesan } from '$lib/i18n';
import { teksJenjang, type TemaPendidikan } from './jenjang';

export interface BarisGrafik {
	kunci: KunciPesan;
	angka: Angka;
	isKebutuhan: boolean;
}

type JalurKursi = 'negeri' | 'spmb';

export function susunBarisGrafik(detail: DetailPendidikan, tema: TemaPendidikan): BarisGrafik[] {
	const teks = teksJenjang[tema];
	const lulusan = detail.calon_lulusan
		? [{ kunci: teks.lulusan, angka: detail.calon_lulusan, isKebutuhan: true }]
		: [];
	const madrasah = detail.kursi_madrasah
		? [{ kunci: teks.kursiMadrasah, angka: detail.kursi_madrasah, isKebutuhan: false }]
		: [];
	return [
		...lulusan,
		...barisKursi(detail.kursi_negeri_per_bentuk, 'negeri', detail.kursi_negeri, teks.kursiNegeri),
		...barisKursi(
			detail.kursi_spmb_bersama_per_bentuk,
			'spmb',
			detail.kursi_spmb_bersama,
			teks.kursiSPMB
		),
		...madrasah
	];
}

function barisKursi(
	perBentuk: KursiBentuk[] | undefined,
	jalur: JalurKursi,
	total: Angka | null,
	kunciTotal: KunciPesan
): BarisGrafik[] {
	if (perBentuk?.length) {
		return perBentuk.map((k) => ({
			kunci: kunciKursiBentuk(jalur, k.bentuk) ?? kunciTotal,
			angka: k.kursi,
			isKebutuhan: false
		}));
	}
	return total ? [{ kunci: kunciTotal, angka: total, isKebutuhan: false }] : [];
}

function kunciKursiBentuk(jalur: JalurKursi, bentuk: BentukSekolah): KunciPesan | null {
	const kunci = `pen.seats.${jalur}.${bentuk}`;
	return isKunciPesan(kunci) ? kunci : null;
}

export function persenLebar(nilai: number, maksimum: number): number {
	return maksimum > 0 ? Math.round((nilai / maksimum) * 100) : 0;
}

const kunciSekolahUmum = {
	SMP: { negeri: 'pen.school.negeri', swasta: 'pen.school.swasta', spmb: 'pen.school.spmb' },
	SMA: {
		negeri: 'pen.school.sma.negeri',
		swasta: 'pen.school.sma.swasta',
		spmb: 'pen.school.sma.spmb'
	},
	SMK: {
		negeri: 'pen.school.smk.negeri',
		swasta: 'pen.school.smk.swasta',
		spmb: 'pen.school.smk.spmb'
	}
} as const satisfies Record<string, Record<'negeri' | 'swasta' | 'spmb', KunciPesan>>;

const kunciMadrasah: Partial<Record<BentukSekolah, KunciPesan>> = {
	MTs: 'pen.school.madrasah',
	MA: 'pen.school.ma'
};

function isBentukUmum(bentuk: BentukSekolah): bentuk is keyof typeof kunciSekolahUmum {
	return bentuk in kunciSekolahUmum;
}

export function kunciJenisSekolah(sekolah: Sekolah): KunciPesan {
	const madrasah = kunciMadrasah[sekolah.bentuk];
	if (madrasah) return madrasah;
	const kunci = kunciSekolahUmum[isBentukUmum(sekolah.bentuk) ? sekolah.bentuk : 'SMP'];
	return sekolah.jalur === 'spmb_bersama' ? kunci.spmb : kunci[sekolah.status];
}

export interface BarisRumus {
	kunci: KunciPesan;
	parameter: ParameterPesan;
	hitungan: string | null;
}

const nilai = (angka: Angka) => tampilkanAngka(angka).nilai;

export function susunRumus(detail: DetailPendidikan): BarisRumus[] {
	const r = detail.rumus;
	const lulusan = detail.calon_lulusan;
	const kursi = detail.kursi_negeri;
	const baris: BarisRumus[] = [
		{
			kunci: 'pen.formula.graduates',
			parameter: { tingkat: r.tingkat_asal, bentuk: r.bentuk_asal },
			hitungan: lulusan ? nilai(lulusan) : null
		},
		{
			kunci: 'pen.formula.seats',
			parameter: { tujuan: r.bentuk_tujuan.join(' + ') },
			hitungan: kursi ? hitunganKursi(detail.kursi_negeri_per_bentuk, kursi) : null
		},
		{
			kunci: 'pen.formula.gap',
			parameter: {},
			hitungan:
				lulusan && kursi
					? `${nilai(lulusan)} − ${nilai(kursi)} = ${nilai(detail.kesenjangan)}`
					: null
		}
	];
	if (detail.rasio && lulusan && kursi) {
		const hitungan = `${nilai(kursi)} ÷ ${nilai(lulusan)} = ${nilai(detail.rasio)}`;
		baris.push({ kunci: 'pen.formula.ratio', parameter: {}, hitungan });
	}
	if (detail.kesenjangan_dengan_spmb_bersama && detail.kursi_spmb_bersama) {
		const hitungan = `${nilai(detail.kesenjangan)} − ${nilai(detail.kursi_spmb_bersama)} = ${nilai(detail.kesenjangan_dengan_spmb_bersama)}`;
		baris.push({ kunci: 'pen.formula.gapSpmb', parameter: {}, hitungan });
	}
	if (r.batas_murid_per_rombel) {
		const parameter = {
			madrasah: r.bentuk_madrasah_tujuan,
			tingkat: r.tingkat_tujuan,
			batas: r.batas_murid_per_rombel
		};
		const hitungan = detail.kursi_madrasah ? nilai(detail.kursi_madrasah) : null;
		baris.push({ kunci: 'pen.formula.madrasah', parameter, hitungan });
	}
	return baris;
}

function hitunganKursi(perBentuk: KursiBentuk[] | undefined, total: Angka): string {
	if (!perBentuk || perBentuk.length < 2) return nilai(total);
	return `${perBentuk.map((k) => nilai(k.kursi)).join(' + ')} = ${nilai(total)}`;
}

export const ukuranGrafikTren = {
	lebar: 600,
	tinggi: 240,
	kiri: 32,
	kanan: 32,
	atas: 36,
	bawah: 44
};
const ruangDiAtasPuncak = 1.15;
const rasioSeimbang = 1;

export function skalaTren<T>(daftar: T[], rasioDari: (item: T) => number) {
	const { lebar, tinggi, kiri, kanan, atas, bawah } = ukuranGrafikTren;
	const puncak = Math.max(rasioSeimbang, ...daftar.map(rasioDari)) * ruangDiAtasPuncak;
	const jarak = daftar.length > 1 ? (lebar - kiri - kanan) / (daftar.length - 1) : 0;
	const y = (nilai: number) => tinggi - bawah - (nilai / puncak) * (tinggi - atas - bawah);
	return {
		titik: daftar.map((item, i) => ({ item, x: kiri + i * jarak, y: y(rasioDari(item)) })),
		ySeimbang: y(rasioSeimbang)
	};
}
