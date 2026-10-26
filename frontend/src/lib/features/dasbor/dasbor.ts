import {
	ambilDasborLayanan,
	ambilDasborTanggungJawab,
	ambilDasborUang,
	ambilPetaTema,
	type Fetch
} from '$lib/api/klien';
import type {
	Angka,
	DasborLayanan,
	DasborTanggungJawab,
	DasborUang,
	GeometriPeta,
	KelasRasio,
	PetaTema,
	Porsi,
	RekapKursi,
	Seri
} from '$lib/api/tipe';

export const daftarLensa = ['layanan', 'uang', 'tanggung-jawab'] as const;
export type Lensa = (typeof daftarLensa)[number];

export const temaDasbor = 'pendidikan';

export type DataDasbor =
	| { lensa: 'layanan'; dasbor: DasborLayanan; peta: PetaTema }
	| { lensa: 'uang'; dasbor: DasborUang; peta: PetaTema }
	| { lensa: 'tanggung-jawab'; dasbor: DasborTanggungJawab };

export function bacaLensa(nilai: string | null): Lensa {
	return daftarLensa.find((l) => l === nilai) ?? 'layanan';
}

const polaTahun = /^\d{4}$/;

export function bacaTahun(nilai: string | null): number | undefined {
	return nilai && polaTahun.test(nilai) ? Number(nilai) : undefined;
}

export function tautanLensa(url: URL, lensa: Lensa): string {
	const tujuan = new URL(url);
	tujuan.searchParams.delete('tahun');
	if (lensa === 'layanan') tujuan.searchParams.delete('lensa');
	else tujuan.searchParams.set('lensa', lensa);
	return tujuan.pathname + tujuan.search;
}

export function tautanTahun(url: URL, tahun: number): string {
	const tujuan = new URL(url);
	tujuan.searchParams.set('tahun', String(tahun));
	return tujuan.pathname + tujuan.search;
}

export async function muatDasbor(
	fetchFn: Fetch,
	lensa: Lensa,
	tahun?: number,
	temaLayanan: string = temaDasbor
): Promise<DataDasbor> {
	if (lensa === 'tanggung-jawab') {
		return { lensa, dasbor: await ambilDasborTanggungJawab(fetchFn, tahun) };
	}
	if (lensa === 'uang') {
		const [dasbor, petaTema] = await Promise.all([
			ambilDasborUang(fetchFn, tahun),
			ambilPetaTema(fetchFn, temaDasbor)
		]);
		return { lensa, dasbor, peta: petaTema };
	}
	const [dasbor, petaTema] = await Promise.all([
		ambilDasborLayanan(fetchFn, temaLayanan),
		ambilPetaTema(fetchFn, temaLayanan)
	]);
	return { lensa, dasbor, peta: petaTema };
}

export function angkaDari(
	rujukan: Pick<Seri, 'label' | 'provenance_id'>,
	nilai: number,
	satuan: string
): Angka | null {
	if (!rujukan.label || !rujukan.provenance_id) return null;
	return { nilai, satuan, label: rujukan.label, provenance_id: rujukan.provenance_id };
}

export function jumlahkan<T>(daftar: T[], nilaiDari: (item: T) => number): number {
	return daftar.reduce((total, item) => total + nilaiDari(item), 0);
}

export function daftarTahun(seri: Seri): number[] {
	return seri.isi.map((p) => Number(p.kunci)).filter(Number.isInteger);
}

export function porsiTahun(seri: Seri, tahun: number): Porsi | undefined {
	return seri.isi.find((p) => p.kunci === String(tahun));
}

export function tahunTersedia(data: DataDasbor): { aktif: number; daftar: number[] } | null {
	if (data.lensa === 'uang')
		return { aktif: data.dasbor.tahun, daftar: daftarTahun(data.dasbor.tahunan) };
	if (data.lensa === 'tanggung-jawab') {
		return { aktif: data.dasbor.tahun, daftar: daftarTahun(data.dasbor.kegiatan_tahunan) };
	}
	return null;
}

const awalanNamaKota = /^(Kota|Kabupaten) Administrasi /;

export function namaWilayahPendek(nama: string): string {
	return nama.replace(awalanNamaKota, '');
}

