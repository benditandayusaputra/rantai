import { bacaPesanPush, tujuanKlikPush } from '$lib/features/pantau/pesanPush';

interface PeristiwaTertunda extends Event {
	waitUntil(janji: Promise<unknown>): void;
}

interface PeristiwaPush extends PeristiwaTertunda {
	data: { text(): string } | null;
}

interface PeristiwaKlikNotifikasi extends PeristiwaTertunda {
	notification: Notification;
}

interface LingkupPekerja {
	registration: ServiceWorkerRegistration;
	clients: { openWindow(url: string): Promise<unknown> };
	addEventListener(jenis: 'push', penangan: (event: PeristiwaPush) => void): void;
	addEventListener(
		jenis: 'notificationclick',
		penangan: (event: PeristiwaKlikNotifikasi) => void
	): void;
}

const pekerja = self as unknown as LingkupPekerja;
const ikon = '/favicon.svg';

pekerja.addEventListener('push', (event) => {
	const pesan = bacaPesanPush(event.data?.text() ?? '');
	if (!pesan) return;
	event.waitUntil(
		pekerja.registration.showNotification(pesan.judul, {
			body: pesan.isi,
			icon: ikon,
			data: { tujuan: tujuanKlikPush(pesan.tautan) }
		})
	);
});

pekerja.addEventListener('notificationclick', (event) => {
	event.notification.close();
	const tujuan = (event.notification.data as { tujuan?: string } | null)?.tujuan ?? '/';
	event.waitUntil(pekerja.clients.openWindow(tujuan));
});
