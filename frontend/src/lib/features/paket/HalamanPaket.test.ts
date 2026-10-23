import { render, screen, within } from '@testing-library/svelte';
import { describe, expect, it, vi } from 'vitest';
import type { Angka, DetailPaket, PaketRingkas } from '$lib/api/tipe';
import PackageCard from '$lib/components/PackageCard.svelte';
import HalamanPaket from './HalamanPaket.svelte';

vi.mock('$lib/features/provenance/provenance', async (asli) => ({
	...(await asli<typeof import('$lib/features/provenance/provenance')>()),
	muatProvenance: vi.fn(() => new Promise(() => {}))
}));

const rupiah = (nilai: number): Angka => ({
	nilai,
	satuan: 'rupiah',
	label: 'simulasi',
	provenance_id: 'prv_sim'
});

const lokasi = {
	tingkat: 'kecamatan',
	kode_wilayah: '31.72.03',
	nama_wilayah: 'Cilincing'
} as const;

const paketSimulasi = {
	kode: 'simulasi-1',
	nama: 'Rehabilitasi Ruang Kelas SMP Negeri Contoh 1',
	tahap: 'Pelaksanaan',
	is_simulation: true,
	is_batal: false,
	sektor: 'pendidikan',
	instansi: 'Instansi Contoh Simulasi',
	tahun_anggaran: 2026,
	metode: 'Tender',
	jenis_pengadaan: 'Pekerjaan Konstruksi',
	hps: rupiah(2_408_700_000),
	nilai_kontrak: null,
	pagu: null,
	jumlah_peserta: null,
	url_sumber: '',
	rup: null,
	penyedia: { id: 1, nama: 'CV Contoh Simulasi' },
	lokasi,
	pemeliharaan: {
		berlaku: true,
		status: 'pelaksanaan',
		akhir_pelaksanaan: null,
		akhir_pemeliharaan: null,
		sisa_hari: null,
		label: 'estimasi',
		provenance_id: 'prv_rantai_pemeliharaan'
	},
	rantai: [
		{ kode: 'rencana', status: 'putus', alasan: 'rencana tidak ditemukan' },
		{ kode: 'pemilihan', status: 'ada' },
		{ kode: 'pemenang', status: 'ada' },
		{ kode: 'bukti_warga', status: 'menunggu' }
	],
	sinyal: []
} as unknown as DetailPaket;

describe('HalamanPaket', () => {
	it('memberi pita Data simulasi pada paket simulasi', () => {
		render(HalamanPaket, { paket: paketSimulasi });
		expect(screen.getByText('Data simulasi')).toBeVisible();
	});

	it('menuliskan status mata rantai putus sebagai teks beserta alasannya', () => {
		render(HalamanPaket, { paket: paketSimulasi });
		const rantai = screen.getByRole('list', { name: '2 dari 4 tersambung, 1 putus' });
		const rencana = within(rantai).getAllByRole('listitem')[0] as HTMLElement;
		expect(within(rencana).getByText('Putus')).toBeVisible();
		expect(within(rencana).getByText('rencana tidak ditemukan')).toBeVisible();
	});
});

describe('PackageCard', () => {
	it('memberi pita Data simulasi pada kartu paket simulasi', () => {
		const ringkas = {
			...paketSimulasi,
			penyedia: 'CV Contoh Simulasi',
			tema: []
		} as unknown as PaketRingkas;
		render(PackageCard, { paket: ringkas });
		expect(screen.getByText('Data simulasi')).toBeVisible();
		expect(screen.getByRole('link')).toHaveAttribute('href', '/id/paket/simulasi-1');
	});
});
