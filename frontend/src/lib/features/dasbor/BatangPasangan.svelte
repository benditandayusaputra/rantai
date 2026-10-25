<script lang="ts">
	import type { RekapKursi } from '$lib/api/tipe';
	import { formatAngka } from '$lib/format';
	import { bahasaAktif, t } from '$lib/i18n';
	import { persenDari } from './dasbor';

	interface Props {
		daftar: RekapKursi[];
		namaDari: (r: RekapKursi) => string;
	}

	let { daftar, namaDari }: Props = $props();

	const maksimum = $derived(Math.max(0, ...daftar.flatMap((r) => [r.kebutuhan, r.kapasitas])));
	const ukuran = [
		{ kunci: 'kebutuhan', label: t('dash.svc.graduates'), warna: 'bg-fg' },
		{ kunci: 'kapasitas', label: t('dash.svc.seats'), warna: 'bg-primary' }
	] as const;
</script>

<div class="flex flex-col gap-4">
	<ul class="flex flex-wrap gap-4 text-caption text-muted" aria-hidden="true">
		{#each ukuran as u (u.kunci)}
			<li class="inline-flex items-center gap-2">
				<span class="h-3 w-3 rounded-xs {u.warna}"></span>{u.label}
			</li>
		{/each}
	</ul>
	<ol class="flex flex-col gap-5">
		{#each daftar as r (r.kunci)}
			<li class="flex flex-col gap-1.5">
				<span class="text-small font-semibold">{namaDari(r)}</span>
				{#each ukuran as u (u.kunci)}
					<div class="flex items-center gap-3">
						<div class="h-3 flex-1 rounded-sm bg-sunken" aria-hidden="true">
							<div
								class="h-full rounded-sm {u.warna}"
								style:width="{persenDari(r[u.kunci], maksimum)}%"
							></div>
						</div>
						<span class="num w-20 shrink-0 text-right text-small font-bold">
							<span class="sr-only">{u.label}:</span>
							{formatAngka(r[u.kunci], bahasaAktif)}
						</span>
					</div>
				{/each}
			</li>
		{/each}
	</ol>
</div>
