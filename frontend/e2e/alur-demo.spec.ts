import { existsSync, readFileSync, statSync } from 'node:fs';
import { expect, test, type Page } from '@playwright/test';
import { alamatAPI, logAPI } from './lingkungan';

const kodePaketSimulasi = 'simulasi-ruang-kelas-1';
const emailPetugas = 'petugas@demo.rantai.test';
const pertanyaanDemo = 'Kenapa SMP negeri di kecamatanku penuh?';
const masalahDemo = 'Kursi SMP negeri kurang';
const polaKodeMasuk = /kode masuk mode demo: (\d{6})/g;
const jumlahDataBelumDibuka = 4;

interface Kecamatan {
	kode: string;
	nama: string;
	namaPaketSimulasi: string;
}

let kecamatan: Kecamatan | null = null;

test.describe.configure({ mode: 'serial' });

test.beforeAll(async ({ request }) => {
	test.skip(process.env.RANTAI_MODE !== 'demo', 'alur demo hanya berjalan dengan RANTAI_MODE=demo');
	const paket = await request.get(`${alamatAPI}/paket/${kodePaketSimulasi}`);
	test.skip(!paket.ok(), 'paket simulasi belum dimuat; jalankan make seed-demo');
	const detail = (await paket.json()).data;
	const kode = String(detail.lokasi.kode_wilayah).split('.').slice(0, 3).join('.');
	const wilayah = (await (await request.get(`${alamatAPI}/wilayah/${kode}`)).json()).data;
	kecamatan = { kode, nama: wilayah.nama, namaPaketSimulasi: detail.nama };
});

function demo(): Kecamatan {
	if (!kecamatan) throw new Error('skenario demo belum dimuat');
	return kecamatan;
}

function teksTerlihat(page: Page, teks: string | RegExp) {
	return page.getByText(teks).filter({ visible: true }).first();
}

function isLayarLebar(page: Page): boolean {
	return (page.viewportSize()?.width ?? 0) >= 1024;
}

async function bukaBagian(page: Page, bagian: 'tema' | 'anggaran' | 'wakilku') {
	await page.goto(`/id/kecamatan/${demo().kode}?bagian=${bagian}`);
	await page.locator('main h1').waitFor();
}

function ukuranLog(): number {
	return existsSync(logAPI) ? statSync(logAPI).size : 0;
}

function kodeMasukSetelah(posisi: number): string | null {
	if (!existsSync(logAPI)) return null;
	const isiBaru = readFileSync(logAPI, 'utf8').slice(posisi);
	return [...isiBaru.matchAll(polaKodeMasuk)].at(-1)?.[1] ?? null;
}

async function masukSebagai(page: Page, email: string, lanjut: string) {
	const posisiLog = ukuranLog();
	await page.goto(`/id/masuk?lanjut=${encodeURIComponent(lanjut)}`, { waitUntil: 'networkidle' });
	await page.getByLabel('Email').fill(email);
	await page.getByRole('button', { name: 'Kirim kode' }).click();
	let kode: string | null = null;
	await expect
		.poll(() => (kode = kodeMasukSetelah(posisiLog)), {
			message: `kode masuk tidak muncul di ${logAPI}`
		})
		.not.toBeNull();
	await page.getByLabel('Kode masuk').fill(kode ?? '');
	await page.waitForURL(`**${lanjut}`);
}

test('1. warga mencari kecamatan di Jakarta Utara dari beranda', async ({ page }) => {
	const { kode, nama } = demo();
	await page.goto('/id');
	await page.getByRole('combobox').fill(nama);
	await page
		.getByRole('option', { name: new RegExp(nama) })
		.first()
		.click();
	await expect(page).toHaveURL(new RegExp(`/id/kecamatan/${kode}`));
	await expect(page.locator('main h1')).toContainText(nama);
});

test('2. tema Pendidikan menampilkan kursi SMP negeri berlabel resmi', async ({ page }) => {
	await bukaBagian(page, 'tema');
	await expect(page.getByText('Kursi SMP negeri dibanding lulusan SD')).toBeVisible();
	await expect(page.getByRole('button', { name: /^Sumber: .*Resmi/ }).first()).toBeVisible();
	await page.getByRole('link', { name: 'Lihat detail sekolah' }).click();
	await expect(page).toHaveURL(new RegExp(`/kecamatan/${demo().kode}/pendidikan`));
	await expect(page.getByRole('table').first()).toBeVisible();
});

