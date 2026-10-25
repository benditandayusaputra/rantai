import {
	CircleArrowUp,
	CircleCheck,
	CircleX,
	Clock,
	Forward,
	Reply,
	ShieldCheck
} from '@lucide/svelte';
import type { StatusLaporan } from '$lib/api/tipe';
import type { KunciPesan } from '$lib/i18n';

interface TampilanStatus {
	label: KunciPesan;
	kelas: string;
	ikon: typeof Clock;
}

export const tampilanStatus: Record<StatusLaporan, TampilanStatus> = {
	DIKIRIM: {
		label: 'report.dikirim',
		kelas: 'border-st-sent-bd bg-st-sent-bg text-st-sent-fg',
		ikon: CircleArrowUp
	},
	MENUNGGU_SAKSI: {
		label: 'report.menunggu-saksi',
		kelas: 'border-dashed border-st-wait-bd bg-st-wait-bg text-st-wait-fg',
		ikon: Clock
	},
	TERVERIFIKASI: {
		label: 'report.terverifikasi',
		kelas: 'border-st-ver-bd bg-st-ver-bg text-st-ver-fg',
		ikon: ShieldCheck
	},
	DITERUSKAN: {
		label: 'report.diteruskan',
		kelas: 'border-st-fwd-bd bg-st-fwd-bg text-st-fwd-fg',
		ikon: Forward
	},
	DITANGGAPI: {
		label: 'report.ditanggapi',
		kelas: 'border-st-resp-bd bg-st-resp-bg text-st-resp-fg',
		ikon: Reply
	},
	SELESAI: {
		label: 'report.selesai',
		kelas: 'border-st-done-bd bg-st-done-bg text-st-done-fg',
		ikon: CircleCheck
	},
	DITOLAK: {
		label: 'report.ditolak',
		kelas: 'border-st-rej-bd bg-st-rej-bg text-st-rej-fg',
		ikon: CircleX
	}
};
