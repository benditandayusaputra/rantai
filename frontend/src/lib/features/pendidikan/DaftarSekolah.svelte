<script lang="ts">
	import { MapPin } from '@lucide/svelte';
	import TableScroll from '$lib/components/TableScroll.svelte';
	import type { Sekolah } from '$lib/api/tipe';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import FotoSekolah from '$lib/components/FotoSekolah.svelte';
	import KreditFoto from '$lib/components/KreditFoto.svelte';
	import ProvenanceBadge from '$lib/components/ProvenanceBadge.svelte';
	import { tampilkanAngka } from '$lib/features/kecamatan/kecamatan';
	import { bahasaAktif, t } from '$lib/i18n';
	import { teksJenjang, temaPendidikanBawaan, type TemaPendidikan } from './jenjang';
	import { kunciJenisSekolah } from './pendidikan';

	let {
		sekolah,
		tema = temaPendidikanBawaan,
		kodeKecamatan
	}: { sekolah: Sekolah[]; tema?: TemaPendidikan; kodeKecamatan?: string } = $props();

	const teks = $derived(teksJenjang[tema]);
</script>

{#if sekolah.length}
	<TableScroll label={t(teks.judulSekolah)}>
		<table class="w-full min-w-max text-small">
			<caption class="sr-only">{t(teks.judulSekolah)}</caption>
			<thead class="bg-sunken text-caption font-semibold">
				<tr>
					<th scope="col" class="sticky left-0 bg-sunken p-3 text-left">{t('pen.school.name')}</th>
					<th scope="col" class="p-3 text-left">{t('pen.school.type')}</th>
					<th scope="col" class="p-3 text-right">{t(teks.kolomKursi)}</th>
					<th scope="col" class="p-3 text-left"><span class="sr-only">{t('nav.map')}</span></th>
				</tr>
			</thead>
			<tbody>
				{#each sekolah as item (item.npsn)}
					{@const tampil = item.daya_tampung ? tampilkanAngka(item.daya_tampung) : null}
					<tr class="border-b border-line last:border-b-0 hover:bg-sunken">
						<th scope="row" class="sticky left-0 bg-surface p-3 text-left font-semibold">
							<span class="flex items-center gap-3">
								<FotoSekolah foto={item.foto} />
								<span class="flex flex-col gap-0.5">
									{item.nama}
									{#if item.foto}<KreditFoto foto={item.foto} />{/if}
								</span>
							</span>
						</th>
						<td class="p-3">{t(kunciJenisSekolah(item))}</td>
						<td class="p-3 text-right">
							{#if item.daya_tampung && tampil}
								<span class="inline-flex items-center gap-2">
									<span class="num font-bold">{tampil.nilai}</span>
									<ProvenanceBadge angka={item.daya_tampung} nilai={tampil.nilai} compact />
								</span>
							{:else}
								<span class="text-muted">{t('pen.school.noSeats')}</span>
							{/if}
						</td>
						<td class="p-3">
							{#if item.lokasi}
								<a
									href={`/${bahasaAktif}/peta?lat=${item.lokasi.lat}&lng=${item.lokasi.lng}&tema=${tema}${kodeKecamatan ? `&w=${encodeURIComponent(kodeKecamatan)}` : ''}`}
									class="link inline-flex min-h-touch items-center gap-1 whitespace-nowrap"
								>
									<MapPin size={16} strokeWidth={1.75} aria-hidden="true" />{t('pen.school.onMap')}
								</a>
							{/if}
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</TableScroll>
{:else}
	<EmptyState title={t(teks.sekolahKosong)} />
{/if}
