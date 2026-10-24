<script lang="ts">
	import { ambilIsuAnggota, tangkap } from '$lib/api/klien';
	import type { IsuDisuarakan } from '$lib/api/tipe';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import { t } from '$lib/i18n';
	import KartuUmpan from './KartuUmpan.svelte';

	let { idAnggota }: { idAnggota: string } = $props();

	let daftar = $state<IsuDisuarakan[] | null>(null);

	$effect(() => {
		void tangkap(ambilIsuAnggota(fetch, idAnggota)).then((hasil) => {
			daftar = hasil.ok ? hasil.data : [];
		});
	});
</script>

<section aria-labelledby="judul-isu" class="flex flex-col gap-4">
	<h2 id="judul-isu" class="t-h2">{t('issue.title')}</h2>
	<p class="text-small text-muted">{t('issue.lead')}</p>
	{#if daftar === null}
		<p class="text-muted" aria-live="polite">{t('state.loading')}</p>
	{:else if daftar.length === 0}
		<EmptyState title={t('issue.empty')} />
	{:else}
		<ul class="flex flex-col gap-3">
			{#each daftar as isu (isu.tag)}
				<li>
					<details class="rounded-lg border border-line bg-surface">
						<summary class="flex min-h-touch cursor-pointer flex-col gap-1 p-4">
							<span class="font-bold capitalize">{isu.tag}</span>
							<span class="text-small text-muted">
								{t('issue.count', {
									n: isu.jumlah_sumber,
									resmi: isu.jumlah_rilis_resmi,
									media: isu.jumlah_media
								})}
							</span>
						</summary>
						<ul class="flex flex-col gap-2 px-4 pb-4">
							{#each isu.item as item (item.id)}<KartuUmpan {item} />{/each}
						</ul>
					</details>
				</li>
			{/each}
		</ul>
	{/if}
</section>
