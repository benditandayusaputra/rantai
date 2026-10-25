import { defineConfig } from '@playwright/test';

try {
	process.loadEnvFile('../.env');
} catch {
	process.env.RANTAI_MODE ??= 'demo';
}

const portAPI = 18090;
const portWeb = 4174;
const alamatWeb = `http://localhost:${portWeb}`;

export default defineConfig({
	testDir: 'e2e',
	timeout: 60_000,
	retries: process.env.CI ? 1 : 0,
	reporter: 'list',
	use: { baseURL: alamatWeb, trace: 'retain-on-failure' },
	projects: [
		{
			name: 'ponsel',
			use: { viewport: { width: 390, height: 844 }, isMobile: true, hasTouch: true }
		},
		{ name: 'desktop', use: { viewport: { width: 1440, height: 900 } } }
	],
	webServer: [
		{
			command: 'go run ./cmd/api',
			cwd: '../backend',
			url: `http://localhost:${portAPI}/api/v1/kesehatan`,
			timeout: 180_000,
			reuseExistingServer: !process.env.CI,
			env: { ...process.env, API_PORT: String(portAPI), CORS_ORIGIN: alamatWeb } as Record<
				string,
				string
			>
		},
		{
			command: `npm run dev -- --port ${portWeb} --strictPort`,
			url: alamatWeb,
			timeout: 120_000,
			reuseExistingServer: !process.env.CI,
			env: {
				...process.env,
				PUBLIC_API_URL: `http://localhost:${portAPI}/api/v1`
			} as Record<string, string>
		}
	]
});
