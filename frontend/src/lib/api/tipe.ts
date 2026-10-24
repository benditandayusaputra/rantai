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

export interface FiturPeta<G, P> {
	type: 'Feature';
	geometry: G;
	properties: P;
}

export interface KoleksiFitur<F> {
	type: 'FeatureCollection';
	features: F[];
}

export interface PropertiTema {
	kode: string;
	nama: string;
	periode: string;
	kelas: KelasRasio;
	label: Label;
	provenance_id: string;
	kebutuhan: Angka;
	kapasitas: Angka;
	kesenjangan: Angka;
	rasio: Angka | null;
}

export type GeometriPeta = { type: string; coordinates: unknown };

export type PetaTema = KoleksiFitur<FiturPeta<GeometriPeta, PropertiTema>>;

export type Sektor = 'pendidikan' | 'air' | 'jalan' | '';

export type PropertiPaket =
	| { jenis_lokasi: 'titik'; kode: string; nama: string; sektor: Sektor; is_simulation: boolean }
	| { jenis_lokasi: 'area'; kode_kecamatan: string; nama_kecamatan: string; paket: string[] };

export type PetaPaket = KoleksiFitur<FiturPeta<GeometriPeta, PropertiPaket>>;

export interface SekolahTerdekat extends Sekolah {
	jarak_m: number;
}

export interface SekitarTitik {
	titik: Titik;
	radius_m: number;
	smp_negeri: SekolahTerdekat[];
	spmb_bersama: SekolahTerdekat[];
	madrasah: SekolahTerdekat[];
	tekanan_radius: {
		siswa_kelas_6: Angka | null;
		kursi_negeri: Angka | null;
		rasio: Angka | null;
	};
	keterbatasan: Keterbatasan[];
}

export interface LokasiPaket {
	tingkat: Tingkat | '';
	kode_wilayah: string;
	nama_wilayah: string;
	titik?: Titik;
	rt?: number;
	rw?: number;
}

export interface PaketRingkas {
	kode: string;
	nama: string;
	tahap: string;
	jenis_pengadaan: string;
	metode: string;
	tahun_anggaran: number;
	hps: Angka | null;
	nilai_kontrak: Angka | null;
	sektor: Sektor;
	tema: string[];
	lokasi: LokasiPaket;
	penyedia: string;
	is_simulation: boolean;
	provenance_id: string;
}

export type StatusMata = 'ada' | 'putus' | 'menunggu' | 'estimasi';
export type KodeMata =
	| 'rencana'
	| 'pemilihan'
	| 'pemenang'
	| 'pelaksanaan'
	| 'lokasi'
	| 'serah_terima'
	| 'masa_pemeliharaan'
	| 'bukti_warga';

export interface MataRantai {
	kode: KodeMata;
	status: StatusMata;
	alasan?: string;
}

export type JenisSinyal = 'S-01' | 'S-02' | 'S-03' | 'S-04' | 'S-05';

export interface Sinyal {
	jenis: JenisSinyal;
	alasan: string;
	pembanding: Record<string, unknown>;
}

export type StatusPemeliharaan =
	'pelaksanaan' | 'pemeliharaan' | 'pemeliharaan_berakhir' | 'tidak_dapat_diestimasi';

export interface Pemeliharaan {
	berlaku: boolean;
	status: StatusPemeliharaan;
	akhir_pelaksanaan: string | null;
	akhir_pemeliharaan: string | null;
	sisa_hari: Angka | null;
	label: 'resmi' | 'estimasi';
	provenance_id: string;
}

export interface RUP {
	kode: string;
	nama: string;
	sumber_dana: string;
	pagu: Angka | null;
	detail_tersedia: boolean;
	pelaksanaan_mulai: string | null;
	pelaksanaan_selesai: string | null;
	url: string;
	provenance_id: string;
}

