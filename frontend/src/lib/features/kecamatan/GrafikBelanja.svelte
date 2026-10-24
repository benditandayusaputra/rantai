<script lang="ts">
	import { Info } from '@lucide/svelte';
	import { ambilBelanjaTema, tangkap } from '$lib/api/klien';
	import type { BelanjaTema } from '$lib/api/tipe';
	import ErrorState from '$lib/components/ErrorState.svelte';
	import ProvenanceBadge from '$lib/components/ProvenanceBadge.svelte';
	import TableToggle from '$lib/components/TableToggle.svelte';
	import { formatAngka, formatRupiah } from '$lib/format';
	import { bahasaBawaan, t } from '$lib/i18n';
	import {
		batasSumbu,
		jumlahTanpaRasio,
		titikSebar,
		ukuranUntuk,
		urutBaris,
		type TitikSebar
	} from './belanja';
	import { tampilkanAngka } from './kecamatan';

	let { kodeKecamatan, kodeTema }: { kodeKecamatan: string; kodeTema: string } = $props();

	const id = $props.id();
	const jumlahGaris = 4;
	let data = $state<BelanjaTema | null>(null);
	let isGagal = $state(false);
	let isTabel = $state(false);
	let disorot = $state<TitikSebar | null>(null);
	let lebarKontainer = $state(0);

	$effect(() => {
		void tangkap(ambilBelanjaTema(fetch, kodeTema)).then((hasil) => {
			if (hasil.ok) data = hasil.data;
			else isGagal = true;
		});
	});

	const ukuran = $derived(ukuranUntuk(lebarKontainer));
	const titik = $derived(data ? titikSebar(data, kodeKecamatan, ukuran) : []);
	const sumbu = $derived(data ? batasSumbu(data) : { maksX: 1, maksY: 1 });
	const garis = $derived(Array.from({ length: jumlahGaris + 1 }, (_, i) => i / jumlahGaris));
	const aktif = $derived(titik.find((p) => p.isAktif) ?? null);
	const kotaAktif = $derived(
		data?.luar_kecamatan.find(
			(l) =>
				l.wilayah.kode ===
				data?.kecamatan.find((k) => k.kecamatan.kode === kodeKecamatan)?.kode_kota
		) ?? null
	);
	const { lebar, tinggi, kiri, kanan, atas, bawah } = $derived(ukuran);
	const lebarPlot = $derived(lebar - kiri - kanan);
	const tinggiPlot = $derived(tinggi - atas - bawah);
</script>

<section
	aria-labelledby="{id}-judul"
	class="flex flex-col gap-4 rounded-lg border border-line bg-surface p-4 md:p-5"
