import type { Bahasa } from '$lib/i18n';

const tagBahasa: Record<Bahasa, string> = { id: 'id-ID' };
const zonaWaktu = 'Asia/Jakarta';

export function formatAngka(n: number, bahasa: Bahasa, desimal = 0): string {
	return new Intl.NumberFormat(tagBahasa[bahasa], {
		minimumFractionDigits: desimal,
		maximumFractionDigits: desimal
	}).format(n);
}

export function formatTanggal(d: Date, bahasa: Bahasa, gaya: 'short' | 'long' = 'short'): string {
	return new Intl.DateTimeFormat(tagBahasa[bahasa], {
		day: 'numeric',
		month: gaya,
		year: 'numeric',
		timeZone: zonaWaktu
	})
		.format(d)
		.replace(/\.$/, '');
}

export function formatJam(d: Date): string {
	return new Intl.DateTimeFormat('en-GB', {
		hour: '2-digit',
		minute: '2-digit',
		hour12: false,
		timeZone: zonaWaktu
	})
		.format(d)
		.replace(':', '.');
}

const meterPerKilometer = 1000;

export function formatJarak(meter: number, bahasa: Bahasa): string {
	if (meter < meterPerKilometer) return `${formatAngka(Math.round(meter), bahasa)} m`;
	return `${formatAngka(meter / meterPerKilometer, bahasa, 1)} km`;
}
