<script lang="ts">
	import {
		ArrowRight,
		Calendar,
		EyeOff,
		History,
		Images,
		MapPin,
		Package,
		Quote,
		Scale
	} from '@lucide/svelte';
	import { tangkap } from '$lib/api/klien';
	import { konfirmasiLaporan } from '$lib/api/laporan';
	import type { DetailLaporan } from '$lib/api/tipe';
	import ErrorState from '$lib/components/ErrorState.svelte';
	import ReportStatusChip from '$lib/components/ReportStatusChip.svelte';
	import RightOfReply from '$lib/components/RightOfReply.svelte';
	import ReportTimeline from '$lib/components/ReportTimeline.svelte';
	import SimRibbon from '$lib/components/SimRibbon.svelte';
	import WitnessPanel from '$lib/components/WitnessPanel.svelte';
	import { sesiAktif } from '$lib/features/akun/sesi.svelte';
	import { cariPosisi, susunLinimasa } from '$lib/features/lapor/lapor';
	import { formatJam, formatTanggal } from '$lib/format';
	import { bahasaAktif, t } from '$lib/i18n';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	let laporan = $state<DetailLaporan | null>(null);
	let pesanSaksi = $state('');
	let isMengonfirmasi = $state(false);
	let isGagalMuat = $state(false);

	$effect(() => {
		void data.laporan.then((hasil) => {
			if (hasil.ok) laporan = hasil.data;
			else isGagalMuat = true;
		});
	});

	const sesi = $derived(sesiAktif());
	const bisaSaksi = $derived(laporan?.status === 'MENUNGGU_SAKSI' && sesi !== null);

	async function jadiSaksi() {
		if (!laporan || !sesi) return;
		isMengonfirmasi = true;
		const posisi = await cariPosisi(navigator.geolocation);
		if (posisi.status !== 'akurat') {
			isMengonfirmasi = false;
			pesanSaksi = t(posisi.status === 'kasar' ? 'report.s1.poorTitle' : 'report.s1.failed');
			return;
		}
		const hasil = await tangkap(
			konfirmasiLaporan(fetch, sesi.token_akses, laporan.id, {
				titik: posisi.titik,
				akurasiM: posisi.akurasiM,
				kondisi: laporan.kondisi
			})
		);
		isMengonfirmasi = false;
		if (!hasil.ok) {
			pesanSaksi = hasil.galat.message;
			return;
		}
		pesanSaksi = t('report.witness.done');
		laporan.status = hasil.data.status;
		laporan.saksi = hasil.data.saksi;
	}
</script>

<svelte:head>
	<title
		>{laporan ? t(`report.cond.${laporan.kondisi}`) : t('report.flowTitle')} · {t(
			'brand.name'
		)}</title
	>
</svelte:head>

