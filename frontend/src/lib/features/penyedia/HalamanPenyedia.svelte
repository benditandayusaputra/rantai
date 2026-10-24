<script lang="ts">
	import { Info } from '@lucide/svelte';
	import type { ProfilPenyedia } from '$lib/api/tipe';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import KeyNumberCard from '$lib/components/KeyNumberCard.svelte';
	import PackageCard from '$lib/components/PackageCard.svelte';
	import ProvenanceBadge from '$lib/components/ProvenanceBadge.svelte';
	import TableToggle from '$lib/components/TableToggle.svelte';
	import { tampilkanAngka } from '$lib/features/kecamatan/kecamatan';
	import { t } from '$lib/i18n';
	import { inisialPenyedia, tahunTerbaru, tinggiBatang } from './penyedia';

	let { penyedia }: { penyedia: ProfilPenyedia } = $props();

	let isTabel = $state(false);
	const terbaru = $derived(tahunTerbaru(penyedia.per_tahun));
	const tinggi = $derived(tinggiBatang(penyedia.per_tahun));
	const urutTahun = $derived([...penyedia.per_tahun].sort((a, b) => a.tahun - b.tahun));
</script>

<svelte:head>
	<title>{penyedia.nama} · {t('brand.name')}</title>
</svelte:head>

<article class="wrap flex flex-col gap-8 py-6 lg:py-10">
	<header class="flex flex-col gap-3">
		<p class="t-over">{t('vendor.overline')}</p>
		<div class="flex items-center gap-4">
			<span
				class="inline-flex h-16 w-16 shrink-0 items-center justify-center rounded-md bg-secondary text-h3 font-bold text-secondary-on"
				aria-hidden="true">{inisialPenyedia(penyedia.nama)}</span
			>
			<h1 class="t-h1">{penyedia.nama}</h1>
		</div>
		{#if penyedia.alias.length > 1}
			<p class="text-small text-muted">
				{t('vendor.alias', { nama: penyedia.alias.filter((a) => a !== penyedia.nama).join(', ') })}
			</p>
		{/if}
		<p class="flex max-w-prose gap-2 text-small text-muted">
			<Info size={18} strokeWidth={1.75} class="mt-0.5 shrink-0" aria-hidden="true" />
			{t('vendor.noScore')}
		</p>
	</header>

	{#if terbaru}
		<div class="grid grid-cols-2 gap-3 lg:max-w-prose lg:gap-5">
			<KeyNumberCard
				label={t('vendor.latestCount', { tahun: terbaru.tahun })}
				angka={terbaru.jumlah_paket}
			/>
			<KeyNumberCard
				label={t('vendor.latestValue', { tahun: terbaru.tahun })}
				angka={terbaru.total_nilai}
			/>
		</div>
	{/if}

	<div class="grid gap-8 lg:grid-cols-12">
		<div class="flex flex-col gap-8 lg:col-span-8">
			{#if urutTahun.length}
				<section aria-labelledby="judul-per-tahun" class="flex flex-col gap-3">
					<h2 id="judul-per-tahun" class="t-h2">{t('vendor.perYear')}</h2>
					<p class="text-small text-muted">{t('vendor.perYearNote')}</p>
					<TableToggle bind:pressed={isTabel} controls="isi-per-tahun" />
					<div id="isi-per-tahun">
						{#if isTabel}
							<div class="overflow-x-auto rounded-lg border border-line bg-surface">
								<table class="w-full text-small">
									<caption class="sr-only">{t('vendor.perYear')}</caption>
									<thead class="bg-sunken text-caption uppercase">
										<tr>
											<th scope="col" class="p-3 text-left">{t('vendor.col.year')}</th>
											<th scope="col" class="p-3 text-right">{t('vendor.col.count')}</th>
											<th scope="col" class="p-3 text-right">{t('vendor.col.value')}</th>
										</tr>
									</thead>
									<tbody>
										{#each urutTahun as item (`${item.tahun}-${item.sumber}`)}
											<tr class="border-b border-line">
												<th scope="row" class="p-3 text-left num">{item.tahun}</th>
												<td class="p-3 text-right num">{tampilkanAngka(item.jumlah_paket).nilai}</td
												>
												<td class="p-3 text-right">
													<ProvenanceBadge
														angka={item.total_nilai}
														nilai={tampilkanAngka(item.total_nilai).nilai}
														compact
													/>
												</td>
											</tr>
										{/each}
									</tbody>
								</table>
							</div>
						{:else}
							<ol
								class="flex h-56 items-end gap-2 overflow-x-auto rounded-lg border border-line bg-surface p-4"
								role="img"
								aria-label={t('vendor.perYearNote')}
							>
								{#each urutTahun as item (`${item.tahun}-${item.sumber}`)}
									<li class="flex h-full min-w-11 flex-1 flex-col items-center justify-end gap-1">
										<span
											class="flex w-full items-start justify-center rounded-t-sm bg-primary pt-1 text-caption font-bold text-primary-on num"
											style:height="{Math.max(12, (tinggi.get(item.tahun) ?? 0) * 100)}%"
											>{tampilkanAngka(item.jumlah_paket).nilai}</span
										>
										<span class="text-caption text-muted num">{item.tahun}</span>
									</li>
								{/each}
							</ol>
						{/if}
					</div>
				</section>
			{/if}

			<section aria-labelledby="judul-paket-penyedia" class="flex flex-col gap-3">
				<h2 id="judul-paket-penyedia" class="t-h2">{t('vendor.packages')}</h2>
				{#if penyedia.paket.length}
					<div class="grid gap-3 md:grid-cols-2">
						{#each penyedia.paket as paket (paket.kode)}
							<PackageCard {paket} />
						{/each}
					</div>
				{:else}
					<EmptyState title={t('vendor.noPackages')} />
				{/if}
			</section>

			{#if penyedia.riwayat_opentender.length}
				<section aria-labelledby="judul-riwayat" class="flex flex-col gap-3">
					<h2 id="judul-riwayat" class="t-h2">{t('vendor.history')}</h2>
					<ul class="flex flex-col divide-y divide-line rounded-lg border border-line bg-surface">
						{#each penyedia.riwayat_opentender as r, i (`${r.tahun}-${i}`)}
							<li class="flex flex-col gap-1 p-4">
								<span class="font-semibold">{r.nama_paket}</span>
								<span class="text-small text-muted"
									>{r.tahun} · {r.pembeli}{#if r.nilai}
										· {tampilkanAngka(r.nilai).nilai}{/if}</span
								>
							</li>
						{/each}
					</ul>
				</section>
			{/if}
		</div>

		<aside class="flex flex-col gap-6 lg:col-span-4">
			{#each [{ judul: t('vendor.categories'), isi: penyedia.kategori_utama }, { judul: t('vendor.agencies'), isi: penyedia.instansi_pemberi }] as blok (blok.judul)}
				{#if blok.isi.length}
					<section class="flex flex-col gap-2 rounded-lg border border-line bg-surface p-4">
						<h2 class="t-h3">{blok.judul}</h2>
						<ul class="flex flex-col gap-1 text-small">
							{#each blok.isi as item (item)}<li>{item}</li>{/each}
						</ul>
					</section>
				{/if}
			{/each}
			<section class="flex flex-col gap-2 rounded-lg border border-line bg-surface p-4">
				<h2 class="t-h3">{t('vendor.attribution')}</h2>
				<p class="text-small text-muted">{penyedia.atribusi}</p>
			</section>
		</aside>
	</div>
</article>
