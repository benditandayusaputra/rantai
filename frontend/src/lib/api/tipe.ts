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