export function hurufAwalBesar(teks: string): string {
	return teks.charAt(0).toUpperCase() + teks.slice(1);
}

const persenPenuh = 100;

export function rataPerBulanTercatat(bulanan: Porsi[]): number {
	const tercatat = bulanan.filter((p) => p.jumlah > 0);
	return tercatat.length === 0 ? 0 : jumlahkan(tercatat, (p) => p.jumlah) / tercatat.length;
}

const langkahRapi = [1, 1.2, 1.5, 2, 2.5, 3, 4, 5, 6, 8, 10];

export function batasRapi(maksimum: number): number {
	if (maksimum <= 0) return 1;
	const pangkat = 10 ** Math.floor(Math.log10(maksimum));
	const langkah = langkahRapi.find((l) => l * pangkat >= maksimum) ?? 10;
	return langkah * pangkat;
}

export function persenDari(nilai: number, maksimum: number): number {
	return maksimum > 0 ? Math.min(persenPenuh, (nilai / maksimum) * persenPenuh) : 0;
}

export function ambangKuantil(nilai: number[], jumlahKelas: number): number[] {
	const urut = nilai.filter((n) => n > 0).sort((a, b) => a - b);
	if (urut.length === 0) return [];
	const ambang = Array.from({ length: jumlahKelas - 1 }, (_, i) => {
		const posisi = Math.floor(((i + 1) * urut.length) / jumlahKelas);
		return urut[Math.min(posisi, urut.length - 1)] ?? 0;
	});
	return [...new Set(ambang)];
}

export const tanpaNilai = -1;

export function tingkatKuantil(nilai: number, ambang: number[]): number {
	if (nilai <= 0) return tanpaNilai;
	return ambang.filter((a) => nilai >= a).length;
}

export const urutanKelas: KelasRasio[] = ['sangat_kurang', 'kurang', 'hampir_cukup', 'cukup'];
export const kotakKelasPeta: Record<KelasRasio, string> = {
	sangat_kurang: 'bg-gap-1',
	kurang: 'bg-gap-2',
	hampir_cukup: 'bg-gap-3',
	cukup: 'bg-gap-4',
	tidak_ada_data: 'border-ctrl border-dashed border-line-strong bg-sunken'
};
export const kelasKurang: KelasRasio[] = ['sangat_kurang', 'kurang'];
export const kelasCukup: KelasRasio[] = ['hampir_cukup', 'cukup'];

export function hitungKelas(rekap: RekapKursi, daftarKelas: KelasRasio[]): number {
	return jumlahkan(daftarKelas, (kelas) => rekap.per_kelas[kelas] ?? 0);
}

export function totalKecamatan(rekap: RekapKursi): number {
	return jumlahkan(Object.values(rekap.per_kelas), (n) => n ?? 0);
}

type Posisi = [number, number];
type Poligon = Posisi[][];

export interface AreaPeta {
	kode: string;
	nama: string;
	d: string;
	tengah: Posisi;
}

export interface Kotak {
	x: number;
	y: number;
	lebar: number;
	tinggi: number;
}

export interface SusunanPeta {
	lebar: number;
	tinggi: number;
	area: AreaPeta[];
	sisipan: Kotak | null;
}

export interface IsiArea {
	kelas: string;
	opasitas?: number;
	keterangan: string;
}

export interface FiturWilayah {
	geometry: GeometriPeta;
	properties: { kode: string; nama: string };
}

interface Batas {
	barat: number;
	timur: number;
	utara: number;
	selatan: number;
}

export const ukuranPeta = {
	lebarUtama: 520,
	lebarSisipan: 96,
	jarak: 24,
	tinggiPetak: 64,
	celahPetak: 8
};
const awalanKepulauanSeribu = '31.01';
const derajatKeRadian = Math.PI / 180;
const desimalJalur = 1;

function poligonDari(geometri: GeometriPeta): Poligon[] {
	if (geometri.type === 'Polygon') return [geometri.coordinates as Poligon];
	if (geometri.type === 'MultiPolygon') return geometri.coordinates as Poligon[];
	return [];
}

function batasDari(semuaPoligon: Poligon[]): Batas {
	const titik = semuaPoligon.flat(2);
	const lng = titik.map((p) => p[0]);
	const lat = titik.map((p) => p[1]);
	return {
		barat: Math.min(...lng),
		timur: Math.max(...lng),
		utara: Math.max(...lat),
		selatan: Math.min(...lat)
	};
}

