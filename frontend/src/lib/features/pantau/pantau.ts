import {
	CalendarClock,
	Camera,
	ChartColumn,
	FileText,
	Hourglass,
	MessageSquareReply
} from '@lucide/svelte';
import type { JenisKejadian, Notifikasi } from '$lib/api/tipe';
import { formatTanggal } from '$lib/format';
import { bahasaBawaan, t, type KunciPesan } from '$lib/i18n';

const zonaWaktu = 'Asia/Jakarta';
const milidetikSehari = 86_400_000;

export const tampilanKejadian: Record<
	JenisKejadian,
	{ label: KunciPesan; ikon: typeof Camera; kelas: string }
> = {
	paket_baru: {
		label: 'watch.ev.newPackage',
		ikon: FileText,
		kelas: 'bg-secondary text-secondary-on'
	},
	melewati_jadwal: {
		label: 'watch.ev.pastSchedule',
		ikon: CalendarClock,
		kelas: 'bg-attn-bg text-attn'
	},
	pemeliharaan_berakhir: {
		label: 'watch.ev.maintEnding',
		ikon: Hourglass,
		kelas: 'bg-attn-bg text-attn'
	},
	laporan_terverifikasi: {
		label: 'watch.ev.reportVerified',
		ikon: Camera,
		kelas: 'bg-success-bg text-success'
	},
	hak_jawab: {
		label: 'watch.ev.newReply',
		ikon: MessageSquareReply,
		kelas: 'bg-info-bg text-info'
	},
	indikator_diperbarui: {
		label: 'watch.ev.indicator',
		ikon: ChartColumn,
		kelas: 'bg-info-bg text-info'
	}
};

function kunciHari(d: Date): string {
	return new Intl.DateTimeFormat('en-CA', { timeZone: zonaWaktu }).format(d);
}

export function labelHari(d: Date, sekarang: Date): string {
	const hari = kunciHari(d);
	if (hari === kunciHari(sekarang)) return t('watch.today');
	if (hari === kunciHari(new Date(sekarang.getTime() - milidetikSehari)))
		return t('watch.yesterday');
	return formatTanggal(d, bahasaBawaan, 'long');
}

export function kelompokkanPerHari(
	daftar: Notifikasi[],
	sekarang: Date
): { hari: string; item: Notifikasi[] }[] {
	const kelompok = new Map<string, Notifikasi[]>();
	for (const n of daftar) {
		const label = labelHari(new Date(n.terjadi_pada), sekarang);
		kelompok.set(label, [...(kelompok.get(label) ?? []), n]);
	}
	return [...kelompok].map(([hari, item]) => ({ hari, item }));
}
