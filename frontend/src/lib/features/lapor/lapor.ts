import type {
	KondisiLaporan,
	PaketTerdekat,
	RiwayatLaporan,
	StatusLaporan,
	Titik,
	WilayahDitemukan
} from '$lib/api/tipe';
import { ambilWilayahDiTitik, tangkap, type Fetch } from '$lib/api/klien';
import { ambangAkurasiKasarM, isIzinDitolak, mintaPosisi } from '$lib/features/wilayah/lokasi';

export const totalLangkah = 6;
export const batasFoto = 4;
export const batasCatatan = 500;
export const radiusProyekM = 500;
export const jumlahSaksiDibutuhkan = 3;

export const daftarKondisi: KondisiLaporan[] = [
	'belum_dimulai',
	'sedang_dikerjakan',
	'terhenti',
	'selesai_baik',
	'selesai_rusak'
];

export const paketTidakTerdaftar = '';

export interface FotoDraf {
	id: string;
	berkas: Blob;
	pratinjau: string;
	diambilPada: Date;
}

export interface DrafLaporan {
	titik: Titik | null;
	akurasiM: number;
	paketKode: string | null;
	paketTerdekat: PaketTerdekat[] | null;
	foto: FotoDraf[];
	kondisi: KondisiLaporan | null;
	catatan: string;
	setuju: boolean;
}

export function drafKosong(): DrafLaporan {
	return {
		titik: null,
		akurasiM: 0,
		paketKode: null,
		paketTerdekat: null,
		foto: [],
		kondisi: null,
		catatan: '',
		setuju: false
	};
}

export function bacaLangkah(nilai: string | null): number {
	const langkah = Number(nilai);
	return Number.isInteger(langkah) && langkah >= 1 && langkah <= totalLangkah ? langkah : 1;
}

const syaratLangkah: Record<number, (draf: DrafLaporan) => boolean> = {
	1: (draf) => draf.titik !== null && draf.akurasiM <= ambangAkurasiKasarM,
	2: (draf) => draf.paketKode !== null,
	3: (draf) => draf.foto.length >= 1 && draf.foto.length <= batasFoto,
	4: (draf) => draf.kondisi !== null,
	5: (draf) => draf.catatan.length <= batasCatatan,
	6: (draf) => draf.setuju
};

export function isLangkahLengkap(langkah: number, draf: DrafLaporan): boolean {
	return syaratLangkah[langkah]?.(draf) ?? false;
}

export function langkahTerjauh(draf: DrafLaporan): number {
	let langkah = 1;
	while (langkah < totalLangkah && isLangkahLengkap(langkah, draf)) langkah += 1;
	return langkah;
}

export function susunFormulir(draf: DrafLaporan): FormData {
	const formulir = new FormData();
	formulir.set('lat', String(draf.titik?.lat ?? ''));
	formulir.set('lng', String(draf.titik?.lng ?? ''));
	formulir.set('akurasi', String(Math.round(draf.akurasiM)));
	formulir.set('kondisi', draf.kondisi ?? '');
	formulir.set('catatan', draf.catatan.trim());
	if (draf.paketKode) formulir.set('paket_kode', draf.paketKode);
	draf.foto.forEach((foto, i) => formulir.append('foto', foto.berkas, `bukti-${i + 1}.jpg`));
	return formulir;
}

export function ukuranSkala(
	lebar: number,
	tinggi: number,
	sisiTerpanjang: number
): { lebar: number; tinggi: number } {
	const skala = Math.min(1, sisiTerpanjang / Math.max(lebar, tinggi));
	return {
		lebar: Math.max(1, Math.round(lebar * skala)),
		tinggi: Math.max(1, Math.round(tinggi * skala))
	};
}

export type HasilPosisi =
	| { status: 'akurat'; titik: Titik; akurasiM: number }
	| { status: 'kasar'; titik: Titik; akurasiM: number }
	| { status: 'ditolak' }
	| { status: 'gagal' };

export function nilaiPosisi(titik: Titik, akurasiM: number): HasilPosisi {
	return akurasiM <= ambangAkurasiKasarM
		? { status: 'akurat', titik, akurasiM }
		: { status: 'kasar', titik, akurasiM };
}

export async function cariPosisi(geolokasi: Geolocation | undefined): Promise<HasilPosisi> {
	if (!geolokasi) return { status: 'gagal' };
	try {
		const { coords } = await mintaPosisi(geolokasi);
		return nilaiPosisi({ lat: coords.latitude, lng: coords.longitude }, coords.accuracy);
	} catch (galat) {
		return isIzinDitolak(galat) ? { status: 'ditolak' } : { status: 'gagal' };
	}
}

export type HasilTempat =
	{ status: 'dikenali'; wilayah: WilayahDitemukan } | { status: 'luar' } | { status: 'gagal' };

const kodeDiLuarJakarta = 404;

export async function kenaliTempat(fetchFn: Fetch, titik: Titik): Promise<HasilTempat> {
	const hasil = await tangkap(ambilWilayahDiTitik(fetchFn, titik.lat, titik.lng));
	if (hasil.ok) return { status: 'dikenali', wilayah: hasil.data };
	return { status: hasil.galat.kode === kodeDiLuarJakarta ? 'luar' : 'gagal' };
}

export function isDrafTerisi(draf: DrafLaporan): boolean {
	return draf.foto.length > 0 || draf.kondisi !== null || draf.catatan.trim() !== '';
}

export const jalurStatusUtama: StatusLaporan[] = [
	'DIKIRIM',
	'MENUNGGU_SAKSI',
	'TERVERIFIKASI',
	'DITERUSKAN',
	'DITANGGAPI',
	'SELESAI'
];

export interface PeristiwaLinimasa {
	status: StatusLaporan;
	pada: Date | null;
	catatan: string;
	isTerjadi: boolean;
}

export function susunLinimasa(riwayat: RiwayatLaporan[]): PeristiwaLinimasa[] {
	const terjadi = riwayat.map((r) => ({
		status: r.status_baru,
		pada: new Date(r.dibuat_pada),
		catatan: r.alasan,
		isTerjadi: true
	}));
	const terakhir = terjadi.at(-1)?.status ?? 'DIKIRIM';
	const posisi = jalurStatusUtama.indexOf(terakhir);
	if (posisi < 0) return terjadi;
	const mendatang = jalurStatusUtama.slice(posisi + 1).map((status) => ({
		status,
		pada: null,
		catatan: '',
		isTerjadi: false
	}));
	return [...terjadi, ...mendatang];
}
