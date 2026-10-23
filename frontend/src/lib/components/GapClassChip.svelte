<script lang="ts">
	import type { KelasRasio } from '$lib/api/tipe';
	import { kunciKelas } from '$lib/features/kecamatan/kecamatan';
	import { t } from '$lib/i18n';

	let { kelas }: { kelas: KelasRasio } = $props();

	const warna: Record<KelasRasio, string> = {
		sangat_kurang: 'bg-gap-1 text-gap-1-on',
		kurang: 'bg-gap-2 text-gap-2-on',
		hampir_cukup: 'bg-gap-3 text-gap-3-on',
		cukup: 'bg-gap-4 text-gap-4-on',
		tidak_ada_data: 'border-ctrl border-dashed border-line-strong text-muted'
	};
</script>

<span
	class="inline-flex h-7 items-center rounded-pill px-3 text-caption font-bold {warna[kelas]}"
	data-kelas={kelas}>{t(kunciKelas[kelas])}</span
>
