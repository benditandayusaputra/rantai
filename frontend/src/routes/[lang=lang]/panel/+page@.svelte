<script lang="ts">
	import {
		ArrowLeft,
		Building2,
		CircleCheck,
		CircleHelp,
		Eye,
		FileWarning,
		Hammer,
		Handshake,
		House,
		LoaderCircle,
		LogIn,
		Monitor,
		Newspaper,
		Scale
	} from '@lucide/svelte';
	import { onMount } from 'svelte';
	import { tangkap } from '$lib/api/klien';
	import {
		ambilAntreanPanel,
		isiSerahTerima,
		tandaiLaporanSelesai,
		terbitkanHakJawab
	} from '$lib/api/panel';
	import type { AntreanPanel } from '$lib/api/tipe';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import Logo from '$lib/components/Logo.svelte';
	import MaintenanceCounter from '$lib/components/MaintenanceCounter.svelte';
	import ReportStatusChip from '$lib/components/ReportStatusChip.svelte';
	import TextArea from '$lib/components/TextArea.svelte';
	import { inisial, pulihkanSesi, sesiAktif } from '$lib/features/akun/sesi.svelte';
	import {
		bagianAwal,
		bagianUntuk,
		daftarSorotan,
		isAkunPanel,
		isSiapTerbit,
		jumlahSegeraBerakhir,
		panjangMaksHakJawab,
		tanggalLokal,
		type BagianPanel
	} from '$lib/features/panel/panel';
	import { formatTanggal } from '$lib/format';
	import { bahasaAktif, t } from '$lib/i18n';

	const ikonBagian = { sinyal: CircleHelp, laporan: FileWarning, paket: Hammer, umpan: Newspaper };

	let isSesiDiperiksa = $state(false);
	let antrean = $state<AntreanPanel | null>(null);
	let bagian = $state<BagianPanel>('sinyal');
	let indeks = $state(0);
	let isi = $state('');
	let tautan = $state('');
	let isBerwenang = $state(false);
	let tanggalSerahTerima = $state('');
	let nomorBeritaAcara = $state('');
	let isSibuk = $state(false);
	let pesan = $state('');

	const sesi = $derived(sesiAktif());
	const isAnggota = $derived(sesi?.pengguna.peran === 'anggota');
	const sorotan = $derived(antrean ? daftarSorotan(antrean, bagian, bahasaAktif) : []);
	const terpilih = $derived(sorotan[indeks] ?? null);
	const paketTerpilih = $derived(
		bagian === 'paket' ? (antrean?.paket.find((p) => p.kode === terpilih?.id) ?? null) : null
	);
	const hariIni = tanggalLokal(new Date());

	onMount(async () => {
		await pulihkanSesi(fetch);
		isSesiDiperiksa = true;
		await muat();
	});

	async function muat() {
		const s = sesiAktif();
		if (!s || !isAkunPanel(s.pengguna.peran)) return;
		const hasil = await tangkap(ambilAntreanPanel(fetch, s.token_akses));
		if (!hasil.ok) {
			pesan = hasil.galat.message;
			return;
		}
		if (!antrean) bagian = bagianAwal(hasil.data.peran);
		antrean = hasil.data;
	}

	function pilih(baru: BagianPanel, i: number) {
		bagian = baru;
		indeks = i;
		isi = '';
		tautan = '';
		isBerwenang = false;
		tanggalSerahTerima = '';
		nomorBeritaAcara = '';
	}

	async function jalankan(aksi: (token: string) => Promise<unknown>, berhasil: string) {
		if (!sesi) return;
		isSibuk = true;
		const hasil = await tangkap(aksi(sesi.token_akses));
		isSibuk = false;
		pesan = hasil.ok ? berhasil : hasil.galat.message;
		if (hasil.ok) {
			await muat();
			pilih(bagian, indeks);
		}
	}

	function terbitkan(event: SubmitEvent) {
		event.preventDefault();
		const target = terpilih;
		if (!target) return;
		void jalankan(
			(token) =>
				terbitkanHakJawab(fetch, token, {
					target_jenis: target.jenis,
					target_id: target.id,
					isi,
					tautan_dokumen: tautan
				}),
			t('panel.reply.done')
		);
	}

	function simpanSerahTerima(event: SubmitEvent) {
		event.preventDefault();
		const paket = paketTerpilih;
		if (!paket) return;
		void jalankan(
			(token) =>
				isiSerahTerima(fetch, token, paket.kode, {
					tanggal: tanggalSerahTerima,
					nomor_berita_acara: nomorBeritaAcara
				}),
			t('panel.handover.done')
		);
	}
