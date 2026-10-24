<script lang="ts">
	import { ambilPaketTerkait, tangkap } from '$lib/api/klien';
	import type { HasilPaketTerkait } from '$lib/api/tipe';
	import { bahasaBawaan, t } from '$lib/i18n';

	let { idAspirasi }: { idAspirasi: number } = $props();

	let hasil = $state<HasilPaketTerkait | null>(null);
	let isGagal = $state(false);
	let isDimuat = false;

	async function muat(event: Event & { currentTarget: HTMLDetailsElement }) {
		if (!event.currentTarget.open || isDimuat) return;
		isDimuat = true;
		const jawaban = await tangkap(ambilPaketTerkait(fetch, idAspirasi));
		if (jawaban.ok) hasil = jawaban.data;
		else isGagal = true;
	}
</script>

<details class="rounded-md bg-sunken" ontoggle={muat}>
	<summary class="flex min-h-touch cursor-pointer items-center px-3 text-small font-semibold"
		>{t('related.toggle')}</summary
	>
	<div class="flex flex-col gap-2 px-3 pb-3 text-small" aria-live="polite">
		{#if isGagal}
			<p class="text-muted">{t('related.failed')}</p>
		{:else if !hasil}
			<p class="text-muted">{t('state.loading')}</p>
		{:else}
			<p class="text-muted">{t('related.lead', { tahun: hasil.tahun_anggaran })}</p>
			{#if hasil.paket.length === 0}
				<p class="font-semibold">{t('related.none')}</p>
			{:else}
				<ul class="flex flex-col gap-2">
					{#each hasil.paket as paket (paket.kode)}
						<li class="flex flex-col gap-1 rounded-md bg-surface p-3">
							<a href="/{bahasaBawaan}/paket/{paket.kode}" class="font-semibold underline"
								>{paket.nama}</a
							>
							<span
								class="self-start rounded-xs border-ctrl border-dashed border-line-strong px-2 text-caption font-bold"
								>{t(`related.match.${paket.dasar_cocok}`)}</span
							>
						</li>
					{/each}
				</ul>
			{/if}
		{/if}
	</div>
</details>
