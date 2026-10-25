<script lang="ts">
	import {
		ChartColumnStacked,
		ChartNoAxesCombined,
		Info,
		Map as IkonPeta,
		Users
	} from '@lucide/svelte';
	import type { DasborLayanan, PetaTema, RekapKursi } from '$lib/api/tipe';
	import KeyNumberCard from '$lib/components/KeyNumberCard.svelte';
	import { kunciKelas, tampilkanAngka } from '$lib/features/kecamatan/kecamatan';
	import { formatAngka } from '$lib/format';
	import {
		teksJenjang,
		temaPendidikanBawaan,
		type TemaPendidikan
	} from '$lib/features/pendidikan/jenjang';
	import { bahasaAktif, t } from '$lib/i18n';
	import BatangKelas from './BatangKelas.svelte';
	import BatangPasangan from './BatangPasangan.svelte';
	import KartuGrafik from './KartuGrafik.svelte';
	import PetaWilayah from './PetaWilayah.svelte';
	import {
		angkaDari,
		hitungKelas,
		isiKelasPeta,
		kelasCukup,
		kelasKurang,
		kotakKelasPeta,
		namaWilayahPendek,
		totalKecamatan,
		urutanKelas,
		type IsiArea
	} from './dasbor';

	interface Props {
		dasbor: DasborLayanan;
		peta: PetaTema;
		jenjang?: TemaPendidikan;
	}

	let { dasbor, peta, jenjang = temaPendidikanBawaan }: Props = $props();

	const teks = $derived(teksJenjang[jenjang]);

	const terbaru = $derived(dasbor.periode.at(-1));
	const properti = $derived(new Map(peta.features.map((f) => [f.properties.kode, f.properties])));
	const kecamatan = $derived(
		[...properti.values()].sort((a, b) => a.nama.localeCompare(b.nama, 'id'))
	);
	const angka = (nilai: number, satuan: string) => angkaDari(dasbor, nilai, satuan);
	const angkaBulat = (n: number) => formatAngka(n, bahasaAktif);
	const namaKota = (r: RekapKursi) => namaWilayahPendek(r.nama ?? r.kunci);

	function rasioTampil(kode: string): string {
		const rasio = properti.get(kode)?.rasio;
		return rasio ? tampilkanAngka(rasio).nilai : '-';
	}

	function isiDari(kode: string): IsiArea {
		const kelas = properti.get(kode)?.kelas ?? 'tidak_ada_data';
		return {
			kelas: isiKelasPeta[kelas],
			keterangan: `${t(kunciKelas[kelas])} · ${rasioTampil(kode)} ${t('satuan.rasio')}`
		};
	}

	const kartuAngka = $derived.by(() => {
		if (!terbaru) return [];
		const ringkas = t('dash.svc.ofTotal', { n: totalKecamatan(terbaru), periode: terbaru.kunci });
		return [
			{
				label: t(teks.dasborLulusan),
				angka: angka(terbaru.kebutuhan, 'siswa'),
				catatan: terbaru.kunci
			},
			{
				label: t(teks.dasborKursi),
				angka: angka(terbaru.kapasitas, 'kursi'),
				catatan: terbaru.kunci
			},
			{
				label: t('dash.svc.short'),
				angka: angka(hitungKelas(terbaru, kelasKurang), 'kecamatan'),
				catatan: ringkas
			},
			{
				label: t('dash.svc.enough'),
				angka: angka(hitungKelas(terbaru, kelasCukup), 'kecamatan'),
				catatan: ringkas
			}
		];
	});

	const kolomKelas = urutanKelas.map((k) => t(kunciKelas[k]));
	const barisKelas = (r: RekapKursi, nama: string) => [
		nama,
		...urutanKelas.map((k) => angkaBulat(r.per_kelas[k] ?? 0))
	];
</script>

