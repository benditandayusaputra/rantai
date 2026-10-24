import { ambilData, kirimBerwenang, type Fetch } from './klien';

export function ambilKunciPublikPush(fetchFn: Fetch): Promise<{ kunci_publik: string }> {
	return ambilData<{ kunci_publik: string }>(fetchFn, '/webpush/kunci-publik');
}

export function aktifkanPush(
	fetchFn: Fetch,
	token: string,
	langganan: PushSubscriptionJSON
): Promise<{ jumlah_kecamatan: number }> {
	return kirimBerwenang(fetchFn, '/webpush/langganan', token, langganan);
}

export function matikanPush(fetchFn: Fetch, token: string, endpoint: string): Promise<null> {
	return kirimBerwenang<null>(fetchFn, '/webpush/berhenti', token, { endpoint });
}
