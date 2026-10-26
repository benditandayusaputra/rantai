import type { KondisiLaporan, Sektor, StatusLaporan } from '$lib/api/tipe';
import { kunciStatusSingkat } from '$lib/features/lapor/status';
import { kunciSektor } from '$lib/features/paket/paket';
import { t, type Bahasa } from '$lib/i18n';
import type { IsiPopup } from './PopupPenanda.svelte';

export const idPenanda = ['sekolah-titik', 'proyek-titik', 'laporan-titik'] as const;
type IdPenanda = (typeof idPenanda)[number];

type Properti = Record<string, unknown>;

function teks(nilai: unknown): string {
	return typeof nilai === 'string' || typeof nilai === 'number' ? String(nilai) : '';
}

function popupProyek(p: Properti, bahasa: Bahasa): IsiPopup {
	const sektor = (teks(p.sektor) in kunciSektor ? teks(p.sektor) : '') as Sektor;
	return {
		keterangan: `${t('map.layer.proyek')} · ${t(kunciSektor[sektor])}`,
		judul: teks(p.nama),
		tautan: {
			href: `/${bahasa}/paket/${encodeURIComponent(teks(p.kode))}`,
			label: t('map.popup.package')
		}
	};
}

function popupLaporan(p: Properti, bahasa: Bahasa): IsiPopup {
	const status = t(kunciStatusSingkat(teks(p.status) as StatusLaporan));
	return {
		keterangan: p.is_simulation ? `${status} · ${t('cert.simulasi')}` : status,
		judul: t(`report.cond.${teks(p.kondisi) as KondisiLaporan}`),
		tautan: {
			href: `/${bahasa}/laporan/${encodeURIComponent(teks(p.id))}`,
			label: t('map.popup.report')
		}
	};
}

function popupSekolah(p: Properti): IsiPopup {
	const status = p.status === 'swasta' ? t('map.popup.swasta') : t('map.popup.negeri');
	return { keterangan: `${teks(p.bentuk)} ${status}`, judul: teks(p.nama) };
}

export function isiPopup(lapisan: IdPenanda, properti: Properti, bahasa: Bahasa): IsiPopup {
	if (lapisan === 'proyek-titik') return popupProyek(properti, bahasa);
	if (lapisan === 'laporan-titik') return popupLaporan(properti, bahasa);
	return popupSekolah(properti);
}

export function isIdPenanda(id: string): id is IdPenanda {
	return (idPenanda as readonly string[]).includes(id);
}
