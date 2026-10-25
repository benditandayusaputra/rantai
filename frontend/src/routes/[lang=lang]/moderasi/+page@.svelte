<script lang="ts">
	import {
		ArrowLeft,
		ExternalLink,
		EyeOff,
		FileWarning,
		Keyboard,
		LogIn,
		Monitor,
		Newspaper,
		ShieldCheck
	} from '@lucide/svelte';
	import { onMount } from 'svelte';
	import { tangkap } from '$lib/api/klien';
	import { ambilAntreanModerasi, moderasiLaporan, sembunyikanUmpan } from '$lib/api/moderasi';
	import type { AksiModerasi, AntreanModerasi } from '$lib/api/tipe';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import Logo from '$lib/components/Logo.svelte';
	import { inisial, pulihkanSesi, sesiAktif } from '$lib/features/akun/sesi.svelte';
	import DetailLaporanModerasi from '$lib/features/moderasi/DetailLaporanModerasi.svelte';
	import { bacaPintasan, geser, sisihkan } from '$lib/features/moderasi/moderasi';
	import { formatJam, formatTanggal } from '$lib/format';
	import { bahasaAktif, t } from '$lib/i18n';

	let isSesiDiperiksa = $state(false);
	let antrean = $state<AntreanModerasi | null>(null);
	let bagian = $state<'laporan' | 'umpan'>('laporan');
	let indeks = $state(0);
	let isSibuk = $state(false);
	let pesan = $state('');
	let alasanUmpan = $state('');
	let detail = $state<DetailLaporanModerasi>();

	const sesi = $derived(sesiAktif());
	const isModerator = $derived(
		sesi?.pengguna.peran === 'moderator' || sesi?.pengguna.peran === 'admin'
	);
	const terpilih = $derived(antrean?.laporan[indeks] ?? null);

	onMount(async () => {
		await pulihkanSesi(fetch);
		isSesiDiperiksa = true;
		await muat();
	});

	async function muat() {
		const s = sesiAktif();
		if (!s) return;
		const hasil = await tangkap(ambilAntreanModerasi(fetch, s.token_akses));
		if (hasil.ok) antrean = hasil.data;
		else pesan = hasil.galat.message;
	}

	async function jalankan(id: number, aksi: AksiModerasi) {
		if (!sesi || !antrean) return;
		isSibuk = true;
		const hasil = await tangkap(moderasiLaporan(fetch, sesi.token_akses, id, aksi));
		isSibuk = false;
		pesan = hasil.ok ? t('mod.done') : hasil.galat.message;
		if (!hasil.ok) return;
		if (aksi.aksi === 'sembunyikan_foto') {
			const foto = antrean.laporan
				.find((l) => l.id === id)
				?.foto.find((f) => f.id === aksi.foto_id);
			if (foto) foto.disembunyikan = true;
			return;
		}
		antrean.laporan = sisihkan(antrean.laporan, id);
		indeks = Math.min(indeks, Math.max(0, antrean.laporan.length - 1));
	}

	async function sembunyikan(id: number) {
		if (!sesi || !antrean) return;
		isSibuk = true;
		const hasil = await tangkap(sembunyikanUmpan(fetch, sesi.token_akses, id, alasanUmpan.trim()));
		isSibuk = false;
		pesan = hasil.ok ? t('mod.done') : hasil.galat.message;
		if (hasil.ok) antrean.umpan = sisihkan(antrean.umpan, id);
	}

	function tanganiTombol(event: KeyboardEvent) {
		if (bagian !== 'laporan' || !antrean) return;
		const pintasan = bacaPintasan(event);
		if (!pintasan) return;
		event.preventDefault();
		if (pintasan === 'berikut') indeks = geser(indeks, 1, antrean.laporan.length);
		if (pintasan === 'sebelum') indeks = geser(indeks, -1, antrean.laporan.length);
		if (pintasan === 'verifikasi' && terpilih) void jalankan(terpilih.id, { aksi: 'verifikasi' });
		if (pintasan === 'tolak') detail?.fokusTolak();
	}

	const navigasi = $derived([
		{
			kunci: 'laporan' as const,
			label: t('mod.queue.reports'),
			ikon: FileWarning,
			jumlah: antrean?.laporan.length
		},
		{
			kunci: 'umpan' as const,
			label: t('mod.queue.feed'),
			ikon: Newspaper,
			jumlah: antrean?.umpan.length
		}
	]);
