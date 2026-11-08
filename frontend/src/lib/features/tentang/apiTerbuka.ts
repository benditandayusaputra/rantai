import { susunURL, type Kueri } from '$lib/api/klien';
import type { KunciPesan } from '$lib/i18n';

export interface TitikAkhirTerbuka {
	path: string;
	kueri: Kueri;
	deskripsi: KunciPesan;
}

const kecamatanContoh = '31.72.06';
const kotaContoh = '31.72';
const tahunContoh = 2026;

export const titikAkhirTerbuka: TitikAkhirTerbuka[] = [
	{ path: '/wilayah/cari', kueri: { q: 'kelapa gading' }, deskripsi: 'about.api.search' },
	{ path: `/kecamatan/${kecamatanContoh}/ringkasan`, kueri: {}, deskripsi: 'about.api.summary' },
	{
		path: `/kecamatan/${kecamatanContoh}/paket`,
		kueri: { tahun: tahunContoh },
		deskripsi: 'about.api.packages'
	},
	{
		path: `/tema/pendidikan/kecamatan/${kecamatanContoh}`,
		kueri: {},
		deskripsi: 'about.api.seats'
	},
	{ path: '/jaringan', kueri: { wilayah: kotaContoh }, deskripsi: 'about.api.network' },
	{ path: '/apbd', kueri: { tahun: tahunContoh }, deskripsi: 'about.api.budget' },
	{
		path: '/unduh/paket.csv',
		kueri: { kode_wilayah: kecamatanContoh, tahun: tahunContoh },
		deskripsi: 'about.api.csv'
	},
	{ path: '/sumber-data', kueri: {}, deskripsi: 'about.api.sources' }
];

export function alamatLengkap(titik: TitikAkhirTerbuka, asal: string): string {
	return new URL(susunURL(titik.path, titik.kueri), asal).toString();
}

export function contohPath(titik: TitikAkhirTerbuka): string {
	const url = new URL(alamatLengkap(titik, 'http://contoh.invalid'));
	return `${url.pathname}${url.search}`;
}