export interface DetailPaket {
	kode: string;
	kode_spse: string;
	nama: string;
	tahap: string;
	jenis_pengadaan: string;
	metode: string;
	jenis_kontrak: string;
	is_tender: boolean;
	is_batal: boolean;
	is_simulation: boolean;
	tahun_anggaran: number;
	instansi: string;
	kategori: string;
	sektor: Sektor;
	tema: string[];
	pagu: Angka | null;
	hps: Angka | null;
	nilai_kontrak: Angka | null;
	jumlah_peserta: Angka | null;
	tanggal_dibuat: string | null;
	tanggal_kontrak: string | null;
	tanggal_serah_terima: string | null;
	url_sumber: string;
	rup: RUP | null;
	penyedia: { id: number; nama: string } | null;
	lokasi: LokasiPaket;
	pemeliharaan: Pemeliharaan;
	rantai: MataRantai[];
	sinyal: Sinyal[];
	provenance_id: string;
}

export interface PengambilanSumber {
	job: string;
	status_terakhir: 'berhasil' | 'gagal' | 'berjalan';
	dijalankan_pada: string;
	berhasil_pada: string | null;
}

export interface InfoSumberData {
	pengambilan: PengambilanSumber[];
	media: { nama: string; domain: string }[];
	atribusi_opentender: string;
}

export interface TanggungJawab {
	masalah: string;
	label: Teks;
	instansi_pelaksana: string;
	unit_wilayah: string;
	komisi: string;
	dasar_hukum: string;
	dasar_hukum_url: string;
}

export interface TemaDetail {
	kode: string;
	nama: Teks;
	tanggung_jawab: TanggungJawab[];
	provenance_id: string;
}

export interface PermintaanAspirasi {
	kode_kecamatan: string;
	tema: string;
	masalah: string;
	bahasa: string;
}

export interface Dokumen {
	teks: string;
	pdf_url: string;
}

export type Peran = 'warga' | 'anggota' | 'instansi' | 'moderator' | 'admin';

export interface Pengguna {
	id: number;
	email: string;
	nama_tampilan: string;
	peran: Peran;
	instansi_id: number | null;
	anggota_id: string | null;
}

export interface Sesi {
	token_akses: string;
	pengguna: Pengguna;
}

export interface PPID {
	nama: string;
	alamat: string;
	url: string;
}

export interface DataBelumDibuka {
	kode: string;
	nama: Teks;
	alasan: Teks;
	pengganti: Teks;
	fitur_pengganti: string;
}

export interface DaftarDataBelumDibuka {
	ppid: PPID;
	data: DataBelumDibuka[];
}

export interface PaketTerdekat {
	kode: string;
	nama: string;
	sektor: string;
	is_simulation: boolean;
	titik: Titik;
	jarak: Angka;
}

export type KondisiLaporan =
	'belum_dimulai' | 'sedang_dikerjakan' | 'terhenti' | 'selesai_baik' | 'selesai_rusak';

export type StatusLaporan =
	| 'DIKIRIM'
	| 'MENUNGGU_SAKSI'
	| 'TERVERIFIKASI'
	| 'DITERUSKAN'
	| 'DITANGGAPI'
	| 'SELESAI'
	| 'DITOLAK';

export interface LaporanDiterima {
	id: number;
	status: StatusLaporan;
}

export interface SaksiLaporan {
	jumlah: number;
	dibutuhkan: number;
}

export interface RiwayatLaporan {
	status_lama: StatusLaporan | null;
	status_baru: StatusLaporan;
	alasan: string;
	dibuat_pada: string;
}

export interface DetailLaporan {
	id: number;
	status: StatusLaporan;
	kondisi: KondisiLaporan;
	catatan: string;
	kode_wilayah: string;
	nama_wilayah: string;
	lokasi_publik: Titik;
	paket: { kode: string; nama: string } | null;
	is_simulation: boolean;
	dibuat_pada: string;
	saksi: SaksiLaporan;
	foto: { url: string; url_thumbnail: string; lebar: number; tinggi: number }[];
	riwayat: RiwayatLaporan[];
}

export interface HasilKonfirmasi {
	status: StatusLaporan;
	saksi: SaksiLaporan;
}

export interface PropertiLaporan {
	id: number;
	status: StatusLaporan;
	kondisi: KondisiLaporan;
	is_simulation: boolean;
}

