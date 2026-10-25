<script lang="ts">
	import { CalendarDays, ChartColumn, Info, Landmark, Megaphone, Tags } from '@lucide/svelte';
	import type { DasborTanggungJawab, Porsi, Seri } from '$lib/api/tipe';
	import KeyNumberCard from '$lib/components/KeyNumberCard.svelte';
	import { formatAngka, formatBulan, formatRupiah } from '$lib/format';
	import { bahasaAktif, t } from '$lib/i18n';
	import GrafikBatang from './GrafikBatang.svelte';
	import GrafikKolom from './GrafikKolom.svelte';
	import KartuGrafik from './KartuGrafik.svelte';
	import { angkaDari, hurufAwalBesar, jumlahkan, porsiTahun, rataPerBulanTercatat } from './dasbor';

	let { dasbor }: { dasbor: DasborTanggungJawab } = $props();

	const tahun = $derived(dasbor.tahun);
	const kegiatanIni = $derived(porsiTahun(dasbor.kegiatan_tahunan, tahun)?.jumlah ?? 0);
	const bilangan = (n: number) => formatAngka(n, bahasaAktif);
	const rupiah = (n: number) => formatRupiah(n, bahasaAktif);
	const jumlahDari = (p: Porsi) => p.jumlah;
	const nilaiDari = (p: Porsi) => p.nilai;
	const total = (seri: Seri) => jumlahkan(seri.isi, jumlahDari);
	const barisJumlah = (isi: Porsi[], namaDari: (p: Porsi) => string) =>
		isi.map((p) => [namaDari(p), bilangan(p.jumlah)]);

	const kartuAngka = $derived([
		{
			label: t('dash.acc.activity'),
			angka: angkaDari(dasbor.kegiatan_tahunan, kegiatanIni, 'kegiatan'),
			catatan: t('dash.acc.perMonth', {
				n: bilangan(Math.round(rataPerBulanTercatat(dasbor.kegiatan_bulanan.isi)))
			})
		},
		{
			label: t('dash.acc.issues'),
			angka: angkaDari(dasbor.topik, total(dasbor.topik), 'sebutan'),
			catatan: String(tahun)
		},
		{
			label: t('dash.acc.reses'),
			angka: angkaDari(dasbor.reses, total(dasbor.reses), 'aspirasi'),
			catatan: String(tahun)
		}
	]);
	const grafikIsu = $derived([
		{
			judul: t('dash.acc.topic.title', { tahun }),
			keterangan: t('dash.acc.topic.lead'),
			seri: dasbor.topik,
			ikon: Tags,
			namaDari: (p: Porsi) => p.kunci
		},
		{
			judul: t('dash.acc.reses.title', { tahun }),
			keterangan: t('dash.acc.reses.lead'),
			seri: dasbor.reses,
			ikon: Megaphone,
			namaDari: (p: Porsi) => hurufAwalBesar(p.kunci)
		}
	]);
</script>

<div class="flex flex-col gap-6">
	<div class="grid grid-cols-2 gap-3 lg:grid-cols-3">
		{#each kartuAngka as kartu (kartu.label)}
			{#if kartu.angka}
				<KeyNumberCard label={kartu.label} angka={kartu.angka} catatan={kartu.catatan} />
			{/if}
		{/each}
	</div>

	<div class="grid gap-6 lg:grid-cols-12">
		<KartuGrafik
			class="lg:col-span-8"
			judul={t('dash.acc.monthly.title', { tahun })}
			keterangan={t('dash.acc.monthly.lead')}
			ikon={CalendarDays}
			angka={angkaDari(dasbor.kegiatan_bulanan, kegiatanIni, 'kegiatan')}
			kolom={[t('dash.col.period'), t('dash.col.count')]}
			baris={kegiatanIni > 0
				? barisJumlah(dasbor.kegiatan_bulanan.isi, (p) => formatBulan(p.kunci, bahasaAktif))
				: []}
		>
			<GrafikKolom
				isi={dasbor.kegiatan_bulanan.isi}
				nilaiDari={jumlahDari}
				format={bilangan}
				labelDari={(p) => formatBulan(p.kunci, bahasaAktif)}
			/>
		</KartuGrafik>

		<KartuGrafik
			class="lg:col-span-4"
			judul={t('dash.acc.yearly.title')}
			keterangan={t('dash.acc.yearly.lead')}
			ikon={ChartColumn}
			angka={angkaDari(dasbor.kegiatan_tahunan, total(dasbor.kegiatan_tahunan), 'kegiatan')}
			kolom={[t('dash.col.period'), t('dash.col.count')]}
			baris={barisJumlah(dasbor.kegiatan_tahunan.isi, (p) => p.kunci)}
		>
			<GrafikKolom
				isi={dasbor.kegiatan_tahunan.isi}
				nilaiDari={jumlahDari}
				format={bilangan}
				labelDari={(p) => p.kunci}
				sorot={String(tahun)}
			/>
		</KartuGrafik>
	</div>

	<div class="grid gap-6 lg:grid-cols-2">
		{#each grafikIsu as grafik (grafik.judul)}
			<KartuGrafik
				judul={grafik.judul}
				keterangan={grafik.keterangan}
				ikon={grafik.ikon}
				angka={angkaDari(grafik.seri, total(grafik.seri), 'kegiatan')}
				kolom={[t('dash.col.name'), t('dash.col.count')]}
				baris={barisJumlah(grafik.seri.isi, grafik.namaDari)}
			>
				<GrafikBatang
					isi={grafik.seri.isi}
					nilaiDari={jumlahDari}
					format={bilangan}
					namaDari={grafik.namaDari}
				/>
			</KartuGrafik>
		{/each}
	</div>

	<KartuGrafik
		judul={t('dash.acc.agency.title', { tahun })}
		keterangan={t('dash.acc.agency.lead')}
		ikon={Landmark}
		angka={angkaDari(dasbor.instansi, jumlahkan(dasbor.instansi.isi, nilaiDari), 'rupiah')}
		kolom={[t('dash.col.name'), t('dash.col.count'), t('dash.col.value')]}
		baris={dasbor.instansi.isi.map((p) => [p.nama ?? p.kunci, bilangan(p.jumlah), rupiah(p.nilai)])}
	>
		<GrafikBatang
			isi={dasbor.instansi.isi}
			{nilaiDari}
			format={rupiah}
			namaDari={(p) => p.nama ?? p.kunci}
			rincian={(p) => t('dash.money.packages', { n: bilangan(p.jumlah) })}
		/>
		{#snippet catatan()}
			<p class="text-small text-muted">{t('dash.acc.agency.note')}</p>
		{/snippet}
	</KartuGrafik>

	<p class="flex gap-2 rounded-md bg-info-bg p-3 text-small">
		<Info size={18} strokeWidth={1.75} class="mt-0.5 shrink-0 text-info" aria-hidden="true" />
		{t('dash.acc.noScore')}
	</p>
</div>
