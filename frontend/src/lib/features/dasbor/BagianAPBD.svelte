<script lang="ts">
	import { Landmark } from '@lucide/svelte';
	import { ambilBelanjaAPBD, tangkap } from '$lib/api/klien';
	import type { AkunAPBD, BelanjaAPBD } from '$lib/api/tipe';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import ErrorState from '$lib/components/ErrorState.svelte';
	import { formatAngka, formatBulan, formatRupiah } from '$lib/format';
	import { bahasaAktif, t } from '$lib/i18n';
	import KartuGrafik from './KartuGrafik.svelte';
	import { persenDari } from './dasbor';

	let { tahun }: { tahun: number } = $props();

	const kodeTidakDitemukan = 404;
	const persenPenuh = 100;

	let belanja = $state<BelanjaAPBD | null>(null);
	let status = $state<'memuat' | 'siap' | 'kosong' | 'gagal'>('memuat');

	async function muat(tahunDiminta: number) {
		status = 'memuat';
		const hasil = await tangkap(ambilBelanjaAPBD(fetch, tahunDiminta));
		belanja = hasil.ok ? hasil.data : null;
		if (hasil.ok) status = 'siap';
		else status = hasil.galat.kode === kodeTidakDitemukan ? 'kosong' : 'gagal';
	}

	$effect(() => {
		void muat(tahun);
	});

	const rupiah = (n: number) => formatRupiah(n, bahasaAktif);
	const serapan = (a: AkunAPBD) =>
		a.anggaran.nilai > 0 ? (a.realisasi.nilai / a.anggaran.nilai) * persenPenuh : 0;
	const teksSerapan = (a: AkunAPBD) => `${formatAngka(serapan(a), bahasaAktif, 1)}%`;
	const bulanRealisasi = (b: BelanjaAPBD) =>
		`${formatBulan(`${b.tahun}-${b.bulan_realisasi}`, bahasaAktif)} ${b.tahun}`;
</script>

{#if status === 'memuat'}
	<div class="skel h-64 rounded-lg" aria-busy="true"></div>
{:else if status === 'gagal'}
	<ErrorState title={t('apbd.failed')} onretry={() => muat(tahun)} />
{:else if !belanja}
	<EmptyState title={t('apbd.empty', { tahun })} />
{:else}
	<KartuGrafik
		judul={t('apbd.title', { tahun: belanja.tahun })}
		keterangan={t('apbd.lead', { bulan: bulanRealisasi(belanja) })}
		ikon={Landmark}
		angka={belanja.belanja.realisasi}
		kolom={[t('apbd.col.type'), t('apbd.col.budget'), t('apbd.col.spent'), t('apbd.col.share')]}
		baris={[belanja.belanja, ...belanja.jenis].map((a) => [
			a.nama,
			rupiah(a.anggaran.nilai),
			rupiah(a.realisasi.nilai),
			teksSerapan(a)
		])}
	>
		<p class="mb-4 text-small">
			{t('apbd.summary', {
				realisasi: rupiah(belanja.belanja.realisasi.nilai),
				anggaran: rupiah(belanja.belanja.anggaran.nilai),
				persen: teksSerapan(belanja.belanja)
			})}
		</p>
		<ol class="flex flex-col gap-4">
			{#each belanja.jenis as akun (akun.nama)}
				<li class="flex flex-col gap-1">
					<div class="flex items-baseline justify-between gap-3 text-small">
						<span class="min-w-0 break-words font-semibold">{akun.nama}</span>
						<span class="num shrink-0 font-bold">{teksSerapan(akun)}</span>
					</div>
					<div class="h-3 rounded-sm bg-sunken" aria-hidden="true">
						<div
							class="h-full rounded-sm bg-chart"
							style:width="{persenDari(akun.realisasi.nilai, akun.anggaran.nilai)}%"
						></div>
					</div>
					<p class="num text-caption text-muted">
						{t('apbd.item', {
							realisasi: rupiah(akun.realisasi.nilai),
							anggaran: rupiah(akun.anggaran.nilai)
						})}
					</p>
				</li>
			{/each}
		</ol>
	</KartuGrafik>
{/if}