</script>

<svelte:head>
	<title>{t('panel.title')} · {t('brand.name')}</title>
</svelte:head>

{#if !isSesiDiperiksa}
	<p class="p-8 text-muted" aria-live="polite">{t('state.loading')}</p>
{:else if !sesi || !isAkunPanel(sesi.pengguna.peran)}
	<div class="flex min-h-dvh flex-col bg-page">
		<header class="wrap flex h-16 items-center justify-between gap-4 lg:h-18">
			<Logo />
			<a href="/{bahasaAktif}" class="btn btn-ghost btn-sm">
				<ArrowLeft size={18} strokeWidth={1.75} aria-hidden="true" />
				{t('nav.home')}
			</a>
		</header>
		<main
			id="konten"
			tabindex="-1"
			class="wrap flex flex-1 items-center justify-center pb-16 outline-none"
		>
			<div class="card flex w-full max-w-form flex-col items-start gap-4 p-6 shadow-2 md:p-8">
				<span class="icon-tile h-14 w-14"
					><Building2 size={28} strokeWidth={1.75} aria-hidden="true" /></span
				>
				<h1 class="t-h1">{t('panel.title')}</h1>
				<p class="text-muted">{t('panel.forbidden')}</p>
				{#if !sesi}
					<a
						href="/{bahasaAktif}/masuk?lanjut={encodeURIComponent(`/${bahasaAktif}/panel`)}"
						class="btn btn-primary btn-lg btn-block"
					>
						<LogIn size={20} strokeWidth={1.75} aria-hidden="true" />
						{t('action.signIn')}
					</a>
				{/if}
			</div>
		</main>
	</div>
{:else}
	<div class="flex min-h-dvh items-center justify-center p-4 lg:hidden">
		<div class="card flex w-full max-w-form flex-col items-start gap-4 p-6">
			<span class="icon-tile"><Monitor size={20} strokeWidth={1.75} aria-hidden="true" /></span>
			<p>{t('panel.desktopOnly')}</p>
			<a href="/{bahasaAktif}" class="btn btn-outline">
				<House size={18} strokeWidth={1.75} aria-hidden="true" />
				{t('panel.readOnly')}
			</a>
		</div>
	</div>
	<div class="hidden min-h-dvh lg:flex">
		<nav
			aria-label={t('panel.title')}
			class="sticky top-0 flex h-dvh w-66 shrink-0 flex-col gap-6 border-r border-line bg-surface p-5"
		>
			<Logo />
			{#if antrean}
				<ul class="flex flex-col gap-1">
					{#each bagianUntuk(antrean.peran) as kunci (kunci)}
						{@const Ikon = ikonBagian[kunci]}
						<li>
							<button
								type="button"
								aria-current={bagian === kunci ? 'page' : undefined}
								class="group flex min-h-touch w-full items-center gap-3 rounded-md px-3 text-left text-small font-semibold text-muted transition-colors duration-fast ease-standard hover:bg-sunken hover:text-fg aria-[current=page]:bg-secondary aria-[current=page]:text-secondary-on"
								onclick={() => pilih(kunci, 0)}
							>
								<Ikon size={20} strokeWidth={1.75} aria-hidden="true" />
								<span class="flex-1">{t(`panel.queue.${kunci}`)}</span>
								<span
									class="num inline-flex h-6 min-w-6 items-center justify-center rounded-pill bg-sunken px-2 text-caption group-aria-[current=page]:bg-raised"
									>{antrean[kunci].length}</span
								>
							</button>
						</li>
					{/each}
				</ul>
			{/if}
			<p class="mt-auto flex items-center gap-3 border-t border-line pt-4 text-small">
				<span
					class="inline-flex h-9 w-9 shrink-0 items-center justify-center rounded-pill border border-line-strong bg-sunken font-semibold text-fg"
					>{inisial(sesi.pengguna)}</span
				>
				<span class="min-w-0 truncate font-semibold"
					>{sesi.pengguna.nama_tampilan || sesi.pengguna.email}</span
				>
			</p>
		</nav>

		<main id="konten" class="flex min-w-0 flex-1 flex-col gap-6 bg-page p-8 xl:p-10">
			<header class="flex flex-col gap-2">
				<h1 class="t-h1">
					{isAnggota ? t('panel.memberTitle') : antrean?.nama_instansi || t('panel.title')}
				</h1>
				<p class="max-w-prose text-muted">{isAnggota ? t('panel.memberLead') : t('panel.lead')}</p>
			</header>
			{#if pesan}<p role="status" class="card px-4 py-3 text-small font-semibold">{pesan}</p>{/if}

			{#if !antrean}
				<p class="text-muted">{t('state.loading')}</p>
			{:else}
				{#if !isAnggota}
					<dl class="grid max-w-content grid-cols-3 gap-4">
						{#each [{ label: t('panel.queue.sinyal'), n: antrean.sinyal.length }, { label: t('panel.queue.laporan'), n: antrean.laporan.length }, { label: t('panel.kpi.soon'), n: jumlahSegeraBerakhir(antrean) }] as angka (angka.label)}
							<div class="card p-4 lg:p-5">
								<dt class="text-small text-muted">{angka.label}</dt>
								<dd class="n-md">{angka.n}</dd>
							</div>
						{/each}
					</dl>
				{/if}

				<div class="grid gap-6 xl:grid-cols-2">
					<section aria-labelledby="judul-antrean" class="flex flex-col gap-3">
						<h2 id="judul-antrean" class="t-h3">{t(`panel.queue.${bagian}`)}</h2>
						{#if sorotan.length === 0}
							<EmptyState title={t('panel.queue.empty')} />
						{:else}
							<ul
								role="listbox"
								aria-label={t(`panel.queue.${bagian}`)}
								class="flex flex-col gap-2"
							>
								{#each sorotan as item, i (item.id)}
									<li role="option" aria-selected={i === indeks}>
										<button
											type="button"
											class="flex w-full flex-col gap-1 rounded-lg border-ctrl p-4 text-left shadow-1 transition-colors duration-fast ease-standard hover:border-line-strong {i ===
											indeks
												? 'border-primary bg-secondary'
												: 'border-line bg-surface'}"
											onclick={() => pilih(bagian, i)}
										>
											<span class="font-semibold">{item.judul}</span>
											<span class="flex flex-wrap items-center gap-2 text-small text-muted">
												{#if item.status}<ReportStatusChip status={item.status} />{/if}
												<span class:font-bold={item.isPenting} class:text-attn={item.isPenting}
													>{item.keterangan}</span
												>
											</span>
											{#if item.isSimulation}
												<span class="text-caption font-bold text-cert-sim">{t('sim.ribbon')}</span>
											{/if}
										</button>
									</li>
								{/each}
							</ul>
						{/if}
					</section>

					<section aria-labelledby="judul-formulir" class="flex flex-col gap-6">
						{#if !terpilih}
							<h2 id="judul-formulir" class="sr-only">{t('panel.reply.title')}</h2>
							<div class="card flex flex-col items-center gap-3 px-5 py-10 text-center">
								<span class="icon-tile"
									><Scale size={20} strokeWidth={1.75} aria-hidden="true" /></span
								>
								<p class="text-muted">{t('panel.choose')}</p>
							</div>
						{:else}
							{#if paketTerpilih}
								<MaintenanceCounter pemeliharaan={paketTerpilih.pemeliharaan} />
								<form class="card flex flex-col gap-3 p-5" onsubmit={simpanSerahTerima}>
									<h2 class="t-h3 flex items-center gap-2">
										<Handshake size={20} strokeWidth={1.75} aria-hidden="true" class="text-muted" />
										{t('panel.handover.title')}
									</h2>
									{#if paketTerpilih.tanggal_serah_terima}
										<p class="text-small text-muted">
											{t('panel.handover.current', {
												date: formatTanggal(
													new Date(paketTerpilih.tanggal_serah_terima),
													bahasaAktif
												)
											})}
										</p>
									{/if}
									<label class="flex flex-col gap-1 font-semibold">
										{t('panel.handover.date')}
										<input
											type="date"
											required
											max={hariIni}
											bind:value={tanggalSerahTerima}
											class="field font-normal"
										/>
									</label>
									<label class="flex flex-col gap-1 font-semibold">
										{t('panel.handover.ba')}
										<input
											maxlength={100}
											bind:value={nomorBeritaAcara}
											class="field font-normal"
										/>
									</label>
									<button type="submit" class="btn btn-primary self-start" disabled={isSibuk}>
										<Handshake size={18} strokeWidth={1.75} aria-hidden="true" />
										{t('panel.handover.submit')}
									</button>
								</form>
							{/if}

							<form class="card flex flex-col gap-4 p-5" onsubmit={terbitkan}>
								<h2 id="judul-formulir" class="t-h3 flex items-center gap-2">
									<Scale size={20} strokeWidth={1.75} aria-hidden="true" class="text-muted" />
									{t('panel.reply.title')}
								</h2>
								<p class="font-semibold">{terpilih.judul}</p>
								{#if bagian === 'sinyal'}
									<p class="rounded-md border-ctrl border-attn-fill bg-attn-bg p-3 text-small">
										{terpilih.catatan}<br />{t('panel.signalNote')}
									</p>
								{:else if terpilih.catatan && bagian === 'laporan'}
									<blockquote class="border-l-link border-line-strong pl-3 text-small">
										{terpilih.catatan}
									</blockquote>
								{/if}
								<TextArea
									id="isi-hak-jawab"
									label={t('panel.reply.label')}
									help={t('panel.reply.help')}
									maxlength={panjangMaksHakJawab}
									rows={8}
									bind:value={isi}
								/>
								<label class="flex flex-col gap-1 font-semibold">
									{t('panel.reply.link')}
									<input type="url" inputmode="url" bind:value={tautan} class="field font-normal" />
								</label>
								<label class="flex items-start gap-3">
									<input
										type="checkbox"
										bind:checked={isBerwenang}
										class="mt-0.5 h-5 w-5 shrink-0 accent-primary"
									/>
									<span
										>{isAnggota
											? t('panel.reply.authorityMember')
											: t('panel.reply.authority')}</span
									>
								</label>
								<button
									type="submit"
									class="btn btn-primary btn-lg"
									disabled={isSibuk || !isSiapTerbit(isi, isBerwenang)}
									aria-busy={isSibuk}
								>
									{#if isSibuk}
										<LoaderCircle
											size={20}
											strokeWidth={1.75}
											aria-hidden="true"
											class="animate-spin"
										/>
										{t('panel.reply.sending')}
									{:else}
										<Scale size={20} strokeWidth={1.75} aria-hidden="true" />
										{t('panel.reply.submit')}
									{/if}
								</button>
							</form>

							{#if terpilih.status}
								{@const idLaporan = Number(terpilih.id)}
								<div class="flex flex-wrap items-center gap-3">
									<button
										type="button"
										class="btn btn-outline"
										disabled={isSibuk || terpilih.status !== 'DITANGGAPI'}
										aria-describedby="petunjuk-selesai"
										onclick={() =>
											jalankan(
												(token) => tandaiLaporanSelesai(fetch, token, idLaporan),
												t('panel.report.marked')
											)}
									>
										<CircleCheck size={18} strokeWidth={1.75} aria-hidden="true" />
										{t('panel.report.done')}
									</button>
									<span id="petunjuk-selesai" class="text-small text-muted"
										>{t('panel.report.doneHint')}</span
									>
								</div>
							{/if}

							<a href={terpilih.tautanPublik} class="btn btn-text self-start">
								<Eye size={18} strokeWidth={1.75} aria-hidden="true" />
								{t('panel.openPublic')}
							</a>
						{/if}
					</section>
				</div>
			{/if}
		</main>
	</div>
{/if}
