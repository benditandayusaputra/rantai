import { akhiriSesi, segarkanSesi } from '$lib/api/auth';
import { tangkap, type Fetch } from '$lib/api/klien';
import type { Pengguna, Sesi } from '$lib/api/tipe';

const keadaan = $state<{ sesi: Sesi | null; isSudahDipulihkan: boolean }>({
	sesi: null,
	isSudahDipulihkan: false
});

export function sesiAktif(): Sesi | null {
	return keadaan.sesi;
}

export function simpanSesi(sesi: Sesi): void {
	keadaan.sesi = sesi;
}

export async function pulihkanSesi(fetchFn: Fetch): Promise<void> {
	if (keadaan.isSudahDipulihkan || keadaan.sesi) return;
	keadaan.isSudahDipulihkan = true;
	const hasil = await tangkap(segarkanSesi(fetchFn));
	if (hasil.ok) keadaan.sesi = hasil.data;
}

export async function keluar(fetchFn: Fetch): Promise<void> {
	await tangkap(akhiriSesi(fetchFn));
	keadaan.sesi = null;
}

export function inisial(pengguna: Pengguna): string {
	const kata = pengguna.nama_tampilan.trim().split(/\s+/).filter(Boolean);
	const sumber = kata.length > 0 ? kata.slice(0, 2).map((k) => k[0]) : [pengguna.email[0]];
	return sumber.join('').toUpperCase();
}
