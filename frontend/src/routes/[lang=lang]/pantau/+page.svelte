<script lang="ts">
	import { Bell, ChevronRight, Info, LogIn, MapPin, Plus, X } from '@lucide/svelte';
	import { onMount } from 'svelte';
	import { cariWilayah, tangkap } from '$lib/api/klien';
	import {
		ambilNotifikasi,
		berhentiLangganan,
		daftarLangganan,
		langgani,
		tandaiDibaca
	} from '$lib/api/pantau';
	import type { HalamanNotifikasi, Langganan, Notifikasi, WilayahDitemukan } from '$lib/api/tipe';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import ErrorState from '$lib/components/ErrorState.svelte';
	import RegionSearch from '$lib/components/RegionSearch.svelte';
	import { pulihkanSesi, sesiAktif } from '$lib/features/akun/sesi.svelte';
	import { kelompokkanPerHari, tampilanKejadian } from '$lib/features/pantau/pantau';
	import SakelarPush from '$lib/features/pantau/SakelarPush.svelte';
	import SakelarRingkasanEmail from '$lib/features/pantau/SakelarRingkasanEmail.svelte';
	import { formatJam } from '$lib/format';
	import { bahasaAktif, t } from '$lib/i18n';

	let isSesiDiperiksa = $state(false);
	let langganan = $state<Langganan[]>([]);
	let kabar = $state<HalamanNotifikasi | null>(null);
	let isGagal = $state(false);
	let kata = $state('');

	const sesi = $derived(sesiAktif());
	const kelompok = $derived(kabar ? kelompokkanPerHari(kabar.item, new Date()) : []);

	onMount(async () => {
		await pulihkanSesi(fetch);
		isSesiDiperiksa = true;
		await muat();
	});

	async function muat() {
		const s = sesiAktif();
		if (!s) return;
		const [daftar, feed] = await Promise.all([
			tangkap(daftarLangganan(fetch, s.token_akses)),
			tangkap(ambilNotifikasi(fetch, s.token_akses))
		]);
		if (daftar.ok) langganan = daftar.data;
		if (feed.ok) kabar = feed.data;
		isGagal = !daftar.ok || !feed.ok;
	}

	async function tambah(wilayah: WilayahDitemukan) {
		if (!sesi) return;
		const hasil = await tangkap(langgani(fetch, sesi.token_akses, wilayah.kode));
		if (!hasil.ok) return;
		kata = '';
		await muat();
	}

	async function hapus(l: Langganan) {
		if (!sesi) return;
		const hasil = await tangkap(berhentiLangganan(fetch, sesi.token_akses, l.id));
		if (hasil.ok) langganan = langganan.filter((item) => item.id !== l.id);
	}

	function buka(n: Notifikasi) {
		if (!sesi || n.dibaca || !kabar) return;
		n.dibaca = true;
		kabar.belum_dibaca = Math.max(0, kabar.belum_dibaca - 1);
		void tangkap(tandaiDibaca(fetch, sesi.token_akses, n.id));
	}
</script>

<svelte:head>
	<title>{t('watch.title')} · {t('brand.name')}</title>
</svelte:head>

