<script lang="ts">
	import { MapPin, School } from '@lucide/svelte';
	import type { SekolahTerdekat } from '$lib/api/tipe';
	import ProvenanceBadge from '$lib/components/ProvenanceBadge.svelte';
	import { tampilkanAngka } from '$lib/features/kecamatan/kecamatan';
	import { kunciJenisSekolah } from '$lib/features/pendidikan/pendidikan';
	import { formatJarak } from '$lib/format';
	import { bahasaAktif, t } from '$lib/i18n';

	let { sekolah, urutan }: { sekolah: SekolahTerdekat; urutan?: number } = $props();

	const kursi = $derived(sekolah.daya_tampung ? tampilkanAngka(sekolah.daya_tampung) : null);
</script>

<li class="card flex gap-3 p-4">
	<span class="icon-tile num h-10 w-10 self-start font-bold text-secondary-on" aria-hidden="true">
		{#if urutan}{urutan}{:else}<School size={18} strokeWidth={1.75} />{/if}
	</span>
	<span class="flex min-w-0 flex-1 flex-col gap-1">
		<span class="font-semibold">{sekolah.nama}</span>
		<span class="flex flex-wrap items-center gap-x-2 text-small text-muted">
			{t(kunciJenisSekolah(sekolah))}
			<span class="inline-flex items-center gap-1">
				<MapPin size={14} strokeWidth={1.75} aria-hidden="true" />
				{t('check.distance', { jarak: formatJarak(sekolah.jarak_m, bahasaAktif) })}
			</span>
		</span>
		<span class="flex flex-wrap items-center gap-2 text-small">
			{t('pen.school.seats')}:
			{#if sekolah.daya_tampung && kursi}
				<span class="num font-bold">{kursi.nilai}</span>
				<ProvenanceBadge angka={sekolah.daya_tampung} nilai={kursi.nilai} compact />
			{:else}
				<span class="text-muted">{t('pen.school.noSeats')}</span>
			{/if}
		</span>
	</span>
</li>
