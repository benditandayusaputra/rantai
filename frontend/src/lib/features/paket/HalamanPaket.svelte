<script lang="ts">
	import {
		Building2,
		Calendar,
		Info,
		CircleHelp,
		CircleX,
		ExternalLink,
		FileText,
		Link2,
		MapPin,
		Scale
	} from '@lucide/svelte';
	import type { DetailPaket } from '$lib/api/tipe';
	import TombolKembali from '$lib/components/TombolKembali.svelte';
	import RingkasanProyek from './RingkasanProyek.svelte';
	import ChainStepper from '$lib/components/ChainStepper.svelte';
	import KeyNumberCard from '$lib/components/KeyNumberCard.svelte';
	import MaintenanceCounter from '$lib/components/MaintenanceCounter.svelte';
	import RightOfReply from '$lib/components/RightOfReply.svelte';
	import SimRibbon from '$lib/components/SimRibbon.svelte';
	import { bahasaAktif, t } from '$lib/i18n';
	import { ikonSektor, kodeKecamatanDari, kunciSektor } from './paket';

	let { paket }: { paket: DetailPaket } = $props();

	const IkonSektor = $derived(ikonSektor[paket.sektor]);
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
			href: paket.penyedia ? `/${bahasaAktif}/penyedia/${paket.penyedia.id}` : undefined
		},
		{ label: t('pkg.location'), isi: paket.lokasi.nama_wilayah }
	]);
	const tautanLuar = $derived(
		[
			{ href: paket.url_sumber, label: t('pkg.openSource') },
			{ href: paket.rup?.url ?? '', label: t('pkg.openRup') }
		].filter((tautan) => tautan.href)
	);
	const kelasPil =
		'inline-flex items-center gap-1.5 rounded-xs border border-line bg-raised px-3 py-1';
</script>

