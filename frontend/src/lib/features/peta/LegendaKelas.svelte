<script lang="ts">
	import type { KelasRasio } from '$lib/api/tipe';
	import { kunciKelas } from '$lib/features/kecamatan/kecamatan';
	import { t } from '$lib/i18n';

	const kotak: Record<KelasRasio, string> = {
		sangat_kurang: 'bg-gap-1',
		kurang: 'bg-gap-2',
		hampir_cukup: 'bg-gap-3',
		cukup: 'bg-gap-4',
		tidak_ada_data: 'border-ctrl border-dashed border-line-strong bg-sunken'
	};
	const urutan = Object.keys(kotak) as KelasRasio[];
</script>

<div>
	<p class="text-caption font-bold">{t('map.legend')}</p>
	<ul class="mt-2 flex flex-col gap-1">
		{#each urutan as kelas (kelas)}
			<li class="flex items-center gap-2 text-caption">
				<span class="h-4 w-6 rounded-xs {kotak[kelas]}" aria-hidden="true"></span>
				{t(kunciKelas[kelas])}
			</li>
		{/each}
	</ul>
</div>
