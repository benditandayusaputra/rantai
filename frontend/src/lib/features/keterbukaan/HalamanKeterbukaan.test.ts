import { render, screen, within } from '@testing-library/svelte';
import { describe, expect, it, vi } from 'vitest';
import type { DataBelumDibuka } from '$lib/api/tipe';
import { t } from '$lib/i18n';
import Halaman from '../../../routes/[lang=lang]/data-belum-dibuka/+page.svelte';

vi.mock('$app/state', () => ({
	page: { url: new URL('http://rantai.uji/id/data-belum-dibuka?minta=kehadiran_per_anggota') }
}));

function data(kode: string, nama: string, pengganti: string): DataBelumDibuka {
	return {
		kode,
		nama: { id: nama, en: nama },
		alasan: { id: `Alasan ${nama}`, en: '' },
		pengganti: { id: pengganti, en: '' },
		fitur_pengganti: ''
	};
}

const daftar = {
	ppid: {
		nama: 'PPID DPRD Provinsi DKI Jakarta',
		alamat: 'Jl. Kebon Sirih No. 18',
		url: 'https://ppid.uji'
	},
	data: [
		data('suara_per_anggota', 'Suara per anggota', 'Sikap Fraksi'),
		data('pokir_per_anggota', 'Pokok pikiran per anggota', 'Aspirasi ke Anggaran'),
		data('kehadiran_per_anggota', 'Kehadiran per anggota', 'Permohonan informasi dari halaman ini'),
		data('laporan_reses_per_anggota', 'Laporan reses per anggota', 'Aspirasi Reses')
	]
};

describe('Halaman Data yang Belum Dibuka', () => {
	it('menampilkan keempat data dengan penggantinya dan surat untuk data terpilih', async () => {
		render(Halaman, {
			data: {
				pilihan: { kode: 'kehadiran_per_anggota', bahasa: 'id' },
				halaman: Promise.resolve({
					ok: true,
					data: {
						daftar,
						dokumen: { teks: 'Undang-Undang Nomor 14 Tahun 2008', pdf_url: '/api/v1/x.pdf' }
					}
				})
			}
		} as never);

		const kartu = await screen.findAllByRole('listitem');
		const daftarData = kartu.filter((k) => within(k).queryByRole('heading', { level: 2 }));
		expect(daftarData).toHaveLength(4);
		for (const item of daftar.data) {
			expect(screen.getByRole('heading', { name: item.nama.id })).toBeVisible();
			expect(screen.getByText(item.pengganti.id)).toBeVisible();
		}
		expect(screen.getAllByRole('link', { name: t('unopened.request') })).toHaveLength(4);
		const surat = screen.getByRole('region', { name: /Surat permohonan: Kehadiran per anggota/ });
		expect(surat).toHaveTextContent('Undang-Undang Nomor 14 Tahun 2008');
		expect(surat).toHaveTextContent('PPID DPRD Provinsi DKI Jakarta');
	});
});
