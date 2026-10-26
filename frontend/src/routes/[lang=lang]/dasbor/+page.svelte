<script lang="ts">
	import {
		CalendarRange,
		ChartColumn,
		ChartNoAxesCombined,
		GraduationCap,
		HandCoins,
		Landmark,
		TrendingUp,
		Wallet
	} from '@lucide/svelte';
	import { invalidateAll } from '$app/navigation';
	import { navigating, page } from '$app/state';
	import ChoiceChips from '$lib/components/ChoiceChips.svelte';
	import ErrorState from '$lib/components/ErrorState.svelte';
	import LensaLayanan from '$lib/features/dasbor/LensaLayanan.svelte';
	import LensaTanggungJawab from '$lib/features/dasbor/LensaTanggungJawab.svelte';
	import LensaUang from '$lib/features/dasbor/LensaUang.svelte';
	import GrafikGaris from '$lib/features/dasbor/GrafikGaris.svelte';
	import KartuGrafik from '$lib/features/dasbor/KartuGrafik.svelte';
	import {
		angkaDari,
		daftarLensa,
		tahunTersedia,
		tautanLensa,
		tautanTahun,
		type Lensa
	} from '$lib/features/dasbor/dasbor';
	import { pilihanJenjang, teksJenjang } from '$lib/features/pendidikan/jenjang';
	import type { Porsi, Seri } from '$lib/api/tipe';
	import { formatAngka, formatRupiah } from '$lib/format';
	import { bahasaAktif, t, type KunciPesan } from '$lib/i18n';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	const infoLensa: Record<Lensa, { ikon: typeof Wallet; nama: KunciPesan; tanya: KunciPesan }> = {
		layanan: { ikon: GraduationCap, nama: 'lens.services', tanya: 'dash.services.q' },
		uang: { ikon: Wallet, nama: 'lens.money', tanya: 'dash.money.q' },
		'tanggung-jawab': {
			ikon: Landmark,
			nama: 'lens.accountability',
			tanya: 'dash.accountability.q'
		}
	};

	const hasil = $derived(data.hasil);
	const pilihanTahun = $derived(hasil.ok ? tahunTersedia(hasil.data) : null);
	const IkonLensa = $derived(infoLensa[data.lensa].ikon);
	const bilangan = (n: number) => formatAngka(n, bahasaAktif);
	const rupiah = (n: number) => formatRupiah(n, bahasaAktif);
	const jumlahDari = (p: Porsi) => p.jumlah;
	const nilaiDari = (p: Porsi) => p.nilai;

	const grafikTren = $derived(
		data.tren.ok
			? [
					{
						judul: t('dash.money.yearly.count'),
						keterangan: t('dash.money.yearly.lead'),
						ikon: ChartColumn,
						seri: data.tren.data.paket,
						nilai: jumlahDari,
						format: bilangan,
						satuan: 'paket'
					},
					{
						judul: t('dash.money.yearly.value'),
						keterangan: t('dash.money.yearly.lead'),
						ikon: HandCoins,
						seri: data.tren.data.paket,
						nilai: nilaiDari,
						format: rupiah,
						satuan: 'rupiah'
					},
					{
						judul: t('dash.acc.yearly.title'),
						keterangan: t('dash.acc.yearly.lead'),
						ikon: CalendarRange,
						seri: data.tren.data.kegiatan,
						nilai: jumlahDari,
						format: bilangan,
						satuan: 'kegiatan'
					}
				]
			: []
	);

	function angkaTerbaru(seri: Seri, nilai: (p: Porsi) => number, satuan: string) {
		const terakhir = seri.isi.at(-1);
		return terakhir ? angkaDari(seri, nilai(terakhir), satuan) : null;
	}

	const tanyaLensa = $derived(
		data.lensa === 'layanan' ? teksJenjang[data.jenjang].dasborTanya : infoLensa[data.lensa].tanya
	);
</script>

<svelte:head>
	<title>{t('dash.title')} · {t('brand.name')}</title>
</svelte:head>