<div class="wrap flex flex-col gap-6 py-6 lg:gap-8 lg:py-10">
	{#if isGagalMuat}
		<ErrorState title={t('report.notFound')} />
	{:else if laporan}
		<header
			class="card flex flex-col overflow-hidden {laporan.is_simulation
				? 'border-ctrl border-cert-sim'
				: ''}"
		>
			{#if laporan.is_simulation}<SimRibbon />{/if}
			<div class="flex flex-col items-start gap-3 p-5 md:p-6 lg:p-8">
				<ReportStatusChip status={laporan.status} witnesses={laporan.saksi} />
				<h1 class="t-h1">{t(`report.cond.${laporan.kondisi}`)}</h1>
				<ul class="flex flex-wrap gap-x-4 gap-y-2 text-small text-muted">
					<li class="flex items-center gap-1.5">
						<Calendar size={16} strokeWidth={1.75} aria-hidden="true" />
						<span class="num"
							>{formatTanggal(new Date(laporan.dibuat_pada), bahasaAktif, 'long')} · {formatJam(
								new Date(laporan.dibuat_pada)
							)}</span
						>
					</li>
					<li class="flex items-center gap-1.5">
						<MapPin size={16} strokeWidth={1.75} aria-hidden="true" />{laporan.nama_wilayah}
					</li>
					<li class="flex items-center gap-1.5">
						<EyeOff size={16} strokeWidth={1.75} aria-hidden="true" />{t('report.anonymous')}
					</li>
				</ul>
				{#if laporan.paket}
					<a
						href="/{bahasaAktif}/paket/{laporan.paket.kode}"
						class="card-link group mt-1 flex w-full items-center gap-3 p-3 md:max-w-prose"
					>
						<span class="icon-tile h-9 w-9"
							><Package size={18} strokeWidth={1.75} aria-hidden="true" /></span
						>
						<span class="min-w-0 flex-1 break-words font-semibold">{laporan.paket.nama}</span>
						<ArrowRight
							size={18}
							strokeWidth={1.75}
							aria-hidden="true"
							class="shrink-0 text-muted transition-transform duration-fast ease-standard group-hover:translate-x-1"
						/>
					</a>
				{/if}
			</div>
		</header>
		<div class="grid gap-6 lg:grid-cols-12 lg:gap-8">
			<div class="flex flex-col gap-6 lg:col-span-8">
				<section aria-labelledby="judul-foto" class="flex flex-col gap-3">
					<h2 id="judul-foto" class="t-h3 flex items-center gap-2">
						<Images size={20} strokeWidth={1.75} aria-hidden="true" />{t('report.photos')}
					</h2>
					<ul class="grid grid-cols-2 gap-3">
						{#each laporan.foto as foto, i (foto.url)}
							<li>
								<a
									href={foto.url}
									class="block overflow-hidden rounded-lg bg-sunken shadow-1 transition duration-base ease-standard hover:shadow-2"
								>
									<img
										src={foto.url_thumbnail}
										alt={laporan.catatan ||
											t('report.s3.alt', {
												n: i + 1,
												jam: formatJam(new Date(laporan.dibuat_pada))
											})}
										width={foto.lebar}
										height={foto.tinggi}
										class="aspect-square h-auto w-full object-cover"
									/>
								</a>
							</li>
						{/each}
					</ul>
				</section>
				<div class="lg:hidden">
					{@render panelSaksi(laporan)}
				</div>
				{#if laporan.catatan}
					<blockquote class="card flex gap-3 p-5 font-heading text-body-lg">
						<Quote
							size={20}
							strokeWidth={1.75}
							aria-hidden="true"
							class="mt-1 shrink-0 text-line-strong"
						/>
						{laporan.catatan}
					</blockquote>
				{/if}
				<section aria-labelledby="judul-linimasa" class="card flex flex-col gap-4 p-5">
					<h2 id="judul-linimasa" class="t-h3 flex items-center gap-2">
						<History size={20} strokeWidth={1.75} aria-hidden="true" />{t('report.timeline')}
					</h2>
					<ReportTimeline events={susunLinimasa(laporan.riwayat)} witnesses={laporan.saksi} />
				</section>
				{#if laporan.hak_jawab.length}
					<section aria-labelledby="judul-hak-jawab" class="flex flex-col gap-3">
						<h2 id="judul-hak-jawab" class="t-h3 flex items-center gap-2">
							<Scale size={20} strokeWidth={1.75} aria-hidden="true" />{t('reply.title')}
						</h2>
						{#each laporan.hak_jawab as hakJawab (hakJawab.id)}
							<RightOfReply {hakJawab} />
						{/each}
					</section>
				{/if}
			</div>
			<aside class="hidden lg:col-span-4 lg:block">
				<div class="lg:sticky lg:top-24">{@render panelSaksi(laporan)}</div>
			</aside>
		</div>
	{:else}
		<div class="flex flex-col gap-4" aria-busy="true">
			<p class="sr-only" aria-live="polite">{t('state.loading')}</p>
			<div class="skel h-40 rounded-lg"></div>
			<div class="grid gap-6 lg:grid-cols-12">
				<div class="skel h-72 rounded-lg lg:col-span-8"></div>
				<div class="skel h-56 rounded-lg lg:col-span-4"></div>
			</div>
		</div>
	{/if}
</div>

{#snippet panelSaksi(l: DetailLaporan)}
	<div class="flex flex-col gap-2">
		<WitnessPanel
			count={l.saksi.jumlah}
			needed={l.saksi.dibutuhkan}
			canWitness={bisaSaksi}
			reason={sesi ? undefined : t('report.witness.needSignIn')}
			busy={isMengonfirmasi}
			onwitness={jadiSaksi}
		/>
		{#if pesanSaksi}<p role="status" class="text-small">{pesanSaksi}</p>{/if}
	</div>
{/snippet}
