<script lang="ts">
	import TableScroll from '$lib/components/TableScroll.svelte';
	import type { Sekolah } from '$lib/api/tipe';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import ProvenanceBadge from '$lib/components/ProvenanceBadge.svelte';
	import { tampilkanAngka } from '$lib/features/kecamatan/kecamatan';
	import { bahasaBawaan, t } from '$lib/i18n';
	import { kunciJenisSekolah } from './pendidikan';

	let { sekolah }: { sekolah: Sekolah[] } = $props();
</script>

{#if sekolah.length}
	<TableScroll label={t('pen.schoolsHeading')}>
		<table class="w-full min-w-max text-small">
			<caption class="sr-only">{t('pen.schoolsHeading')}</caption>
			<thead class="bg-sunken text-caption uppercase">
				<tr>
					<th scope="col" class="sticky left-0 bg-sunken p-3 text-left">{t('pen.school.name')}</th>
					<th scope="col" class="p-3 text-left">{t('pen.school.type')}</th>
					<th scope="col" class="p-3 text-right">{t('pen.school.seats')}</th>
					<th scope="col" class="p-3 text-left"><span class="sr-only">{t('nav.map')}</span></th>
				</tr>
			</thead>
			<tbody>
				{#each sekolah as item (item.npsn)}
					{@const tampil = item.daya_tampung ? tampilkanAngka(item.daya_tampung) : null}
					<tr class="border-b border-line">
						<th scope="row" class="sticky left-0 bg-surface p-3 text-left font-semibold"
							>{item.nama}</th
						>
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
									href="/{bahasaBawaan}/peta?lat={item.lokasi.lat}&lng={item.lokasi.lng}"
									class="inline-flex min-h-touch items-center font-semibold text-secondary-on underline"
									>{t('pen.school.onMap')}</a
								>
							{/if}
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</TableScroll>
{:else}
	<EmptyState title={t('pen.schoolsEmpty')} />
{/if}
