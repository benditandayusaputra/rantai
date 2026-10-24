import type { AntreanPanel, Peran, StatusLaporan, TargetHakJawab } from '$lib/api/tipe';
import { t, type KunciPesan } from '$lib/i18n';

export type BagianPanel = 'sinyal' | 'laporan' | 'paket' | 'umpan';

export interface Sorotan {
	jenis: TargetHakJawab;
	id: string;
	judul: string;
	catatan: string;
	keterangan: string;
	tautanPublik: string;
	status?: StatusLaporan;
	isPenting: boolean;
	isSimulation: boolean;
}

export const panjangMaksHakJawab = 2000;

const statusPaket: Record<string, KunciPesan> = {
	pelaksanaan: 'maint.status.pelaksanaan',
	pemeliharaan: 'maint.status.pemeliharaan'
};

export function isAkunPanel(peran: Peran | undefined): boolean {
	return peran === 'instansi' || peran === 'anggota';
}

export function bagianUntuk(peran: Peran): BagianPanel[] {
	return peran === 'anggota' ? ['umpan'] : ['sinyal', 'laporan', 'paket'];
}

export function bagianAwal(peran: Peran): BagianPanel {
	return peran === 'anggota' ? 'umpan' : 'sinyal';
}

export function jumlahSegeraBerakhir(antrean: AntreanPanel): number {
	return antrean.paket.filter((p) => p.segera_berakhir).length;
}

export function daftarSorotan(
	antrean: AntreanPanel,
	bagian: BagianPanel,
	bahasa: string
): Sorotan[] {
	switch (bagian) {
		case 'sinyal':
			return antrean.sinyal.map((s) => ({
				jenis: 'sinyal',
				id: String(s.id),
				judul: s.nama_paket,
				catatan: s.alasan,
				keterangan: `${t(`signal.${s.jenis}`)} · ${t('panel.reply.count', { n: s.jumlah_hak_jawab })}`,
				tautanPublik: `/${bahasa}/paket/${s.kode_paket}`,
				isPenting: s.jumlah_hak_jawab === 0,
				isSimulation: false
			}));
		case 'laporan':
			return antrean.laporan.map((l) => ({
				jenis: 'laporan',
				id: String(l.id),
				judul: l.nama_paket,
				catatan: l.catatan,
				keterangan: `${t(`report.cond.${l.kondisi}`)} · ${l.nama_wilayah}`,
				tautanPublik: `/${bahasa}/laporan/${l.id}`,
				status: l.status,
				isPenting: l.status === 'DITERUSKAN',
				isSimulation: l.is_simulation
			}));
		case 'paket':
			return antrean.paket.map((p) => ({
				jenis: 'paket',
				id: p.kode,
				judul: p.nama,
				catatan: '',
				keterangan: p.segera_berakhir
					? t('panel.kpi.soon')
					: t(statusPaket[p.pemeliharaan.status] ?? 'maint.status.pemeliharaan'),
				tautanPublik: `/${bahasa}/paket/${p.kode}`,
				isPenting: p.segera_berakhir,
				isSimulation: p.is_simulation
			}));
		case 'umpan':
			return antrean.umpan.map((u) => ({
				jenis: 'umpan',
				id: String(u.id),
				judul: u.judul,
				catatan: '',
				keterangan: `${u.media_nama} · ${t('panel.reply.count', { n: u.jumlah_hak_jawab })}`,
				tautanPublik: u.url,
				isPenting: u.jumlah_hak_jawab === 0,
				isSimulation: false
			}));
	}
}

export function isSiapTerbit(isi: string, isBerwenang: boolean): boolean {
	const panjang = isi.trim().length;
	return isBerwenang && panjang > 0 && panjang <= panjangMaksHakJawab;
}

export function tanggalLokal(d: Date): string {
	const bulan = String(d.getMonth() + 1).padStart(2, '0');
	const hari = String(d.getDate()).padStart(2, '0');
	return `${d.getFullYear()}-${bulan}-${hari}`;
}
