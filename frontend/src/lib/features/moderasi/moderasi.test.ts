import { describe, expect, it } from 'vitest';
import { bacaPintasan, geser, sisihkan } from './moderasi';

function tekan(key: string, target: EventTarget = document.body): KeyboardEvent {
	const event = new KeyboardEvent('keydown', { key });
	Object.defineProperty(event, 'target', { value: target });
	return event;
}

describe('moderasi', () => {
	it('memetakan J, K, V, R kecuali saat mengetik', () => {
		expect([tekan('j'), tekan('K'), tekan('v'), tekan('r')].map(bacaPintasan)).toEqual([
			'berikut',
			'sebelum',
			'verifikasi',
			'tolak'
		]);
		expect(bacaPintasan(tekan('v', document.createElement('textarea')))).toBeNull();
		expect(bacaPintasan(tekan('x'))).toBeNull();
	});

	it('menggeser pilihan dalam batas antrean', () => {
		expect(geser(0, -1, 3)).toBe(0);
		expect(geser(2, 1, 3)).toBe(2);
		expect(geser(1, 1, 3)).toBe(2);
		expect(geser(0, 1, 0)).toBe(0);
	});

	it('menyisihkan item yang sudah ditangani', () => {
		expect(sisihkan([{ id: 1 }, { id: 2 }], 1)).toEqual([{ id: 2 }]);
	});
});