<div class="border-b border-line bg-surface">
	<div class="wrap pt-6 lg:pt-10">
		<div class="flex items-start gap-4">
			<span class="icon-tile hidden h-14 w-14 md:inline-flex">
				<ChartNoAxesCombined size={28} strokeWidth={1.75} aria-hidden="true" />
			</span>
			<div class="min-w-0">
				<h1 class="t-h1">{t('dash.title')}</h1>
				<p class="t-lead mt-2 max-w-prose">{t('dash.lead')}</p>
			</div>
		</div>
		{#if grafikTren.length > 0}
			<section aria-labelledby="judul-tren" class="mt-8">
				<div class="flex items-center gap-3">
					<TrendingUp size={22} strokeWidth={1.75} aria-hidden="true" class="text-primary" />
					<h2 id="judul-tren" class="t-h2">{t('dash.trend.title')}</h2>
				</div>
				<p class="mt-1 max-w-prose text-small text-muted">{t('dash.trend.lead')}</p>
				<div class="mt-4 grid gap-4 md:grid-cols-2 lg:grid-cols-3 lg:gap-6">
					{#each grafikTren as grafik (grafik.judul)}
						<KartuGrafik
							judul={grafik.judul}
							keterangan={grafik.keterangan}
							ikon={grafik.ikon}
							angka={angkaTerbaru(grafik.seri, grafik.nilai, grafik.satuan)}
							kolom={[t('dash.col.period'), grafik.judul]}
							baris={grafik.seri.isi.map((p) => [p.kunci, grafik.format(grafik.nilai(p))])}
						>
							<GrafikGaris
								isi={grafik.seri.isi}
								nilaiDari={grafik.nilai}
								format={grafik.format}
								labelDari={(p) => p.kunci}
							/>
						</KartuGrafik>
					{/each}
				</div>
			</section>
		{/if}
		<nav aria-label={t('dash.lensNav')} class="mt-8">
			<ul class="grid grid-cols-3 md:flex md:gap-2">
				{#each daftarLensa as lensa (lensa)}
					{@const Ikon = infoLensa[lensa].ikon}
					<li>
						<a
							href={tautanLensa(page.url, lensa)}
							data-sveltekit-noscroll
							aria-current={data.lensa === lensa ? 'page' : undefined}
							class="-mb-px flex h-15 flex-col items-center justify-center gap-1 border-b-link border-transparent px-2 text-small font-semibold text-muted transition-colors duration-fast ease-standard hover:text-fg md:h-14 md:flex-row md:gap-2 md:px-4 md:text-body aria-[current=page]:border-primary aria-[current=page]:text-primary"
						>
							<Ikon size={20} strokeWidth={1.75} aria-hidden="true" />
							{t(infoLensa[lensa].nama)}
						</a>
					</li>
				{/each}
			</ul>
		</nav>
	</div>
</div>

<section
	class="wrap mt-8 transition-opacity duration-base ease-standard lg:mt-10"
	class:opacity-60={navigating.to}
	aria-labelledby="judul-lensa"
	aria-busy={Boolean(navigating.to)}
>
	<div class="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
		<div class="flex items-start gap-4">
			<span class="icon-tile hidden md:inline-flex">
				<IkonLensa size={20} strokeWidth={1.75} aria-hidden="true" />
			</span>
			<div>
				<h2 id="judul-lensa" class="t-h2">{t(infoLensa[data.lensa].nama)}</h2>
				<p class="t-lead mt-1 max-w-prose">{t(tanyaLensa)}</p>
			</div>
		</div>
		{#if data.lensa === 'layanan'}
			<div class="flex flex-wrap items-center gap-3">
				<p id="judul-jenjang-dasbor" class="text-small font-semibold text-muted">
					{t('level.label')}
				</p>
				<ChoiceChips
					labelledby="judul-jenjang-dasbor"
					items={pilihanJenjang(page.url, data.jenjang)}
				/>
			</div>
		{/if}
		{#if pilihanTahun && pilihanTahun.daftar.length > 1}
			<nav aria-label={t('dash.yearNav')}>
				<ul class="flex flex-wrap gap-2">
					{#each pilihanTahun.daftar as tahun (tahun)}
						<li>
							<a
								href={tautanTahun(page.url, tahun)}
								data-sveltekit-noscroll
								aria-current={tahun === pilihanTahun.aktif ? 'true' : undefined}
								class="chip num aria-[current=true]:border-primary aria-[current=true]:bg-secondary aria-[current=true]:text-secondary-on"
								>{tahun}</a
							>
						</li>
					{/each}
				</ul>
			</nav>
		{/if}
	</div>

	<div class="mt-6 lg:mt-8">
		{#if !hasil.ok}
			<ErrorState
				title={t('dash.failed')}
				description={hasil.galat.message}
				onretry={invalidateAll}
			/>
		{:else if hasil.data.lensa === 'layanan'}
			<LensaLayanan dasbor={hasil.data.dasbor} peta={hasil.data.peta} jenjang={data.jenjang} />
		{:else if hasil.data.lensa === 'uang'}
			<LensaUang dasbor={hasil.data.dasbor} peta={hasil.data.peta} />
		{:else}
			<LensaTanggungJawab dasbor={hasil.data.dasbor} />
		{/if}
	</div>
</section>
