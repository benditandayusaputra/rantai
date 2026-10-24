import { kirimData, type Fetch } from './klien';
import type { Sesi } from './tipe';

const denganCookie: RequestInit = { credentials: 'include' };

export function kirimKodeMasuk(fetchFn: Fetch, email: string): Promise<null> {
	return kirimData<null>(fetchFn, '/auth/otp/kirim', { email });
}

export function verifikasiKodeMasuk(fetchFn: Fetch, email: string, kode: string): Promise<Sesi> {
	return kirimData<Sesi>(fetchFn, '/auth/otp/verifikasi', { email, kode }, denganCookie);
}

export function segarkanSesi(fetchFn: Fetch): Promise<Sesi> {
	return kirimData<Sesi>(fetchFn, '/auth/refresh', {}, denganCookie);
}

export function akhiriSesi(fetchFn: Fetch): Promise<null> {
	return kirimData<null>(fetchFn, '/auth/keluar', {}, denganCookie);
}
