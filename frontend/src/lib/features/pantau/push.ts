import { tangkap, type Fetch } from '$lib/api/klien';
import { aktifkanPush, matikanPush } from '$lib/api/webpush';
import { kunciBase64KeBita } from './pesanPush';

export type StatusPush = 'tidak_didukung' | 'mati' | 'aktif' | 'ditolak' | 'gagal';

export function isPushDidukung(): boolean {
	return 'serviceWorker' in navigator && 'PushManager' in window && 'Notification' in window;
}

export async function statusPushAwal(): Promise<StatusPush> {
	if (!isPushDidukung()) return 'tidak_didukung';
	if (Notification.permission === 'denied') return 'ditolak';
	const pendaftaran = await navigator.serviceWorker.getRegistration();
	return (await pendaftaran?.pushManager.getSubscription()) ? 'aktif' : 'mati';
}

export async function nyalakanPush(
	fetchFn: Fetch,
	token: string,
	kunciPublik: string
): Promise<StatusPush> {
	if ((await Notification.requestPermission()) !== 'granted') return 'ditolak';
	const pendaftaran = await navigator.serviceWorker.ready;
	const langganan = await pendaftaran.pushManager.subscribe({
		userVisibleOnly: true,
		applicationServerKey: kunciBase64KeBita(kunciPublik)
	});
	const hasil = await tangkap(aktifkanPush(fetchFn, token, langganan.toJSON()));
	return hasil.ok ? 'aktif' : 'gagal';
}

export async function padamkanPush(fetchFn: Fetch, token: string): Promise<StatusPush> {
	const pendaftaran = await navigator.serviceWorker.getRegistration();
	const langganan = await pendaftaran?.pushManager.getSubscription();
	if (langganan) {
		await tangkap(matikanPush(fetchFn, token, langganan.endpoint));
		await langganan.unsubscribe();
	}
	return 'mati';
}