{#if terbaru}
	<div class="flex flex-col gap-6">
		<div class="grid grid-cols-2 gap-3 lg:grid-cols-4">
			{#each kartuAngka as kartu (kartu.label)}
				{#if kartu.angka}
					<KeyNumberCard label={kartu.label} angka={kartu.angka} catatan={kartu.catatan} />
				{/if}
			{/each}
		</div>

		<div class="grid gap-6 lg:grid-cols-12">
			<KartuGrafik
				class="lg:col-span-7"
				judul={t('dash.svc.map.title')}
				keterangan={t('dash.svc.map.lead', { periode: terbaru.kunci })}
				ikon={IkonPeta}
				angka={angka(terbaru.kebutuhan, 'siswa')}
				kolom={[t('spend.col.area'), t('map.col.class'), t('spend.col.ratio')]}
				baris={kecamatan.map((k) => [k.nama, t(kunciKelas[k.kelas]), rasioTampil(k.kode)])}
			>
				<PetaWilayah fitur={peta.features} {isiDari} judul={t('dash.svc.map.title')}>
					{#snippet legenda()}
						<p class="text-caption font-bold">{t(teks.legenda)}</p>
						<ul class="mt-2 flex flex-wrap gap-x-4 gap-y-2">
							{#each [...urutanKelas, 'tidak_ada_data' as const] as kelas (kelas)}
								<li class="flex items-center gap-2 text-caption">
									<span class="h-3 w-5 rounded-xs {kotakKelasPeta[kelas]}" aria-hidden="true"
									></span>
									{t(kunciKelas[kelas])}
								</li>
							{/each}
						</ul>
					{/snippet}
				</PetaWilayah>
			</KartuGrafik>

			<KartuGrafik
				class="lg:col-span-5"
				judul={t('dash.svc.city.title')}
				keterangan={t('dash.svc.city.lead', { periode: terbaru.kunci })}
				ikon={ChartColumnStacked}
				angka={angka(totalKecamatan(terbaru), 'kecamatan')}
				kolom={[t('dash.col.city'), ...kolomKelas]}
				baris={dasbor.kota.map((r) => barisKelas(r, namaKota(r)))}
			>
				<BatangKelas daftar={dasbor.kota} namaDari={namaKota} />
			</KartuGrafik>
		</div>

		<div class="grid gap-6 lg:grid-cols-2">
			<KartuGrafik
				judul={t('dash.svc.pair.title')}
				keterangan={t('dash.svc.pair.lead', { periode: terbaru.kunci })}
				ikon={Users}
				angka={angka(terbaru.kebutuhan, 'siswa')}
				kolom={[t('dash.col.city'), t(teks.dasborLulusan), t(teks.dasborKursi)]}
				baris={dasbor.kota.map((r) => [
					namaKota(r),
					angkaBulat(r.kebutuhan),
					angkaBulat(r.kapasitas)
				])}
			>
				<BatangPasangan daftar={dasbor.kota} namaDari={namaKota} />
			</KartuGrafik>

			<KartuGrafik
				judul={t('dash.svc.trend.title')}
				keterangan={t('dash.svc.trend.lead')}
				ikon={ChartNoAxesCombined}
				angka={angka(terbaru.kebutuhan, 'siswa')}
				kolom={[t('pen.trend.period'), t(teks.dasborLulusan), t(teks.dasborKursi), ...kolomKelas]}
				baris={dasbor.periode.map((r) => [
					r.kunci,
					angkaBulat(r.kebutuhan),
					angkaBulat(r.kapasitas),
					...barisKelas(r, r.kunci).slice(1)
				])}
			>
				<div class="flex flex-col gap-8">
					<BatangPasangan daftar={dasbor.periode} namaDari={(r) => r.kunci} />
					<div>
						<p class="mb-3 text-small font-bold">{t('dash.svc.dist.title')}</p>
						<BatangKelas daftar={dasbor.periode} namaDari={(r) => r.kunci} />
					</div>
				</div>
				{#snippet catatan()}
					{#if dasbor.periode.length === 1}
						<p class="flex gap-2 rounded-md bg-info-bg p-3 text-small">
							<Info
								size={18}
								strokeWidth={1.75}
								class="mt-0.5 shrink-0 text-info"
								aria-hidden="true"
							/>
							{t('dash.svc.trend.single')}
						</p>
					{/if}
				{/snippet}
			</KartuGrafik>
		</div>
	</div>
{:else}
	<p class="card p-6 text-center text-muted">{t('dash.empty')}</p>
{/if}
