<script lang="ts">
	import {
		ChevronDown,
		History,
		Info,
		LoaderCircle,
		Landmark,
		Layers,
		Newspaper
	} from '@lucide/svelte';
	import { ambilUmpanAnggota, tangkap } from '$lib/api/klien';
	import type { ItemUmpan, JenisSumberUmpan } from '$lib/api/tipe';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import ErrorState from '$lib/components/ErrorState.svelte';
	import { t } from '$lib/i18n';
	import KartuUmpan from './KartuUmpan.svelte';

	let { idAnggota }: { idAnggota: string } = $props();

	const pilihan: (JenisSumberUmpan | 'semua')[] = ['semua', 'resmi_dprd', 'media'];
	const ikonPilihan = { semua: Layers, resmi_dprd: Landmark, media: Newspaper } as const;
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
	<div class="flex items-center gap-3">
		<span class="icon-tile"><History size={20} strokeWidth={1.75} aria-hidden="true" /></span>
		<h2 id="judul-jejak" class="t-h2">{t('member.feed')}</h2>
	</div>
	<p class="flex gap-2 rounded-md bg-info-bg p-3 text-small">
		<Info size={18} strokeWidth={1.75} class="mt-0.5 shrink-0 text-info" aria-hidden="true" />
		{t('member.bias')}
	</p>
	<div role="radiogroup" aria-label={t('member.feed')} class="flex flex-wrap gap-2">
		{#each pilihan as p (p)}
			{@const Ikon = ikonPilihan[p]}
			<label
				class="chip cursor-pointer has-[:checked]:border-primary has-[:checked]:bg-secondary has-[:checked]:text-secondary-on has-[:focus-visible]:shadow-ring"
			>
				<Ikon size={16} strokeWidth={1.75} aria-hidden="true" />
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
	{#if status === 'memuat'}
		<p class="flex items-center gap-2 text-muted" aria-live="polite">
			<LoaderCircle size={18} strokeWidth={1.75} aria-hidden="true" class="animate-spin" />
			{t('state.loading')}
		</p>
	{/if}
	{#if status === 'siap' && halaman < totalHalaman}
		<button type="button" class="btn btn-outline self-start" onclick={() => muat(halaman + 1)}>
			<ChevronDown size={18} strokeWidth={1.75} aria-hidden="true" />
			{t('action.loadMoreShort')}
		</button>
	{/if}
</section>
