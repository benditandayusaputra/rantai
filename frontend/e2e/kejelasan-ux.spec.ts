import AxeBuilder from '@axe-core/playwright';
import { expect, test } from '@playwright/test';

test.use({ storageState: { cookies: [], origins: [] } });

const wilayah = '31.72.06';

test('kunjungan pertama menjelaskan manfaat dan tidak menghalangi pencarian dengan tur', async ({
	page
}, info) => {
	await page.goto('/id');
	await expect(page.getByRole('heading', { level: 1 })).toContainText('Kenali kondisi wilayahmu');
	await expect(page.getByRole('combobox')).toBeVisible();
	await expect(
		page.getByText('Khusus DKI Jakarta · minimal 3 huruf.', { exact: false })
	).toBeVisible();
	await expect(page.locator('dialog[open]')).toHaveCount(0);
	await expect(page.locator('main [aria-busy="true"]')).toHaveCount(0, { timeout: 15_000 });
	await page.screenshot({ path: info.outputPath('beranda-layar-pertama.png') });
	await page.screenshot({ path: info.outputPath('beranda.png'), fullPage: true });
	const audit = await new AxeBuilder({ page })
		.withTags(['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa'])
		.analyze();
	expect(audit.violations.filter((v) => v.impact === 'serious' || v.impact === 'critical')).toEqual(
		[]
	);
	await page.getByRole('button', { name: 'Cara memakai Rantai', exact: true }).first().click();
	await expect(page.getByRole('dialog')).toBeVisible();
});

test('ketik wilayah lalu Enter membuka kecamatan yang tepat', async ({ page }) => {
	await page.goto('/id');
	await page.getByRole('combobox').fill('Kelapa Gading');
	await expect(page.getByRole('option', { name: /Kelapa Gading Barat/ })).toBeVisible();
	await page.getByRole('combobox').press('Enter');
	await expect(page).toHaveURL(new RegExp(`/kecamatan/${wilayah}$`));
	await expect(page.getByRole('heading', { level: 1 })).toContainText('Kelapa Gading');
});

test('ringkasan angka tersambung ke aspirasi dengan wilayah dan masalah terpilih', async ({
	page
}, info) => {
	await page.goto(`/id/kecamatan/${wilayah}`);
	await expect(page.getByText(/kursi SMP negeri tersedia untuk setiap 100/)).toBeVisible();
	await page.screenshot({ path: info.outputPath('kecamatan.png'), fullPage: true });
	await page
		.getByRole('link', { name: 'Susun aspirasi dari data ini', exact: true })
		.first()
		.click();
	await expect(page).toHaveURL(/masalah=kursi_smp_kurang/);
	await expect(page.getByRole('heading', { name: 'Draf aspirasi siap' })).toBeVisible();
	await expect(page.getByText('Belum dikirim.', { exact: false })).toBeVisible();
	await expect(page.getByRole('heading', { name: 'Pilihan kanal penyampaian' })).toBeVisible();
	await expect(page.getByRole('link', { name: 'Unduh PDF', exact: true })).toHaveAttribute(
		'href',
		/31.72.06/
	);
	await page.screenshot({ path: info.outputPath('aspirasi.png'), fullPage: true });
});

test('Tanya dari navigasi global tetap membawa wilayah aktif', async ({ page }) => {
	await page.goto(`/id/kecamatan/${wilayah}`);
	await page
		.getByRole('navigation', { name: 'Navigasi utama', exact: true })
		.getByRole('link', { name: 'Tanya', exact: true })
		.filter({ visible: true })
		.click();
	await expect(page).toHaveURL(new RegExp(`/tanya\\?w=${wilayah}$`));
	await expect(
		page.locator('main').getByText('Konteks: Kelapa Gading', { exact: true })
	).toBeVisible();
});

test('Cek Sekolahku memakai wilayah asal dan menjelaskan titik tengah wilayah', async ({
	page
}) => {
	await page.goto(`/id/kecamatan/${wilayah}`);
	await page
		.getByRole('link', { name: /^Cek Sekolahku/ })
		.first()
		.click();
	await expect(page).toHaveURL(new RegExp(`w=${wilayah}`));
	await expect(
		page.getByText('Jarak dihitung dari titik tengah wilayah', { exact: false })
	).toBeVisible();
	await expect(
		page.getByRole('heading', { name: 'Tekanan kursi dalam radius 5 km' })
	).toBeVisible();
});

test('ponsel dan desktop tidak memiliki konten melebar keluar layar', async ({ page }) => {
	for (const path of [
		'/id',
		`/id/kecamatan/${wilayah}`,
		`/id/aspirasi?w=${wilayah}&tema=pendidikan&masalah=kursi_smp_kurang`
	]) {
		await page.goto(path);
		await expect(page.locator('main h1')).toBeVisible();
		expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth)).toBe(
			true
		);
	}
});

test('layout tetap terbaca pada ponsel kecil, tablet, dan batas navigasi desktop', async ({
	page
}, info) => {
	test.skip(info.project.name !== 'desktop');
	for (const width of [320, 768, 1024]) {
		await page.setViewportSize({ width, height: 900 });
		for (const path of ['/id', `/id/kecamatan/${wilayah}`]) {
			await page.goto(path);
			await expect(page.locator('main h1')).toBeVisible();
			expect(
				await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth),
				`${width}: ${path}`
			).toBe(true);
		}
	}
});
