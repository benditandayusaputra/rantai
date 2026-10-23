export type Label = 'resmi' | 'proksi' | 'estimasi' | 'simulasi';

export interface Angka {
	nilai: number;
	satuan: string;
	label: Label;
	provenance_id: string;
}

export interface Meta {
	page: number;
	per_page: number;
	total: number;
	total_pages: number;
}

export interface Amplop<T> {
	code: number;
	message: string;
	data: T;
	meta?: Meta;
}

export interface Provenance {
	id: string;
	sumber: string;
	url: string;
	diambil_pada: string;
	label: Label;
	hash_konten: string;
	catatan: string;
}

export type Tingkat = 'provinsi' | 'kota' | 'kecamatan' | 'kelurahan';

export interface RingkasanWilayah {
	kode: string;
	nama: string;
	tingkat: Tingkat;
}

export interface WilayahDitemukan extends RingkasanWilayah {
	hierarki: RingkasanWilayah[];
}

export interface Teks {
	id: string;
	en: string;
}

export type KelasRasio = 'sangat_kurang' | 'kurang' | 'hampir_cukup' | 'cukup' | 'tidak_ada_data';

export interface IndikatorRingkas {
	periode: string;
	kebutuhan: Angka;
	kapasitas: Angka;
	kesenjangan: Angka;
	rasio: Angka | null;
	kelas: KelasRasio;
}

export interface TemaKecamatan {
	kode: string;
	nama: Teks;
	judul: Teks;
	indikator: IndikatorRingkas | null;
}

export interface Anggaran {
	tahun: number;
	jumlah_paket: Angka;
	total_hps: Angka;
}

export interface Komisi {
	kode: string;
	nama: string;
	bidang?: string[];
}

export interface Wakil {
	id: string;
	nama: string;
	partai: string;
	fraksi: string;
	komisi: Komisi | null;
	duduk_di_komisi_pengawas: boolean;
}

export interface Dapil {
	id: number;
	nomor: number;
	nama: string;
	tahun_pemilu: number;
	jumlah_kursi: Angka;
	jumlah_anggota: Angka;
	kecamatan: RingkasanWilayah[];
	anggota: Wakil[];
	provenance_id: string;
}

export interface WakilKecamatan {
	kecamatan: RingkasanWilayah;
	komisi_pengawas: string[];
	dapil: Dapil;
}

export interface RingkasanKecamatan {
	kecamatan: WilayahDitemukan;
	kelurahan_disorot: RingkasanWilayah | null;
	tema: TemaKecamatan[];
	anggaran: Anggaran | null;
	wakilku: WakilKecamatan;
}

export interface Titik {
	lat: number;
	lng: number;
}

export type Keterbatasan =
	| 'spmb_tidak_dibatasi_kecamatan'
	| 'calon_lulusan_tanpa_mi'
	| 'kursi_madrasah_tidak_tersedia'
	| 'tidak_ada_data_lulusan'
	| 'jarak_garis_lurus';

export type MetodeCalonLulusan = 'per_sekolah' | 'agregat_resmi' | 'estimasi_total_sd' | '';

export interface DetailPendidikan {
	kode: string;
	nama: string;
	periode: string;
	kelas: KelasRasio;
	label: Label;
	provenance_id: string;
	metode_calon_lulusan: MetodeCalonLulusan;
	calon_lulusan: Angka | null;
	kursi_negeri: Angka | null;
	kursi_spmb_bersama: Angka | null;
	kursi_madrasah: Angka | null;
	kesenjangan: Angka;
	kesenjangan_dengan_spmb_bersama: Angka | null;
	rasio: Angka | null;
	keterbatasan: Keterbatasan[];
}

export interface Sekolah {
	npsn: string;
	nama: string;
	bentuk: 'SD' | 'MI' | 'SMP' | 'MTs';
	status: 'negeri' | 'swasta';
	lokasi: Titik | null;
	jalur: 'negeri' | 'spmb_bersama' | null;
	daya_tampung: Angka | null;
}