<div class="wrap flex flex-col gap-8 pb-10">
	<header class="page-head max-w-prose">
		<span class="icon-tile h-12 w-12"><Bell size={24} strokeWidth={1.75} aria-hidden="true" /></span
		>
		<h1 class="t-h1">{t('watch.title')}</h1>
		<p class="t-lead">{t('watch.lead')}</p>
	</header>

	{#if !isSesiDiperiksa}
		<p class="text-muted" aria-live="polite">{t('state.loading')}</p>
	{:else if !sesi}
		<EmptyState title={t('watch.signIn')} description={t('watch.signInBody')}>
			{#snippet actions()}
				<a
					href="/{bahasaAktif}/masuk?lanjut={encodeURIComponent(`/${bahasaAktif}/pantau`)}"
					class="btn btn-primary"
				>
					<LogIn size={18} strokeWidth={1.75} aria-hidden="true" />
					{t('action.signIn')}
				</a>
			{/snippet}
		</EmptyState>
	{:else}
		<div class="grid gap-8 lg:grid-cols-12">
			<div class="flex flex-col gap-6 lg:col-span-8">
				<section aria-labelledby="judul-dipantau" class="card flex flex-col gap-4 p-5">
					<h2 id="judul-dipantau" class="t-h3 flex items-center gap-2">
						<MapPin size={20} strokeWidth={1.75} aria-hidden="true" class="text-muted" />
						{t('watch.areas')}
					</h2>
					{#if langganan.length}
						<ul class="flex flex-wrap gap-2">
							{#each langganan as l (l.id)}
								<li
									class="inline-flex min-h-touch items-center gap-1 rounded-md bg-secondary pl-4 pr-1 text-small font-semibold text-secondary-on"
								>
									<a
										href="/{bahasaAktif}/kecamatan/{l.kode_wilayah}"
										class="underline-offset-4 hover:underline">{l.nama_wilayah}</a
									>
									<button
										type="button"
										class="inline-flex h-11 w-11 items-center justify-center rounded-md hover:bg-raised"
										aria-label={t('watch.remove', { nama: l.nama_wilayah })}
										onclick={() => hapus(l)}
									>
										<X size={18} strokeWidth={1.75} aria-hidden="true" />
									</button>
								</li>
							{/each}
						</ul>
					{/if}
					<div class="flex max-w-form flex-col gap-2">
						<span class="flex items-center gap-2 text-small font-semibold">
							<Plus size={16} strokeWidth={1.75} aria-hidden="true" />
							{t('watch.addLabel')}
						</span>
						<RegionSearch bind:value={kata} cari={(q) => cariWilayah(fetch, q)} onselect={tambah} />
					</div>
				</section>

				{#if isGagal}
					<ErrorState title={t('watch.failed')} />
				{:else if kabar}
					{#if kabar.belum_dibaca}
						<p class="text-small font-semibold" aria-live="polite">
							{t('watch.unread', { n: kabar.belum_dibaca })}
						</p>
					{/if}
					{#if kabar.item.length === 0}
						<EmptyState
							title={langganan.length ? t('watch.emptyFeed') : t('watch.emptyTitle')}
							description={t('watch.emptyBody')}
						/>
					{/if}
					{#each kelompok as grup (grup.hari)}
						<section class="flex flex-col gap-3" aria-label={grup.hari}>
							<h2 class="t-over">{grup.hari}</h2>
							<ul class="flex flex-col gap-2">
								{#each grup.item as n (n.id)}
									{@const tampil = tampilanKejadian[n.jenis]}
									{@const Ikon = tampil.ikon}
									<li>
										<a
											href="/{bahasaAktif}{n.tautan}"
											onclick={() => buka(n)}
											class="card-link group flex items-center gap-3 p-4"
										>
											<span
												class="inline-flex h-10 w-10 shrink-0 items-center justify-center rounded-pill {tampil.kelas}"
											>
												<Ikon size={18} strokeWidth={1.75} aria-hidden="true" />
											</span>
											<span class="flex min-w-0 flex-1 flex-col gap-0.5">
												<span class="text-caption text-muted">
													{t(tampil.label)} · {formatJam(new Date(n.terjadi_pada))}
													{#if n.is_simulation}· {t('sim.ribbon')}{/if}
												</span>
												<span class={n.dibaca ? '' : 'font-bold'}>{n.judul}</span>
												<span class="text-small text-muted">{n.nama_kecamatan}</span>
											</span>
											<ChevronRight
												size={18}
												strokeWidth={1.75}
												aria-hidden="true"
												class="shrink-0 text-muted transition-transform duration-fast ease-standard group-hover:translate-x-1"
											/>
										</a>
									</li>
								{/each}
							</ul>
						</section>
					{/each}
				{/if}
			</div>
			<aside class="flex flex-col gap-4 lg:col-span-4">
				<SakelarPush token={sesi.token_akses} />
				<SakelarRingkasanEmail token={sesi.token_akses} />
				<section class="flex flex-col gap-2 rounded-lg border border-line bg-info-bg p-5">
					<h2 class="flex items-center gap-2 font-bold">
						<Info size={18} strokeWidth={1.75} class="text-info" aria-hidden="true" />
						{t('watch.meaning')}
					</h2>
					<p class="text-small">{t('watch.meaningBody')}</p>
				</section>
			</aside>
		</div>
	{/if}
</div>
