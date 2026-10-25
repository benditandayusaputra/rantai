<script lang="ts">
	import {
		Building2,
		CalendarDays,
		ChartBarBig,
		ChartColumn,
		HandCoins,
		Info,
		Map as IkonPeta,
		Scale
	} from '@lucide/svelte';
	import type { DasborUang, PetaTema, Porsi } from '$lib/api/tipe';
	import KeyNumberCard from '$lib/components/KeyNumberCard.svelte';
	import { formatAngka, formatBulan, formatRupiah } from '$lib/format';
	import { bahasaAktif, t } from '$lib/i18n';
	import GrafikBatang from './GrafikBatang.svelte';
	import GrafikKolom from './GrafikKolom.svelte';
	import KartuGrafik from './KartuGrafik.svelte';
	import PetaWilayah from './PetaWilayah.svelte';
	import {
		ambangKuantil,
		angkaDari,
		jumlahkan,
		namaWilayahPendek,
		porsiTahun,
		tanpaNilai,
		tingkatKuantil,
		type IsiArea
	} from './dasbor';

	let { dasbor, peta }: { dasbor: DasborUang; peta: PetaTema } = $props();

	type Ukuran = 'jumlah' | 'nilai';
	let ukuranBulanan = $state<Ukuran>('jumlah');

	const tingkatOpasitas = [0.2, 0.4, 0.6, 0.8, 1];
	const tahun = $derived(dasbor.tahun);
	const tahunIni = $derived(porsiTahun(dasbor.tahunan, tahun));
	const terpetakan = $derived(jumlahkan(dasbor.kecamatan.isi, (p) => p.jumlah));
	const bertanggal = $derived(jumlahkan(dasbor.bulanan.isi, (p) => p.jumlah));
	const provinsi = $derived(dasbor.kota.isi.find((p) => !p.kunci.includes('.')));
	const kota = $derived(dasbor.kota.isi.filter((p) => p.kunci.includes('.')));
	const perKecamatan = $derived(new Map(dasbor.kecamatan.isi.map((p) => [p.kunci, p])));
	const ambang = $derived(
		ambangKuantil(
			dasbor.kecamatan.isi.map((p) => p.nilai),
			tingkatOpasitas.length
		)
	);

	const legendaNilai = $derived(
		ambang.length === 0
			? []
			: [undefined, ...ambang].map((bawah, i) => ({
					opasitas: tingkatOpasitas[i] ?? 1,
					label: bawah === undefined ? `< ${rupiah(ambang[0] ?? 0)}` : `≥ ${rupiah(bawah)}`
				}))
	);

	const angka = (nilai: number, satuan: string) => angkaDari(dasbor.tahunan, nilai, satuan);
	const rupiah = (n: number) => formatRupiah(n, bahasaAktif);
	const bilangan = (n: number) => formatAngka(n, bahasaAktif);
	const paket = (p: Porsi) => t('dash.money.packages', { n: bilangan(p.jumlah) });
	const nilaiDari = (p: Porsi) => p.nilai;
	const jumlahDari = (p: Porsi) => p.jumlah;
	const barisPorsi = (isi: Porsi[], namaDari: (p: Porsi) => string) =>
		isi.map((p) => [namaDari(p), bilangan(p.jumlah), rupiah(p.nilai)]);
	const kolomPorsi = (pertama: string) => [pertama, t('dash.col.count'), t('dash.col.value')];
	const namaPorsi = (p: Porsi) => namaWilayahPendek(p.nama ?? p.kunci);

	function isiDari(kode: string): IsiArea {
		const p = perKecamatan.get(kode);
		const tingkat = tingkatKuantil(p?.nilai ?? 0, ambang);
		if (!p || tingkat === tanpaNilai) {
			return { kelas: 'fill-sunken', keterangan: t('dash.money.map.none') };
		}
		return {
			kelas: 'fill-primary',
			opasitas: tingkatOpasitas[tingkat],
			keterangan: `${rupiah(p.nilai)} · ${paket(p)}`
		};
	}

	const grafikTahunan = [
		{
			judul: t('dash.money.yearly.count'),
			nilai: jumlahDari,
			format: bilangan,
			satuan: 'paket',
			ikon: ChartColumn,
			rincian: (p: Porsi) => rupiah(p.nilai)
		},
		{
			judul: t('dash.money.yearly.value'),
			nilai: nilaiDari,
			format: rupiah,
			satuan: 'rupiah',
			ikon: HandCoins,
			rincian: paket
		}
	];
	const pilihanUkuran: { kunci: Ukuran; label: string }[] = [
		{ kunci: 'jumlah', label: t('dash.measure.count') },
		{ kunci: 'nilai', label: t('dash.measure.value') }
	];
	const grafikRincian = $derived([
		{ judul: t('dash.money.type.title', { tahun }), seri: dasbor.jenis, ikon: ChartBarBig },
		{ judul: t('dash.money.method.title', { tahun }), seri: dasbor.metode, ikon: Scale }
	]);

	const kartuAngka = $derived.by(() => {
		if (!tahunIni) return [];
		return [
			{
				label: t('dash.money.count'),
				angka: angka(tahunIni.jumlah, 'paket'),
				catatan: t('dash.money.budgetYear', { tahun })
			},
			{
				label: t('dash.money.value'),
				angka: angka(tahunIni.nilai, 'rupiah'),
				catatan: t('dash.money.budgetYear', { tahun })
			},
			{
				label: t('dash.money.located'),
				angka: angka(terpetakan, 'paket'),
				catatan: t('dash.money.locatedNote', {
					persen: `${formatAngka((terpetakan / Math.max(1, tahunIni.jumlah)) * 100, bahasaAktif, 1)}%`
				})
			}
		];
	});
