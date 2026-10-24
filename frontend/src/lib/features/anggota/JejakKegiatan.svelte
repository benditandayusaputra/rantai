<script lang="ts">
	import { Info } from '@lucide/svelte';
	import { ambilUmpanAnggota, tangkap } from '$lib/api/klien';
	import type { ItemUmpan, JenisSumberUmpan } from '$lib/api/tipe';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import ErrorState from '$lib/components/ErrorState.svelte';
	import { t } from '$lib/i18n';
	import KartuUmpan from './KartuUmpan.svelte';

	let { idAnggota }: { idAnggota: string } = $props();

	const pilihan: (JenisSumberUmpan | 'semua')[] = ['semua', 'resmi_dprd', 'media'];
	let jenis = $state<JenisSumberUmpan | 'semua'>('semua');
	let item = $state<ItemUmpan[]>([]);
	let halaman = $state(1);
	let totalHalaman = $state(1);
	let status = $state<'memuat' | 'siap' | 'gagal'>('memuat');

	async function muat(halamanBaru: number) {
		status = 'memuat';
		const hasil = await tangkap(
			ambilUmpanAnggota(fetch, idAnggota, {
				jenis: jenis === 'semua' ? undefined : jenis,
				halaman: halamanBaru
			})
		);
		if (!hasil.ok) {
			status = 'gagal';
			return;
		}
		item = halamanBaru === 1 ? hasil.data.data : [...item, ...hasil.data.data];
		halaman = halamanBaru;
		totalHalaman = hasil.data.meta?.total_pages ?? 1;
		status = 'siap';
	}

	$effect(() => {
		void jenis;
		void muat(1);
	});
</script>

<section aria-labelledby="judul-jejak" class="flex flex-col gap-4">
	<h2 id="judul-jejak" class="t-h2">{t('member.feed')}</h2>
	<p class="flex gap-2 rounded-md bg-info-bg p-3 text-small">
		<Info size={18} strokeWidth={1.75} class="mt-0.5 shrink-0 text-info" aria-hidden="true" />
		{t('member.bias')}
	</p>
	<div role="radiogroup" aria-label={t('member.feed')} class="flex flex-wrap gap-2">
		{#each pilihan as p (p)}
			<label
				class="inline-flex min-h-touch cursor-pointer items-center rounded-pill border-ctrl border-line-strong px-4 text-small font-semibold has-[:checked]:border-secondary-on has-[:checked]:bg-secondary has-[:focus-visible]:shadow-ring"
			>
				<input type="radio" name="jenis-umpan" value={p} bind:group={jenis} class="sr-only" />
				{t(`member.filter.${p === 'semua' ? 'all' : p}`)}
			</label>
		{/each}
	</div>
	{#if status === 'gagal'}
		<ErrorState title={t('member.feedFailed')} />
	{:else if status === 'siap' && item.length === 0}
		<EmptyState title={t('member.feedEmpty')} />
	{/if}
	<ul class="flex flex-col gap-3">
		{#each item as u (u.id)}<KartuUmpan item={u} />{/each}
	</ul>
	{#if status === 'memuat'}<p class="text-muted" aria-live="polite">{t('state.loading')}</p>{/if}
	{#if status === 'siap' && halaman < totalHalaman}
		<button
			type="button"
			class="inline-flex min-h-touch items-center self-start rounded-md border-ctrl border-line-strong px-4 font-semibold hover:bg-sunken"
			onclick={() => muat(halaman + 1)}>{t('action.loadMoreShort')}</button
		>
	{/if}
</section>