test('3. arus anggaran memuat paket sekolah dan grafik kesenjangan dibanding belanja', async ({
	page
}) => {
	await bukaBagian(page, 'anggaran');
	await expect(page.getByText('Kesenjangan kursi dibanding belanja sekolah')).toBeVisible();
	const tautanPaket = page.getByRole('link', { name: new RegExp(demo().namaPaketSimulasi) });
	await expect(tautanPaket.first()).toBeVisible();
	await tautanPaket.first().click();
	await expect(page).toHaveURL(new RegExp(`/paket/${kodePaketSimulasi}`));
	await expect(page.getByText('Rantai Rupiah')).toBeVisible();
	await expect(teksTerlihat(page, 'Masa pemeliharaan')).toBeVisible();
	await expect(teksTerlihat(page, 'Data simulasi')).toBeVisible();
});

test('4. wakilku menampilkan wakil dapil di Komisi E yang mengawasi Dinas Pendidikan', async ({
	page
}) => {
	await bukaBagian(page, 'wakilku');
	await expect(teksTerlihat(page, /Komisi E/)).toBeVisible();
	await expect(teksTerlihat(page, 'Duduk di komisi yang mengawasi urusan ini')).toBeVisible();
	await page.locator('main a[href*="/id/anggota/"]').first().click();
	await expect(page).toHaveURL(/\/id\/anggota\//);
	await expect(teksTerlihat(page, 'Dinas Pendidikan')).toBeVisible();
	await expect(teksTerlihat(page, 'Isu yang disuarakan')).toBeVisible();
	await expect(teksTerlihat(page, 'Jejak kegiatan')).toBeVisible();
	await expect(teksTerlihat(page, 'Rilis resmi DPRD')).toBeVisible();
});

test('5. Tanya Rantai menjawab dengan rujukan yang membuka sumber', async ({ page }) => {
	await page.goto(`/id/tanya?w=${demo().kode}`);
	await page.getByRole('button', { name: pertanyaanDemo }).click();
	const rujukan = page.getByRole('button', { name: /^Rujukan 1:/ }).first();
	await expect(rujukan).toBeVisible({ timeout: 30_000 });
	await rujukan.click();
	await expect(page.getByRole('dialog').getByText('Sumber angka ini')).toBeVisible();
});

test('6. Paket Aspirasi tersusun dan bisa diunduh sebagai PDF', async ({ page, request }) => {
	await bukaBagian(page, 'wakilku');
	await page.getByRole('link', { name: 'Buat Paket Aspirasi' }).click();
	await expect(page).toHaveURL(/\/id\/aspirasi/);
	await page.getByRole('link', { name: masalahDemo }).click();
	const unduh = page.getByRole('link', { name: 'Unduh PDF' });
	await expect(unduh).toBeVisible();
	const alamatPDF = new URL((await unduh.getAttribute('href')) ?? '', page.url()).toString();
	const pdf = await request.get(alamatPDF);
	expect(pdf.ok()).toBe(true);
	expect(pdf.headers()['content-type']).toContain('application/pdf');
});

test('7. petugas instansi memberi hak jawab yang langsung tampil di halaman publik', async ({
	page
}) => {
	const halamanPaket = `/id/paket/${kodePaketSimulasi}`;
	if (!isLayarLebar(page)) {
		await page.goto(halamanPaket);
		await expect(teksTerlihat(page, 'Hak jawab instansi')).toBeVisible();
		return;
	}
	const jawaban = `Ruang kelas simulasi sudah diperiksa ulang pada uji alur demo ${Date.now()}.`;
	await masukSebagai(page, emailPetugas, '/id/panel');
	await page.getByRole('button', { name: /Paket berjalan/ }).click();
	await page.getByRole('option', { name: new RegExp(demo().namaPaketSimulasi) }).click();
	await page.getByLabel('Isi hak jawab').fill(jawaban);
	await page.getByLabel('Saya berwenang menanggapi atas nama instansi ini.').check();
	await page.getByRole('button', { name: 'Terbitkan hak jawab' }).click();
	await expect(
		page.getByText('Hak jawab terbit dan sudah tampil di halaman publik.')
	).toBeVisible();
	await page.goto(halamanPaket);
	await expect(page.getByText(jawaban)).toBeVisible();
});

test('8. penutup di halaman Data yang Belum Dibuka dengan dasar UU 14/2008', async ({ page }) => {
	await page.goto('/id/data-belum-dibuka');
	await expect(
		page.getByRole('heading', { level: 1, name: 'Data yang belum dibuka' })
	).toBeVisible();
	const permintaan = page.getByRole('link', { name: 'Minta data ini' });
	await expect(permintaan).toHaveCount(jumlahDataBelumDibuka);
	await permintaan.first().click();
	await expect(teksTerlihat(page, /Undang-Undang Nomor 14 Tahun 2008/)).toBeVisible();
});
