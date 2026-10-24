<script lang="ts">
	import RightOfReply from '$lib/components/RightOfReply.svelte';
	import { CircleHelp, ExternalLink } from '@lucide/svelte';
	import type { DetailPaket } from '$lib/api/tipe';
	import ChainStepper from '$lib/components/ChainStepper.svelte';
	import KeyNumberCard from '$lib/components/KeyNumberCard.svelte';
	import MaintenanceCounter from '$lib/components/MaintenanceCounter.svelte';
	import SimRibbon from '$lib/components/SimRibbon.svelte';
	import { bahasaBawaan, t } from '$lib/i18n';
	import { kodeKecamatanDari, kunciSektor } from './paket';

	let { paket }: { paket: DetailPaket } = $props();

	const angkaKunci = $derived(
		[
			{ label: t('pkg.budget'), angka: paket.pagu ?? paket.rup?.pagu ?? null },
			{ label: t('pkg.hps'), angka: paket.hps },
			{ label: t('pkg.contract'), angka: paket.nilai_kontrak },
			{ label: t('pkg.participants'), angka: paket.jumlah_peserta }
		].flatMap((item) => (item.angka ? [{ label: item.label, angka: item.angka }] : []))
	);
	const rincian = $derived([
		{ label: t('pkg.agency'), isi: paket.instansi },
		{ label: t('pkg.year'), isi: String(paket.tahun_anggaran) },
		{ label: t('pkg.method'), isi: paket.metode },
		{ label: t('pkg.type'), isi: paket.jenis_pengadaan },
		{ label: t('pkg.stage'), isi: paket.tahap },
		{
			label: t('pkg.vendor'),
			isi: paket.penyedia?.nama ?? t('pkg.noVendor'),
			href: paket.penyedia ? `/${bahasaBawaan}/penyedia/${paket.penyedia.id}` : undefined
		},
		{ label: t('pkg.location'), isi: paket.lokasi.nama_wilayah }
	]);
	const kelasTautan =
		'inline-flex min-h-touch items-center gap-2 font-semibold text-secondary-on underline';
</script>

<article class="wrap flex flex-col gap-8 pt-6 lg:pt-10">
	<header
		class="flex flex-col overflow-hidden rounded-lg {paket.is_simulation
			? 'border-ctrl border-cert-sim'
			: ''}"
	>
		{#if paket.is_simulation}<SimRibbon />{/if}
		<div class="flex flex-col gap-3 {paket.is_simulation ? 'p-4 md:p-5' : ''}">
			<p class="text-small text-muted">
				{t(kunciSektor[paket.sektor])} · {t('pkg.rup')}
				<span class="font-mono text-mono">{paket.rup?.kode ?? '-'}</span>
			</p>
			<h1 class="t-h1">{paket.nama}</h1>
			<p class="text-muted">{paket.instansi}</p>
			{#if paket.is_batal}<p class="font-semibold text-attn">{t('pkg.cancelled')}</p>{/if}
		</div>
	</header>

	{#if angkaKunci.length}
		<div class="grid grid-cols-2 gap-3 md:grid-cols-4 lg:gap-5">
			{#each angkaKunci as item (item.label)}
				<KeyNumberCard label={item.label} angka={item.angka} />
			{/each}
		</div>
	{/if}

	<section
		aria-labelledby="judul-rantai"
		class="rounded-lg border border-line bg-surface p-4 md:p-6"
	>
		<h2 id="judul-rantai" class="t-h2 mb-6">{t('chain.title')}</h2>
		<ChainStepper rantai={paket.rantai} />
	</section>

	<div class="grid gap-8 lg:grid-cols-12 lg:gap-5">
		<section aria-labelledby="judul-sinyal" class="flex flex-col gap-4 lg:col-span-8">
			<h2 id="judul-sinyal" class="t-h2">{t('signal.sectionTitle')}</h2>
			{#if paket.sinyal.length}
				<ul class="flex flex-col gap-3">
					{#each paket.sinyal as sinyal (sinyal.jenis)}
						<li class="flex flex-col gap-2 rounded-lg border-ctrl border-attn-fill bg-attn-bg p-4">
							<span class="flex items-center gap-2 font-bold text-attn">
								<CircleHelp size={20} strokeWidth={1.75} aria-hidden="true" />
								{t(`signal.${sinyal.jenis}`)}
							</span>
							<span class="text-small">{sinyal.alasan}</span>
							{#each sinyal.hak_jawab as hakJawab (hakJawab.id)}
								<RightOfReply {hakJawab} />
							{/each}
						</li>
					{/each}
				</ul>
				<p class="font-serif italic text-muted">{t('signal.disclaimer')}</p>
			{:else}
				<p class="text-small text-muted">{t('signal.none')}</p>
			{/if}
			<h2 id="judul-hak-jawab" class="t-h2 mt-4">{t('reply.title')}</h2>
			{#each paket.hak_jawab as hakJawab (hakJawab.id)}
				<RightOfReply {hakJawab} />
			{:else}
				<p class="text-small text-muted">{t('reply.none')}</p>
			{/each}
		</section>

		<aside class="flex flex-col gap-4 lg:col-span-4">
			<MaintenanceCounter pemeliharaan={paket.pemeliharaan} />
			<section
				aria-labelledby="judul-detail"
				class="rounded-lg border border-line bg-surface p-4 md:p-5"
			>
				<h2 id="judul-detail" class="font-bold">{t('pkg.details')}</h2>
				<dl class="mt-3 flex flex-col gap-3 text-small">
					{#each rincian as item (item.label)}
						<div>
							<dt class="text-caption text-muted">{item.label}</dt>
							<dd>
								{#if 'href' in item && item.href}
									<a href={item.href} class="font-semibold underline">{item.isi}</a>
								{:else}
									{item.isi || '-'}
								{/if}
							</dd>
						</div>
					{/each}
				</dl>
				<div class="mt-3 flex flex-col">
					{#if paket.url_sumber}
						<a href={paket.url_sumber} target="_blank" rel="noopener" class={kelasTautan}>
							{t('pkg.openSource')}
							<ExternalLink size={16} strokeWidth={1.75} aria-hidden="true" />
							<span class="sr-only">{t('a11y.opensNewTab')}</span>
						</a>
					{/if}
					{#if paket.rup?.url}
						<a href={paket.rup.url} target="_blank" rel="noopener" class={kelasTautan}>
							{t('pkg.openRup')}
							<ExternalLink size={16} strokeWidth={1.75} aria-hidden="true" />
							<span class="sr-only">{t('a11y.opensNewTab')}</span>
						</a>
					{/if}
					{#if paket.lokasi.tingkat === 'kecamatan' || paket.lokasi.tingkat === 'kelurahan'}
						<a
							href="/{bahasaBawaan}/kecamatan/{kodeKecamatanDari(paket.lokasi.kode_wilayah)}"
							class={kelasTautan}>{t('map.openArea')}</a
						>
					{/if}
				</div>
			</section>
		</aside>
	</div>
</article>
