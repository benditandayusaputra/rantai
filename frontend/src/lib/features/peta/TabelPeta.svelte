<script lang="ts">
	import type { PetaPaket, PetaTema, SekolahTerdekat } from '$lib/api/tipe';
	import CertaintyBadge from '$lib/components/CertaintyBadge.svelte';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import GapClassChip from '$lib/components/GapClassChip.svelte';
	import { tampilkanAngka } from '$lib/features/kecamatan/kecamatan';
	import { formatJarak } from '$lib/format';
	import { bahasaBawaan, t } from '$lib/i18n';
	import type { StatusLapisan } from './lapisan';

	interface Props {
		lapisan: StatusLapisan;
		tema: PetaTema | null;
		paket: PetaPaket['features'];
		sekolah: SekolahTerdekat[];
		terpilih: string | null;
	}

	let { lapisan, tema, paket, sekolah, terpilih }: Props = $props();

	const kelasTabel = 'w-full min-w-max text-small';
	const kelasKepala = 'bg-sunken text-caption uppercase';
	const angkaSel = (a: Parameters<typeof tampilkanAngka>[0] | null) =>
		a ? tampilkanAngka(a).nilai : '';
</script>

<div class="flex flex-col gap-8">
	{#if lapisan.tema}
		<section aria-labelledby="tabel-tema">
			<h2 id="tabel-tema" class="t-h3">{t('map.table.kecamatan')}</h2>
			{#if tema?.features.length}
				<div class="mt-3 overflow-x-auto rounded-lg border border-line bg-surface">
					<table class={kelasTabel}>
						<caption class="sr-only">{t('map.table.kecamatan')}</caption>
						<thead class={kelasKepala}>
							<tr>
								<th scope="col" class="p-3 text-left">{t('map.col.name')}</th>
								<th scope="col" class="p-3 text-right">{t('tema.needs')}</th>
								<th scope="col" class="p-3 text-right">{t('tema.capacity')}</th>
								<th scope="col" class="p-3 text-right">{t('tema.ratio')}</th>
								<th scope="col" class="p-3 text-left">{t('map.col.class')}</th>
							</tr>
						</thead>
						<tbody>
							{#each tema.features as fitur (fitur.properties.kode)}
								{@const p = fitur.properties}
								<tr class="border-b border-line {terpilih === p.kode ? 'bg-secondary' : ''}">
									<th scope="row" class="p-3 text-left">
										<a href="/{bahasaBawaan}/kecamatan/{p.kode}" class="font-semibold underline"
											>{p.nama}</a
										>
									</th>
									<td class="num p-3 text-right">{angkaSel(p.kebutuhan)}</td>
									<td class="num p-3 text-right">{angkaSel(p.kapasitas)}</td>
									<td class="num p-3 text-right">{angkaSel(p.rasio)}</td>
									<td class="p-3"><GapClassChip kelas={p.kelas} /></td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			{:else}
				<div class="mt-3"><EmptyState title={t('map.empty')} /></div>
			{/if}
		</section>
	{/if}
	{#if lapisan.proyek}
		<section aria-labelledby="tabel-proyek">
			<h2 id="tabel-proyek" class="t-h3">{t('map.table.proyek')}</h2>
			{#if paket.length}
				<div class="mt-3 overflow-x-auto rounded-lg border border-line bg-surface">
					<table class={kelasTabel}>
						<caption class="sr-only">{t('map.table.proyek')}</caption>
						<thead class={kelasKepala}>
							<tr>
								<th scope="col" class="p-3 text-left">{t('map.col.name')}</th>
								<th scope="col" class="p-3 text-left">{t('map.col.location')}</th>
							</tr>
						</thead>
						<tbody>
							{#each paket as fitur, urutan (urutan)}
								{@const p = fitur.properties}
								<tr class="border-b border-line">
									{#if p.jenis_lokasi === 'titik'}
										<th scope="row" class="p-3 text-left">
											<a href="/{bahasaBawaan}/paket/{p.kode}" class="font-semibold underline"
												>{p.nama}</a
											>
											{#if p.is_simulation}<CertaintyBadge level="simulasi" class="ml-2" />{/if}
										</th>
										<td class="p-3">{t('map.location.titik')}</td>
									{:else}
										<th scope="row" class="p-3 text-left font-semibold">{p.nama_kecamatan}</th>
										<td class="p-3">{t('map.areaPackages', { n: p.paket.length })}</td>
									{/if}
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			{:else}
				<div class="mt-3"><EmptyState title={t('map.empty')} /></div>
			{/if}
		</section>
	{/if}
	{#if lapisan.sekolah}
		<section aria-labelledby="tabel-sekolah">
			<h2 id="tabel-sekolah" class="t-h3">{t('map.table.sekolah')}</h2>
			{#if sekolah.length}
				<div class="mt-3 overflow-x-auto rounded-lg border border-line bg-surface">
					<table class={kelasTabel}>
						<caption class="sr-only">{t('map.table.sekolah')}</caption>
						<thead class={kelasKepala}>
							<tr>
								<th scope="col" class="p-3 text-left">{t('map.col.name')}</th>
								<th scope="col" class="p-3 text-right">{t('pen.school.seats')}</th>
								<th scope="col" class="p-3 text-right">{t('map.col.distance')}</th>
							</tr>
						</thead>
						<tbody>
							{#each sekolah as item (item.npsn)}
								<tr class="border-b border-line">
									<th scope="row" class="p-3 text-left font-semibold">{item.nama}</th>
									<td class="num p-3 text-right">{angkaSel(item.daya_tampung)}</td>
									<td class="num p-3 text-right">{formatJarak(item.jarak_m, bahasaBawaan)}</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			{:else}
				<div class="mt-3"><EmptyState title={t('map.empty')} /></div>
			{/if}
		</section>
	{/if}
</div>