function faktorLintang(batas: Batas): number {
	return Math.cos(((batas.utara + batas.selatan) / 2) * derajatKeRadian);
}

function rasioTinggi(batas: Batas): number {
	return (batas.utara - batas.selatan) / ((batas.timur - batas.barat) * faktorLintang(batas));
}

function proyektor(batas: Batas, kotak: Kotak): (p: Posisi) => Posisi {
	const kx = faktorLintang(batas);
	const skala = Math.min(
		kotak.lebar / ((batas.timur - batas.barat) * kx),
		kotak.tinggi / (batas.utara - batas.selatan)
	);
	return ([lng, lat]) => [
		kotak.x + (lng - batas.barat) * kx * skala,
		kotak.y + (batas.utara - lat) * skala
	];
}

function areaDari(fitur: FiturWilayah, proyeksi: (p: Posisi) => Posisi): AreaPeta {
	const poligon = poligonDari(fitur.geometry);
	const cincin = poligon.flat().map((c) => c.map(proyeksi));
	const d = cincin
		.map((c) => `M${c.map((p) => p.map((n) => n.toFixed(desimalJalur)).join(',')).join('L')}Z`)
		.join('');
	const terluar = cincin.reduce((a, b) => (b.length > a.length ? b : a), []);
	const tengah: Posisi = [
		jumlahkan(terluar, (p) => p[0]) / terluar.length,
		jumlahkan(terluar, (p) => p[1]) / terluar.length
	];
	return { kode: fitur.properties.kode, nama: fitur.properties.nama, d, tengah };
}

function kelompok(fitur: FiturWilayah[]) {
	const isSeribu = (f: FiturWilayah) => f.properties.kode.startsWith(awalanKepulauanSeribu);
	return { daratan: fitur.filter((f) => !isSeribu(f)), seribu: fitur.filter(isSeribu) };
}

function petakSisipan(fitur: FiturWilayah[]): { area: AreaPeta[]; sisipan: Kotak } {
	const { lebarSisipan, tinggiPetak, celahPetak } = ukuranPeta;
	const lebarPetak = lebarSisipan - 2 * celahPetak;
	const urut = [...fitur].sort((a, b) => a.properties.kode.localeCompare(b.properties.kode));
	const area = urut.map((f, i): AreaPeta => {
		const y = celahPetak + i * (tinggiPetak + celahPetak);
		return {
			kode: f.properties.kode,
			nama: f.properties.nama,
			d: `M${celahPetak},${y}h${lebarPetak}v${tinggiPetak}h${-lebarPetak}Z`,
			tengah: [celahPetak + lebarPetak / 2, y + tinggiPetak / 2]
		};
	});
	const tinggi = celahPetak + urut.length * (tinggiPetak + celahPetak);
	return { area, sisipan: { x: 0, y: 0, lebar: lebarSisipan, tinggi } };
}

function batasKelompok(fitur: FiturWilayah[]): Batas {
	return batasDari(fitur.flatMap((f) => poligonDari(f.geometry)));
}

export function susunPeta(semuaFitur: FiturWilayah[]): SusunanPeta {
	const { daratan, seribu } = kelompok(semuaFitur);
	if (daratan.length === 0)
		return { lebar: ukuranPeta.lebarUtama, tinggi: 0, area: [], sisipan: null };
	const { lebarUtama, lebarSisipan, jarak } = ukuranPeta;
	const batasDaratan = batasKelompok(daratan);
	const tinggi = lebarUtama * rasioTinggi(batasDaratan);
	const geser = seribu.length > 0 ? lebarSisipan + jarak : 0;
	const utama: Kotak = { x: geser, y: 0, lebar: lebarUtama, tinggi };
	const area = daratan.map((f) => areaDari(f, proyektor(batasDaratan, utama)));
	if (seribu.length === 0) return { lebar: lebarUtama, tinggi, area, sisipan: null };
	const petak = petakSisipan(seribu);
	return {
		lebar: geser + lebarUtama,
		tinggi,
		area: [...area, ...petak.area],
		sisipan: petak.sisipan
	};
}