</script>

<div class="flex flex-col gap-6">
	<div class="grid grid-cols-2 gap-3 lg:grid-cols-3">
		{#each kartuAngka as kartu (kartu.label)}
			{#if kartu.angka}
				<KeyNumberCard label={kartu.label} angka={kartu.angka} catatan={kartu.catatan} />
			{/if}
		{/each}
	</div>

	<div class="grid gap-6 lg:grid-cols-2">
		{#each grafikTahunan as grafik (grafik.judul)}
			<KartuGrafik
				judul={grafik.judul}
				keterangan={t('dash.money.yearly.lead')}
				ikon={grafik.ikon}
				angka={tahunIni ? angka(grafik.nilai(tahunIni), grafik.satuan) : null}
				kolom={kolomPorsi(t('dash.col.period'))}
				baris={barisPorsi(dasbor.tahunan.isi, (p) => p.kunci)}
			>
				<GrafikKolom
					isi={dasbor.tahunan.isi}
					nilaiDari={grafik.nilai}
					format={grafik.format}
					labelDari={(p) => p.kunci}
					rincian={grafik.rincian}
					sorot={String(tahun)}
				/>
			</KartuGrafik>
		{/each}
	</div>

	<KartuGrafik
		judul={t('dash.money.monthly.title', { tahun })}
		keterangan={t('dash.money.monthly.lead', {
			n: bilangan(bertanggal),
			total: bilangan(tahunIni?.jumlah ?? 0),
			tahun
		})}
		ikon={CalendarDays}
		angka={angka(bertanggal, 'paket')}
		kolom={kolomPorsi(t('dash.col.period'))}
		baris={bertanggal > 0
			? barisPorsi(dasbor.bulanan.isi, (p) => formatBulan(p.kunci, bahasaAktif))
			: []}
	>
		{#snippet aksi()}
			<div role="group" aria-label={t('dash.measure')} class="flex flex-wrap gap-2">
				{#each pilihanUkuran as pilihan (pilihan.kunci)}
					<button
						type="button"
						class="chip"
						aria-pressed={ukuranBulanan === pilihan.kunci}
						onclick={() => (ukuranBulanan = pilihan.kunci)}>{pilihan.label}</button
					>
				{/each}
			</div>
		{/snippet}
		<GrafikKolom
			isi={dasbor.bulanan.isi}
			nilaiDari={ukuranBulanan === 'jumlah' ? jumlahDari : nilaiDari}
			format={ukuranBulanan === 'jumlah' ? bilangan : rupiah}
			labelDari={(p) => formatBulan(p.kunci, bahasaAktif)}
			rincian={ukuranBulanan === 'jumlah' ? (p) => rupiah(p.nilai) : paket}
		/>
	</KartuGrafik>

	<div class="grid gap-6 lg:grid-cols-12">
		<KartuGrafik
			class="lg:col-span-7"
			judul={t('dash.money.map.title', { tahun })}
			keterangan={t('dash.money.map.lead')}
			ikon={IkonPeta}
			angka={angka(jumlahkan(dasbor.kecamatan.isi, nilaiDari), 'rupiah')}
			kolom={kolomPorsi(t('spend.col.area'))}
			baris={barisPorsi(dasbor.kecamatan.isi, namaPorsi)}
		>
			<PetaWilayah fitur={peta.features} {isiDari} judul={t('dash.money.map.title', { tahun })}>
				{#snippet legenda()}
					<p class="text-caption font-bold">{t('dash.col.value')}</p>
					<ul class="mt-2 flex flex-wrap gap-x-4 gap-y-2">
						<li class="flex items-center gap-2 text-caption">
							<span
								class="h-3 w-5 rounded-xs border-ctrl border-dashed border-line-strong bg-sunken"
								aria-hidden="true"
							></span>
							{t('dash.money.map.none')}
						</li>
						{#each legendaNilai as tingkat (tingkat.label)}
							<li class="flex items-center gap-2 text-caption">
								<span
									class="h-3 w-5 rounded-xs bg-primary"
									style:opacity={tingkat.opasitas}
									aria-hidden="true"
								></span>
								<span class="num">{tingkat.label}</span>
							</li>
						{/each}
					</ul>
				{/snippet}
			</PetaWilayah>
		</KartuGrafik>

		<KartuGrafik
			class="lg:col-span-5"
			judul={t('dash.money.city.title', { tahun })}
			ikon={Building2}
			angka={angka(jumlahkan(kota, nilaiDari), 'rupiah')}
			kolom={kolomPorsi(t('dash.col.city'))}
			baris={barisPorsi(kota, namaPorsi)}
		>
			<GrafikBatang isi={kota} {nilaiDari} format={rupiah} namaDari={namaPorsi} rincian={paket} />
			{#snippet catatan()}
				{#if provinsi}
					<p class="text-small text-muted">
						{t('dash.money.city.province', {
							n: bilangan(provinsi.jumlah),
							nilai: rupiah(provinsi.nilai)
						})}
					</p>
				{/if}
			{/snippet}
		</KartuGrafik>
	</div>

	<div class="grid gap-6 lg:grid-cols-2">
		{#each grafikRincian as grafik (grafik.judul)}
			<KartuGrafik
				judul={grafik.judul}
				ikon={grafik.ikon}
				angka={angka(jumlahkan(grafik.seri.isi, nilaiDari), 'rupiah')}
				kolom={kolomPorsi(t('dash.col.name'))}
				baris={barisPorsi(grafik.seri.isi, (p) => p.kunci)}
			>
				<GrafikBatang
					isi={grafik.seri.isi}
					{nilaiDari}
					format={rupiah}
					namaDari={(p) => p.kunci}
					rincian={paket}
				/>
			</KartuGrafik>
		{/each}
	</div>

	<p class="flex gap-2 rounded-md bg-info-bg p-3 text-small">
		<Info size={18} strokeWidth={1.75} class="mt-0.5 shrink-0 text-info" aria-hidden="true" />
		{t('dash.money.hpsNote')}
	</p>
</div>
