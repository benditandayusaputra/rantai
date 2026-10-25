import AxeBuilder from '@axe-core/playwright';
import { expect, test, type Page } from '@playwright/test';
import { alamatAPI } from './lingkungan';

const kecamatanUji = '31.72.06';
const tagWCAG = ['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa'];
const dampakBerat = new Set(['serious', 'critical']);

const halamanTetap = [
	'/id',
	`/id/kecamatan/${kecamatanUji}`,
	`/id/kecamatan/${kecamatanUji}/pendidikan`,
	'/id/peta',
	'/id/dasbor',
	'/id/dasbor?lensa=uang',
	'/id/dasbor?lensa=tanggung-jawab',
	'/id/tanya',
	'/id/tentang-data',
	'/id/data-belum-dibuka',
	'/id/cek-sekolah',
	`/id/aspirasi?w=${kecamatanUji}`,
	'/id/masuk',
	'/id/pantau',
	'/id/lapor',
	'/id/panel',
	'/id/moderasi',
	'/en',
	`/en/kecamatan/${kecamatanUji}`
];

async function pelanggaranBerat(page: Page): Promise<string[]> {
	const hasil = await new AxeBuilder({ page }).withTags(tagWCAG).analyze();
	return hasil.violations
		.filter((v) => dampakBerat.has(v.impact ?? ''))
		.map(
			(v) =>
				`${v.id} (${v.impact}): ${v.nodes
					.map((n) => n.target.join(' '))
					.slice(0, 3)
					.join(' | ')}`
		);
}

async function bukaDanAudit(page: Page, path: string) {
	await page.goto(path);
	await page.locator('main').first().waitFor();
	await page.waitForLoadState('networkidle', { timeout: 15_000 }).catch(() => undefined);
	expect(await pelanggaranBerat(page)).toEqual([]);
}

for (const path of halamanTetap) {
	test(`tanpa pelanggaran serius atau kritis: ${path}`, async ({ page }) => {
		await bukaDanAudit(page, path);
	});
}

test('tanpa pelanggaran serius atau kritis: profil anggota dan paket', async ({
	page,
	request
}) => {
	const wakil = await (await request.get(`${alamatAPI}/kecamatan/${kecamatanUji}/wakil`)).json();
	const anggota = wakil.data?.dapil?.anggota?.[0]?.id;
	test.skip(!anggota, 'data wakil belum ada');
	await bukaDanAudit(page, `/id/anggota/${anggota}`);
	const paket = await (await request.get(`${alamatAPI}/kecamatan/${kecamatanUji}/paket`)).json();
	const kode = paket.data?.[0]?.kode;
	if (kode) await bukaDanAudit(page, `/id/paket/${kode}`);
});