export interface FiturLaporan {
	type: 'Feature';
	geometry: { type: 'Point'; coordinates: [number, number] };
	properties: PropertiLaporan;
}

export interface PetaLaporan {
	type: 'FeatureCollection';
	features: FiturLaporan[];
}

export interface LaporanModerasi {
	id: number;
	kondisi: KondisiLaporan;
	catatan: string;
	nama_wilayah: string;
	paket_nama: string;
	is_simulation: boolean;
	dibuat_pada: string;
	foto: { id: number; url: string; url_thumbnail: string; disembunyikan: boolean }[];
}

export interface UmpanModerasi {
	id: number;
	jenis_sumber: string;
	media_nama: string;
	judul: string;
	url: string;
	terbit_pada: string;
	anggota: string[];
}

export interface AntreanModerasi {
	laporan: LaporanModerasi[];
	umpan: UmpanModerasi[];
	alasan_tolak: string[];
}

export type AksiModerasi =
	| { aksi: 'verifikasi' }
	| { aksi: 'tolak'; alasan: string; pesan: string }
	| { aksi: 'sembunyikan_foto'; foto_id: number; alasan: string };

export interface TahunPenyedia {
	tahun: number;
	sumber: string;
	jumlah_paket: Angka;
	total_nilai: Angka;
}

export interface RiwayatPenyedia {
	tahun: number;
	pembeli: string;
	nama_paket: string;
	nilai: Angka | null;
	sumber: string;
	provenance_id: string;
}

export interface ProfilPenyedia {
	id: number;
	nama: string;
	alias: string[];
	per_tahun: TahunPenyedia[];
	kategori_utama: string[];
	instansi_pemberi: string[];
	paket: PaketRingkas[];
	riwayat_opentender: RiwayatPenyedia[];
	atribusi: string;
}

export interface DetailWilayah extends WilayahDitemukan {
	alias: string[];
	titik_tengah: Titik;
}

export interface Langganan {
	id: number;
	kode_wilayah: string;
	nama_wilayah: string;
	dibuat_pada: string;
}

export type JenisKejadian =
	| 'paket_baru'
	| 'melewati_jadwal'
	| 'pemeliharaan_berakhir'
	| 'laporan_terverifikasi'
	| 'hak_jawab'
	| 'indikator_diperbarui';

export interface Notifikasi {
	id: number;
	jenis: JenisKejadian;
	kode_kecamatan: string;
	nama_kecamatan: string;
	judul: string;
	tautan: string;
	terjadi_pada: string;
	dibaca: boolean;
	is_simulation: boolean;
}

export interface HalamanNotifikasi {
	item: Notifikasi[];
	belum_dibaca: number;
}

export interface BelanjaKecamatan {
	kecamatan: RingkasanWilayah;
	kode_kota: string;
	indikator: IndikatorRingkas | null;
	belanja: Anggaran | null;
}

export interface BelanjaTema {
	tema: string;
	tahun: number;
	kecamatan: BelanjaKecamatan[];
	luar_kecamatan: { wilayah: RingkasanWilayah; belanja: Anggaran | null }[];
}

export interface ProfilAnggota extends Wakil {
	keterangan_komisi?: 'pimpinan_dprd' | 'tidak_tersedia';
	alat_kelengkapan: { alat_kelengkapan: string; jabatan: string }[];
	dapil: Dapil;
	url_profil: string;
	lhkpn_url: string;
	provenance_id: string;
}

export type JenisSumberUmpan = 'resmi_dprd' | 'media';

export interface ItemUmpan {
	id: number;
	jenis_sumber: JenisSumberUmpan;
	media_nama: string;
	media_domain: string;
	judul: string;
	ringkasan: string | null;
	url: string;
	terbit_pada: string;
	tag_isu: string[];
	dasar_cocok: string;
	provenance_id: string;
}

export interface IsuDisuarakan {
	tag: string;
	jumlah_sumber: number;
	jumlah_rilis_resmi: number;
	jumlah_media: number;
	item: ItemUmpan[];
}

export interface WakilDenganUmpan extends WakilKecamatan {
	umpan_isu_pengawas?: Record<string, ItemUmpan[]>;
}