</script>

<svelte:head>
	<title>{t('mod.title')} · {t('brand.name')}</title>
</svelte:head>

<svelte:window onkeydown={tanganiTombol} />

{#if !isSesiDiperiksa}
	<p class="p-8 text-muted" aria-live="polite">{t('state.loading')}</p>
{:else if !sesi || !isModerator}
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
					><ShieldCheck size={28} strokeWidth={1.75} aria-hidden="true" /></span
				>
				<h1 class="t-h1">{t('mod.title')}</h1>
				<p class="text-muted">{t('mod.forbidden')}</p>
				{#if !sesi}
					<a
						href="/{bahasaAktif}/masuk?lanjut={encodeURIComponent(`/${bahasaAktif}/moderasi`)}"
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
			<p>{t('mod.desktopOnly')}</p>
		</div>
	</div>
	<div class="hidden min-h-dvh lg:flex">
		<nav
			aria-label={t('mod.title')}
			class="sticky top-0 flex h-dvh w-66 shrink-0 flex-col gap-6 border-r border-line bg-surface p-5"
		>
			<Logo />
			<ul class="flex flex-col gap-1">
				{#each navigasi as item (item.kunci)}
					{@const Ikon = item.ikon}
					<li>
						<button
							type="button"
							aria-current={bagian === item.kunci ? 'page' : undefined}
							class="group flex min-h-touch w-full items-center gap-3 rounded-md px-3 text-left text-small font-semibold text-muted transition-colors duration-fast ease-standard hover:bg-sunken hover:text-fg aria-[current=page]:bg-secondary aria-[current=page]:text-secondary-on"
							onclick={() => (bagian = item.kunci)}
						>
							<Ikon size={20} strokeWidth={1.75} aria-hidden="true" />
							<span class="flex-1">{item.label}</span>
							{#if item.jumlah !== undefined}<span
									class="num inline-flex h-6 min-w-6 items-center justify-center rounded-pill bg-sunken px-2 text-caption group-aria-[current=page]:bg-raised"
									>{item.jumlah}</span
								>{/if}
						</button>
					</li>
				{/each}
			</ul>
			<p class="mt-auto flex items-center gap-3 border-t border-line pt-4 text-small">
				<span
					class="inline-flex h-9 w-9 shrink-0 items-center justify-center rounded-pill bg-primary font-bold text-primary-on"
					>{inisial(sesi.pengguna)}</span
				>
				<span class="min-w-0 truncate font-semibold"
					>{sesi.pengguna.nama_tampilan || sesi.pengguna.email}</span
				>
			</p>
		</nav>
		<main id="konten" class="flex min-w-0 flex-1 flex-col gap-5 bg-page p-8 xl:p-10">
			<header class="flex flex-wrap items-end justify-between gap-4">
				<h1 class="t-h1">{bagian === 'laporan' ? t('mod.queue.reports') : t('mod.queue.feed')}</h1>
				{#if bagian === 'laporan'}
					<p class="flex items-center gap-2 text-small text-muted">
						<Keyboard size={18} strokeWidth={1.75} aria-hidden="true" />
						{t('mod.shortcuts')}
					</p>
				{/if}
			</header>
			{#if pesan}<p role="status" class="card px-4 py-3 text-small font-semibold">{pesan}</p>{/if}
			{#if !antrean}
				<p class="text-muted">{t('state.loading')}</p>
			{:else if bagian === 'laporan'}
				{#if antrean.laporan.length === 0}
					<EmptyState title={t('mod.queue.empty')} />
				{:else}
					<div class="flex gap-6">
						<ul
							role="listbox"
							aria-label={t('mod.queue.reports')}
							class="flex w-80 shrink-0 flex-col gap-2"
						>
							{#each antrean.laporan as item, i (item.id)}
								<li role="option" aria-selected={i === indeks}>
									<button
										type="button"
										class="flex w-full gap-3 rounded-lg border-ctrl p-3 text-left shadow-1 transition-colors duration-fast ease-standard hover:border-line-strong {i ===
										indeks
											? 'border-primary bg-secondary'
											: 'border-line bg-surface'}"
										onclick={() => (indeks = i)}
									>
										{#if item.foto[0]}
											<img
												src={item.foto[0].url_thumbnail}
												alt=""
												class="h-14 w-14 shrink-0 rounded-md object-cover"
											/>
										{/if}
										<span class="flex flex-col gap-0.5">
											<span class="font-semibold">{t(`report.cond.${item.kondisi}`)}</span>
											<span class="text-small text-muted"
												>{item.nama_wilayah} · {formatTanggal(
													new Date(item.dibuat_pada),
													bahasaAktif
												)}
												{formatJam(new Date(item.dibuat_pada))}</span
											>
											{#if item.is_simulation}<span class="text-caption font-bold text-cert-sim"
													>{t('sim.ribbon')}</span
												>{/if}
										</span>
									</button>
								</li>
							{/each}
						</ul>
						{#if terpilih}
							<div class="max-w-prose flex-1">
								{#key terpilih.id}
									<DetailLaporanModerasi
										bind:this={detail}
										laporan={terpilih}
										alasanTolak={antrean.alasan_tolak}
										busy={isSibuk}
										onverifikasi={() => jalankan(terpilih.id, { aksi: 'verifikasi' })}
										ontolak={(alasan, isiPesan) =>
											jalankan(terpilih.id, { aksi: 'tolak', alasan, pesan: isiPesan })}
										onsembunyikan={(fotoId, alasan) =>
											jalankan(terpilih.id, { aksi: 'sembunyikan_foto', foto_id: fotoId, alasan })}
										onlewati={() => (indeks = geser(indeks, 1, antrean?.laporan.length ?? 0))}
									/>
								{/key}
							</div>
						{/if}
					</div>
				{/if}
			{:else}
				<label class="flex max-w-form flex-col gap-1 text-small font-semibold">
					{t('mod.hideReason')}
					<input bind:value={alasanUmpan} class="field font-normal" />
				</label>
				{#if antrean.umpan.length === 0}
					<EmptyState title={t('mod.queue.empty')} />
				{/if}
				<ul class="flex max-w-prose flex-col gap-3">
					{#each antrean.umpan as item (item.id)}
						<li class="card flex flex-wrap items-start gap-4 p-4">
							<div class="flex min-w-0 flex-1 flex-col gap-1">
								<a href={item.url} class="link inline-flex gap-1" target="_blank" rel="noopener"
									>{item.judul}<ExternalLink
										size={16}
										strokeWidth={1.75}
										aria-hidden="true"
										class="mt-1 shrink-0"
									/><span class="sr-only">{t('a11y.opensNewTab')}</span></a
								>
								<p class="text-small text-muted">
									{item.media_nama} · {formatTanggal(new Date(item.terbit_pada), bahasaAktif)}
								</p>
								<p class="text-small">{t('mod.members', { nama: item.anggota.join(', ') })}</p>
							</div>
							<button
								type="button"
								class="btn btn-outline"
								disabled={isSibuk || !alasanUmpan.trim()}
								onclick={() => sembunyikan(item.id)}
							>
								<EyeOff size={18} strokeWidth={1.75} aria-hidden="true" />
								{t('mod.hideFeed')}
							</button>
						</li>
					{/each}
				</ul>
			{/if}
		</main>
	</div>
{/if}
