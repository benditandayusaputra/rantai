import { describe, expect, it } from 'vitest';
import { saringRespons, susunPermintaanHulu } from './proxy';

const tujuan = { hulu: 'https://api.uji', rahasia: 'rahasia-server' };

describe('susunPermintaanHulu', () => {
	it('meneruskan path, kueri, dan header yang diizinkan saja', () => {
		const masuk = new Request('https://web.uji/api/v1/paket/x?tahun=2026', {
			headers: {
				authorization: 'Bearer t',
				cookie: 'rantai_refresh=r',
				'x-rantai-proxy': 'palsu',
				'x-klien-asli': '1.2.3.4',
				'x-forwarded-for': '9.9.9.9'
			}
		});

		const hulu = susunPermintaanHulu(masuk, 'paket/x', tujuan, '203.0.113.7');

		expect(hulu.url).toBe('https://api.uji/api/v1/paket/x?tahun=2026');
		expect(hulu.headers.get('authorization')).toBe('Bearer t');
		expect(hulu.headers.get('cookie')).toBe('rantai_refresh=r');
		expect(hulu.headers.get('x-rantai-proxy')).toBe('rahasia-server');
		expect(hulu.headers.get('x-klien-asli')).toBe('203.0.113.7');
		expect(hulu.headers.get('x-forwarded-for')).toBeNull();
	});

	it('meneruskan badan pada POST', async () => {
		const masuk = new Request('https://web.uji/api/v1/auth/otp/kirim', {
			method: 'POST',
			body: '{"email":"a@b.id"}',
			headers: { 'content-type': 'application/json' }
		});

		const hulu = susunPermintaanHulu(masuk, 'auth/otp/kirim', tujuan, '203.0.113.7');

		expect(hulu.method).toBe('POST');
		expect(await hulu.text()).toBe('{"email":"a@b.id"}');
	});
});

describe('saringRespons', () => {
	it('mempertahankan semua set-cookie dan membuang header lain', () => {
		const header = new Headers({ 'content-type': 'application/json', server: 'caddy' });
		header.append('set-cookie', 'a=1; Path=/api/v1/auth');
		header.append('set-cookie', 'b=2');
		const keluar = saringRespons(new Response('{}', { status: 201, headers: header }));

		expect(keluar.status).toBe(201);
		expect(keluar.headers.getSetCookie()).toEqual(['a=1; Path=/api/v1/auth', 'b=2']);
		expect(keluar.headers.get('server')).toBeNull();
	});
});
