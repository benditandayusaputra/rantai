<script lang="ts">
	import { tangkap } from '$lib/api/klien';
	import { konfirmasiLaporan } from '$lib/api/laporan';
	import type { DetailLaporan } from '$lib/api/tipe';
	import ErrorState from '$lib/components/ErrorState.svelte';
	import ReportStatusChip from '$lib/components/ReportStatusChip.svelte';
	import ReportTimeline from '$lib/components/ReportTimeline.svelte';
	import SimRibbon from '$lib/components/SimRibbon.svelte';
	import WitnessPanel from '$lib/components/WitnessPanel.svelte';
	import { sesiAktif } from '$lib/features/akun/sesi.svelte';
	import { cariPosisi, susunLinimasa } from '$lib/features/lapor/lapor';
	import { formatJam, formatTanggal } from '$lib/format';
	import { bahasaBawaan, t } from '$lib/i18n';
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

<div class="wrap flex flex-col gap-6 py-6 lg:py-10">
	{#if isGagalMuat}
		<ErrorState title={t('report.notFound')} />
	{:else if laporan}
		{#if laporan.is_simulation}<SimRibbon />{/if}
		<header class="flex flex-col items-start gap-3">
			<ReportStatusChip status={laporan.status} witnesses={laporan.saksi} />
			<h1 class="t-h1">{t(`report.cond.${laporan.kondisi}`)}</h1>
			<p class="text-small text-muted">
				{formatTanggal(new Date(laporan.dibuat_pada), bahasaBawaan, 'long')} · {formatJam(
					new Date(laporan.dibuat_pada)
				)} · {laporan.nama_wilayah} · {t('report.anonymous')}
			</p>
			{#if laporan.paket}
				<a href="/{bahasaBawaan}/paket/{laporan.paket.kode}" class="font-semibold underline"
					>{laporan.paket.nama}</a
				>
			{/if}
		</header>
		<div class="grid gap-6 lg:grid-cols-12">
			<div class="flex flex-col gap-6 lg:col-span-8">
				<section aria-label={t('report.photos')}>
					<ul class="grid grid-cols-2 gap-3">
						{#each laporan.foto as foto, i (foto.url)}
							<li>
								<a href={foto.url} class="block overflow-hidden rounded-lg bg-sunken">
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
					<blockquote class="border-l-link border-line-strong pl-4 font-serif text-body-lg">
						{laporan.catatan}
					</blockquote>
				{/if}
				<section aria-labelledby="judul-linimasa" class="flex flex-col gap-3">
					<h2 id="judul-linimasa" class="t-h3">{t('report.timeline')}</h2>
					<ReportTimeline events={susunLinimasa(laporan.riwayat)} witnesses={laporan.saksi} />
				</section>
			</div>
			<aside class="hidden lg:col-span-4 lg:block">{@render panelSaksi(laporan)}</aside>
		</div>
	{:else}
		<p aria-live="polite" class="text-muted">{t('state.loading')}</p>
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