<article class="wrap flex flex-col gap-6 pt-6 lg:gap-8 lg:pt-10">
	<TombolKembali />
	<header
		class="card flex flex-col overflow-hidden {paket.is_simulation
			? 'border-ctrl border-cert-sim'
			: ''}"
	>
		{#if paket.is_simulation}<SimRibbon />{/if}
		<div class="flex flex-col gap-4 p-5 md:p-6 lg:p-8">
			<div class="flex flex-wrap items-center gap-2 text-small text-muted">
				<span
					class="inline-flex items-center gap-1.5 rounded-xs bg-secondary px-3 py-1 font-semibold text-secondary-on"
				>
					<IkonSektor size={16} strokeWidth={1.75} aria-hidden="true" />
					{t(kunciSektor[paket.sektor])}
				</span>
				<span class={kelasPil}>
					<Calendar size={16} strokeWidth={1.75} aria-hidden="true" />
					<span class="num">{paket.tahun_anggaran}</span>
				</span>
				<span class={kelasPil}>
					<FileText size={16} strokeWidth={1.75} aria-hidden="true" />
					{t('pkg.rup')}
					<span class="font-mono text-mono">{paket.rup?.kode ?? '-'}</span>
				</span>
			</div>
			<h1 class="t-h1 break-words">{paket.nama}</h1>
			<p class="flex items-start gap-2 text-muted">
				<Building2 size={20} strokeWidth={1.75} aria-hidden="true" class="mt-0.5 shrink-0" />
				{paket.instansi}
			</p>
			{#if paket.is_batal}
				<p
					class="inline-flex items-center gap-2 self-start rounded-md bg-attn-bg px-3 py-2 font-semibold text-attn"
				>
					<CircleX size={20} strokeWidth={1.75} aria-hidden="true" />
					{t('pkg.cancelled')}
				</p>
			{/if}
		</div>
	</header>

	{#if angkaKunci.length}
		<div class="grid gap-3 md:grid-cols-4 lg:gap-5">
			{#each angkaKunci as item (item.label)}
				<KeyNumberCard label={item.label} angka={item.angka} />
			{/each}
		</div>
	{/if}

	<RingkasanProyek {paket} />

	<section aria-labelledby="judul-rantai" class="card p-5 md:p-6 lg:p-8">
		<h2 id="judul-rantai" class="t-h2 flex items-center gap-3">
			<span class="icon-tile"><Link2 size={20} strokeWidth={1.75} aria-hidden="true" /></span>
			{t('chain.title')}
		</h2>
		<p class="mb-6 mt-3 max-w-prose text-small text-muted">{t('chain.explainer')}</p>
		<ChainStepper rantai={paket.rantai} />
	</section>

	<div class="grid gap-8 lg:grid-cols-12 lg:gap-6">
		<div class="flex flex-col gap-8 lg:col-span-8">
			<section aria-labelledby="judul-sinyal" class="flex flex-col gap-4">
				<h2 id="judul-sinyal" class="t-h2 flex items-center gap-3">
					<span
						class="icon-tile {paket.sinyal.length
							? 'bg-attn-bg text-attn'
							: 'bg-sunken text-muted'}"
						><CircleHelp size={20} strokeWidth={1.75} aria-hidden="true" /></span
					>
					{t(paket.sinyal.length ? 'signal.sectionTitle' : 'project.assessmentTitle')}
				</h2>
				{#if paket.sinyal.length}
					<ul class="flex flex-col gap-3">
						{#each paket.sinyal as sinyal (sinyal.jenis)}
							<li
								class="flex flex-col gap-2 rounded-lg border-ctrl border-attn-fill bg-attn-bg p-4 md:p-5"
							>
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
					<p class="font-heading italic text-muted">{t('signal.disclaimer')}</p>
				{:else}
					<p
						class="flex items-center gap-2 rounded-lg border border-dashed border-line-strong p-4 text-small text-muted"
					>
						<Info size={18} strokeWidth={1.75} aria-hidden="true" class="shrink-0" />
						{t('signal.none')}
					</p>
				{/if}
			</section>

			<section aria-labelledby="judul-hak-jawab" class="flex flex-col gap-4">
				<h2 id="judul-hak-jawab" class="t-h2 flex items-center gap-3">
					<span class="icon-tile bg-info-bg text-info"
						><Scale size={20} strokeWidth={1.75} aria-hidden="true" /></span
					>
					{t('reply.title')}
				</h2>
				{#each paket.hak_jawab as hakJawab (hakJawab.id)}
					<RightOfReply {hakJawab} />
				{:else}
					<p
						class="flex items-center gap-2 rounded-lg border border-dashed border-line-strong p-4 text-small text-muted"
					>
						<Scale size={18} strokeWidth={1.75} aria-hidden="true" class="shrink-0" />
						{t('reply.none')}
					</p>
				{/each}
			</section>
		</div>

		<aside class="flex flex-col gap-4 lg:col-span-4">
			<MaintenanceCounter pemeliharaan={paket.pemeliharaan} />
			<section aria-labelledby="judul-detail" class="card p-5">
				<h2 id="judul-detail" class="flex items-center gap-2 font-bold">
					<FileText size={20} strokeWidth={1.75} aria-hidden="true" />
					{t('pkg.details')}
				</h2>
				<dl class="mt-3 flex flex-col divide-y divide-line text-small">
					{#each rincian as item (item.label)}
						<div class="flex flex-col gap-0.5 py-2.5">
							<dt class="text-caption text-muted">{item.label}</dt>
							<dd class="break-words font-semibold">
								{#if 'href' in item && item.href}
									<a href={item.href} class="link">{item.isi}</a>
								{:else}
									{item.isi || '-'}
								{/if}
							</dd>
						</div>
					{/each}
				</dl>
				<div class="mt-4 flex flex-col gap-2">
					{#each tautanLuar as tautan (tautan.label)}
						<a href={tautan.href} target="_blank" rel="noopener" class="btn btn-outline btn-block">
							<ExternalLink size={18} strokeWidth={1.75} aria-hidden="true" />
							{tautan.label}
							<span class="sr-only">{t('a11y.opensNewTab')}</span>
						</a>
					{/each}
					{#if paket.lokasi.tingkat === 'kecamatan' || paket.lokasi.tingkat === 'kelurahan'}
						<a
							href="/{bahasaAktif}/kecamatan/{kodeKecamatanDari(paket.lokasi.kode_wilayah)}"
							class="btn btn-secondary btn-block"
						>
							<MapPin size={18} strokeWidth={1.75} aria-hidden="true" />
							{t('map.openArea')}
						</a>
					{/if}
				</div>
			</section>
		</aside>
	</div>
</article>
