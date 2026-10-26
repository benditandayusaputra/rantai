<script lang="ts">
	import { Camera, School } from '@lucide/svelte';
	import type { Sektor } from '$lib/api/tipe';
	import { ikonSektor, kunciSektor } from '$lib/features/paket/paket';
	import { t } from '$lib/i18n';
	import type { StatusLapisan } from './lapisan';
	import { daftarSektor } from './penanda';

	let { lapisan }: { lapisan: StatusLapisan } = $props();

	const kelasSektor: Record<Sektor, string> = {
		pendidikan: 'bg-sector-edu text-sector-edu-on',
		air: 'bg-sector-air text-sector-air-on',
		jalan: 'bg-sector-jalan text-sector-jalan-on',
		'': 'bg-sector-lain text-sector-lain-on'
	};

	const butirSekolah = [
		{ kunci: 'sekolah', ikon: School, kelas: 'bg-fg text-surface', label: t('map.marker.school') }
	];
	const butirProyek = daftarSektor.map((sektor) => ({
		kunci: `proyek-${sektor}`,
		ikon: ikonSektor[sektor],
		kelas: kelasSektor[sektor],
		label: `${t('map.layer.proyek')} · ${t(kunciSektor[sektor])}`
	}));
	const butirLaporan = [
		{ kunci: 'tuntas', kelas: 'bg-success text-success-on', label: t('map.marker.reportDone') },
		{ kunci: 'proses', kelas: 'bg-info text-surface', label: t('map.marker.reportOpen') },
		{ kunci: 'simulasi', kelas: 'bg-cert-sim text-surface', label: t('map.marker.reportSim') }
	].map((butir) => ({ ...butir, ikon: Camera }));

	const butir = $derived([
		...(lapisan.sekolah ? butirSekolah : []),
		...(lapisan.proyek ? butirProyek : []),
		...(lapisan.laporan ? butirLaporan : [])
	]);
</script>

{#if butir.length}
	<div>
		<p class="text-caption font-bold">{t('map.markers')}</p>
		<ul class="mt-2 flex flex-col gap-1">
			{#each butir as item (item.kunci)}
				<li class="flex items-center gap-2 text-caption">
					<span
						class="inline-flex h-5 w-5 shrink-0 items-center justify-center rounded-pill border-2 border-surface shadow-1 {item.kelas}"
						aria-hidden="true"
					>
						<item.ikon size={12} strokeWidth={2.25} />
					</span>
					{item.label}
				</li>
			{/each}
		</ul>
	</div>
{/if}