>
	<div class="flex flex-wrap items-start justify-between gap-3">
		<h3 id="{id}-judul" class="t-h3">{t('spend.title')}</h3>
		{#if data}<TableToggle bind:pressed={isTabel} controls={id} />{/if}
	</div>
	{#if isGagal}
		<ErrorState title={t('spend.failed')} />
	{:else if !data}
		<p class="text-muted" aria-live="polite">{t('state.loading')}</p>
	{:else}
		<p class="text-small text-muted">{t('spend.lead', { tahun: data.tahun })}</p>
		<div {id}>
			{#if isTabel}
				<div class="max-h-96 overflow-auto rounded-md border border-line">
					<table class="w-full text-small">
						<caption class="sr-only">{t('spend.title')}</caption>
						<thead class="sticky top-0 bg-sunken text-caption uppercase">
							<tr>
								<th scope="col" class="p-3 text-left">{t('spend.col.area')}</th>
								<th scope="col" class="p-3 text-right">{t('spend.col.ratio')}</th>
								<th scope="col" class="p-3 text-right">{t('spend.col.count')}</th>
								<th scope="col" class="p-3 text-right">{t('spend.col.value')}</th>
							</tr>
						</thead>
						<tbody>
							{#each urutBaris(data, kodeKecamatan) as baris (baris.kecamatan.kode)}
								<tr
									class="border-b border-line {baris.kecamatan.kode === kodeKecamatan
										? 'bg-secondary font-bold'
										: ''}"
								>
									<th scope="row" class="p-3 text-left font-semibold">{baris.kecamatan.nama}</th>
									<td class="num p-3 text-right">
										{#if baris.indikator?.rasio}
											<ProvenanceBadge
												angka={baris.indikator.rasio}
												nilai={tampilkanAngka(baris.indikator.rasio).nilai}
												compact
											/>
										{:else}-{/if}
									</td>
									<td class="num p-3 text-right"
										>{baris.belanja ? tampilkanAngka(baris.belanja.jumlah_paket).nilai : '-'}</td
									>
									<td class="num p-3 text-right">
										{#if baris.belanja}
											<ProvenanceBadge
												angka={baris.belanja.total_hps}
												nilai={tampilkanAngka(baris.belanja.total_hps).nilai}
												compact
											/>
										{:else}-{/if}
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			{:else}
				<figure class="relative" bind:clientWidth={lebarKontainer}>
					<svg
						width={lebar}
						height={tinggi}
						class="block max-w-full"
						role="img"
						aria-label="{t('spend.title')}. {aktif
							? `${aktif.nama}: ${formatAngka(aktif.rasio, bahasaBawaan, 2)}, ${formatRupiah(aktif.nilai, bahasaBawaan)}`
							: ''}"
					>
						{#each garis as g (g)}
							<line
								x1={kiri}
								x2={lebar - kanan}
								y1={atas + tinggiPlot * (1 - g)}
								y2={atas + tinggiPlot * (1 - g)}
								class="stroke-line"
								stroke-width="1"
							/>
							<text
								x={kiri - 8}
								y={atas + tinggiPlot * (1 - g)}
								text-anchor="end"
								dominant-baseline="middle"
								class="fill-muted text-caption num"
							>
								{formatRupiah(sumbu.maksY * g, bahasaBawaan)}
							</text>
							<text
								x={kiri + lebarPlot * g}
								y={tinggi - bawah + 18}
								text-anchor="middle"
								class="fill-muted text-caption num"
							>
								{formatAngka(sumbu.maksX * g, bahasaBawaan, 1)}
							</text>
						{/each}
						<text
							x={kiri + lebarPlot / 2}
							y={tinggi - 6}
							text-anchor="middle"
							class="fill-muted text-caption">{t('spend.xAxis')}</text
						>
						{#each titik as p (p.kode)}
							<circle
								role="presentation"
								cx={p.x}
								cy={p.y}
								r={p.isAktif ? 8 : 5}
								class="{p.isAktif ? 'fill-primary' : 'fill-line-strong opacity-70'} stroke-surface"
								stroke-width="2"
								onpointerenter={() => (disorot = p)}
								onpointerleave={() => (disorot = null)}
							/>
							<circle
								role="presentation"
								cx={p.x}
								cy={p.y}
								r="12"
								fill="transparent"
								onpointerenter={() => (disorot = p)}
								onpointerleave={() => (disorot = null)}
							/>
						{/each}
						{#if aktif && disorot?.kode !== aktif.kode}
							<text
								x={aktif.x}
								y={aktif.y - 14}
								text-anchor={aktif.x > lebar * 0.7 ? 'end' : 'middle'}
								class="fill-fg text-small font-bold"
							>
								{aktif.nama}
							</text>
						{/if}
					</svg>
					{#if disorot}
						<div
							class="pointer-events-none absolute z-tooltip -translate-x-1/2 -translate-y-full rounded-md bg-raised px-3 py-2 text-small shadow-2"
							style:left="{(disorot.x / lebar) * 100}%"
							style:top="{(disorot.y / tinggi) * 100 - 3}%"
						>
							<p class="font-bold">{disorot.nama}</p>
							<p class="num">
								{t('spend.col.ratio')}: {formatAngka(disorot.rasio, bahasaBawaan, 2)}
							</p>
							<p class="num">{t('spend.col.value')}: {formatRupiah(disorot.nilai, bahasaBawaan)}</p>
						</div>
					{/if}
					<figcaption class="mt-2 flex flex-wrap gap-4 text-caption text-muted">
						<span class="inline-flex items-center gap-2"
							><span class="h-3 w-3 rounded-pill bg-primary" aria-hidden="true"></span>{t(
								'spend.active'
							)}</span
						>
						<span>{t('spend.yAxis')} {data.tahun}</span>
					</figcaption>
				</figure>
			{/if}
		</div>
		<p class="flex gap-2 rounded-md bg-info-bg p-3 text-small">
			<Info size={18} strokeWidth={1.75} class="mt-0.5 shrink-0 text-info" aria-hidden="true" />
			{t('spend.sentence')}
		</p>
		{#if kotaAktif?.belanja}
			<p class="text-small text-muted">
				{t('spend.cityLevel', {
					nama: kotaAktif.wilayah.nama,
					jumlah: tampilkanAngka(kotaAktif.belanja.jumlah_paket).nilai,
					nilai: tampilkanAngka(kotaAktif.belanja.total_hps).nilai
				})}
			</p>
		{/if}
		{#if jumlahTanpaRasio(data)}
			<p class="text-small text-muted">{t('spend.noRatio', { n: jumlahTanpaRasio(data) })}</p>
		{/if}
	{/if}
</section>
