import AxeBuilder from '@axe-core/playwright';
import { expect as expectDefault, test } from '@playwright/test';

const expect = expectDefault.configure({ timeout: 15_000 });

test('layanan kabupaten/kota menyajikan data sebelum penjelasan, dampak, dan kesimpulan', async ({
	page
}, info) => {
	await page.goto('/id/kota/31.72');
	await expect(page.getByRole('heading', { level: 1 })).toContainText(/Jakarta Utara/i);
	await expect(
		page
			.getByRole('table', { name: 'Calon lulusan dan kursi negeri per kecamatan' })
			.getByRole('row')
	).toHaveCount(7);
	await expect(page.getByRole('heading', { name: 'Sebaran sekolah antarkecamatan' })).toBeVisible();
	const urutan = await page.locator('main h2').allTextContents();
	expect(urutan.indexOf('Sebaran sekolah antarkecamatan')).toBeLessThan(
		urutan.indexOf('Bagaimana membaca data ini?')
	);
	expect(urutan.indexOf('Bagaimana membaca data ini?')).toBeLessThan(
		urutan.indexOf('Apa dampaknya pada akses siswa?')
	);
	expect(urutan.indexOf('Apa dampaknya pada akses siswa?')).toBeLessThan(
		urutan.indexOf('Kesimpulan dari data yang tersedia')
	);
	await page.getByRole('button', { name: 'Negeri dan swasta', exact: true }).click();
	await expect(
		page.getByRole('button', { name: 'Negeri dan swasta', exact: true })
	).toHaveAttribute('aria-pressed', 'true');
	expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth)).toBe(true);
	await page.screenshot({ path: info.outputPath('layanan-kota.png'), fullPage: true });
	const audit = await new AxeBuilder({ page })
		.withTags(['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa'])
		.analyze();
	expect(audit.violations.filter((v) => v.impact === 'serious' || v.impact === 'critical')).toEqual(
		[]
	);
	await page.getByRole('link', { name: 'Kelapa Gading', exact: true }).first().click();
	await expect(page).toHaveURL(/kecamatan\/31.72.06\?bagian=tema/);
	await expect(
		page.getByRole('link', { name: 'Lihat analisis kabupaten/kota', exact: true })
	).toBeVisible();
});

test('kabupaten Kepulauan Seribu, pergantian jenjang, dan bahasa tetap tersedia', async ({
	page
}) => {
	await page.goto('/id/kota/31.01');
	await expect(page.getByRole('heading', { level: 1 })).toContainText(/Seribu/i);
	await page.getByRole('link', { name: 'SMA/SMK', exact: true }).click();
	await expect(page).toHaveURL(/tema=pendidikan_sma_smk/);
	await expect(page.getByRole('heading', { name: 'Daya tampung masuk SMA/SMK' })).toBeVisible();
	await page.getByLabel('Pilih kabupaten/kota', { exact: true }).selectOption('31.72');
	await expect(page).toHaveURL(/kota\/31.72\?tema=pendidikan_sma_smk/);
	await page.goto('/en/kota/31.72');
	await expect(
		page.getByRole('heading', { name: 'School distribution across districts' })
	).toBeVisible();
});
