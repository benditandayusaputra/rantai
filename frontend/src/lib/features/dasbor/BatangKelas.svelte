<script lang="ts">
	import type { KelasRasio, RekapKursi } from '$lib/api/tipe';
	import { kunciKelas } from '$lib/features/kecamatan/kecamatan';
	import { t } from '$lib/i18n';
	import { persenDari, totalKecamatan, urutanKelas } from './dasbor';

	interface Props {
		daftar: RekapKursi[];
		namaDari: (r: RekapKursi) => string;
	}

	let { daftar, namaDari }: Props = $props();

	const warna: Record<KelasRasio, string> = {
		sangat_kurang: 'bg-gap-1',
		kurang: 'bg-gap-2',
		hampir_cukup: 'bg-gap-3',
		cukup: 'bg-gap-4',
		tidak_ada_data: 'bg-sunken'
	};
</script>

<div class="flex flex-col gap-4">
	<ul class="flex flex-wrap gap-x-4 gap-y-2 text-caption text-muted" aria-hidden="true">
		{#each urutanKelas as kelas (kelas)}
			<li class="inline-flex items-center gap-2">
				<span class="h-3 w-5 rounded-xs {warna[kelas]}"></span>{t(kunciKelas[kelas])}
			</li>
		{/each}
	</ul>
	<ol class="flex flex-col gap-4">
		{#each daftar as r (r.kunci)}
			{@const total = totalKecamatan(r)}
			{@const segmen = urutanKelas
				.map((kelas) => ({ kelas, n: r.per_kelas[kelas] ?? 0 }))
				.filter((s) => s.n > 0)}
			<li class="flex flex-col gap-1.5">
				<div class="flex items-baseline justify-between gap-3 text-small">
					<span class="font-semibold">{namaDari(r)}</span>
					<span class="num text-muted">{t('dash.svc.kecamatan', { n: total })}</span>
				</div>
				<div class="flex h-5 gap-0.5 overflow-hidden rounded-sm" aria-hidden="true">
					{#each segmen as s (s.kelas)}
						<span class={warna[s.kelas]} style:width="{persenDari(s.n, total)}%"></span>
					{/each}
				</div>
				<div class="flex gap-0.5">
					{#each segmen as s (s.kelas)}
						<span
							class="num text-center text-caption font-bold"
							style:width="{persenDari(s.n, total)}%"
							title="{t(kunciKelas[s.kelas])}: {s.n}"
						>
							<span class="sr-only">{t(kunciKelas[s.kelas])}:</span>{s.n}
						</span>
					{/each}
				</div>
			</li>
		{/each}
	</ol>
</div>
