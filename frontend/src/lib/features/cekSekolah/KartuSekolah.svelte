<script lang="ts">
	import { School } from '@lucide/svelte';
	import type { SekolahTerdekat } from '$lib/api/tipe';
	import ProvenanceBadge from '$lib/components/ProvenanceBadge.svelte';
	import { tampilkanAngka } from '$lib/features/kecamatan/kecamatan';
	import { kunciJenisSekolah } from '$lib/features/pendidikan/pendidikan';
	import { formatJarak } from '$lib/format';
	import { bahasaBawaan, t } from '$lib/i18n';

	let { sekolah, urutan }: { sekolah: SekolahTerdekat; urutan?: number } = $props();

	const kursi = $derived(sekolah.daya_tampung ? tampilkanAngka(sekolah.daya_tampung) : null);
</script>

<li class="flex gap-3 rounded-lg border border-line bg-surface p-4">
	<span
		class="inline-flex h-9 w-9 shrink-0 self-start items-center justify-center rounded-md bg-secondary font-bold text-secondary-on num"
		aria-hidden="true"
	>
		{#if urutan}{urutan}{:else}<School size={18} strokeWidth={1.75} />{/if}
	</span>
	<span class="flex flex-1 flex-col gap-1">
		<span class="font-semibold">{sekolah.nama}</span>
		<span class="text-small text-muted"
			>{t(kunciJenisSekolah(sekolah))} · {t('check.distance', {
				jarak: formatJarak(sekolah.jarak_m, bahasaBawaan)
			})}</span
		>
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
