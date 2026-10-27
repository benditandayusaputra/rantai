<script lang="ts">
	import { Map, Table2 } from '@lucide/svelte';
	import { goto } from '$app/navigation';
	import { ambilPetaTema, ambilSekolahKecamatan, tangkap } from '$lib/api/klien';
	import type { PetaTema, Sekolah } from '$lib/api/tipe';
	import PetaBerlapis from '$lib/features/peta/PetaBerlapis.svelte';
	import LegendaKelas from '$lib/features/peta/LegendaKelas.svelte';
	import { bacaLapisan, bacaPetaDasar } from '$lib/features/peta/lapisan';
	import DaftarSekolah from '$lib/features/pendidikan/DaftarSekolah.svelte';
	import { teksJenjang, type TemaPendidikan } from '$lib/features/pendidikan/jenjang';
	import { bahasaAktif, t } from '$lib/i18n';

	let { kodeKecamatan, kodeTema }: { kodeKecamatan: string; kodeTema: TemaPendidikan } = $props();
	let isTabel = $state(false);
	let peta = $state<PetaTema | null>(null);
	let sekolah = $state<Sekolah[] | null>(null);
	let galatPeta = $state(false);
	let galatSekolah = $state(false);
	let lapisan = $state(bacaLapisan(null));
	let petaDasar = $state(bacaPetaDasar(null));

	$effect(() => {
		let dibatalkan = false;
		peta = null;
		sekolah = null;
		galatPeta = false;
		galatSekolah = false;
		void Promise.all([
			tangkap(ambilPetaTema(fetch, kodeTema)),
			tangkap(ambilSekolahKecamatan(fetch, kodeTema, kodeKecamatan))
		]).then(([hasilPeta, hasilSekolah]) => {
			if (dibatalkan) return;
			peta = hasilPeta.ok ? hasilPeta.data : null;
			sekolah = hasilSekolah.ok ? hasilSekolah.data : null;
			galatPeta = !hasilPeta.ok;
			galatSekolah = !hasilSekolah.ok;
		});
		return () => {
			dibatalkan = true;
		};
	});
</script>

<section aria-labelledby="judul-data-layanan" class="flex flex-col gap-4">
	<div class="flex flex-wrap items-center justify-between gap-3">
		<h3 id="judul-data-layanan" class="t-h3">{t('services.exploreData')}</h3>
		<div role="group" aria-label={t('services.dataView')} class="flex gap-2">
			<button
				type="button"
				class="btn {isTabel ? 'btn-outline' : 'btn-secondary'}"
				aria-pressed={!isTabel}
				onclick={() => (isTabel = false)}
				><Map size={18} aria-hidden="true" />{t('services.mapView')}</button
			>
			<button
				type="button"
				class="btn {isTabel ? 'btn-secondary' : 'btn-outline'}"
				aria-pressed={isTabel}
				onclick={() => (isTabel = true)}
				><Table2 size={18} aria-hidden="true" />{t('services.tableView')}</button
			>
		</div>
	</div>
	{#if isTabel}
		{#if galatSekolah}<p class="text-small text-muted">{t('services.tableFailed')}</p>
		{:else if sekolah}<DaftarSekolah {sekolah} tema={kodeTema} {kodeKecamatan} />
		{:else}<p aria-live="polite" class="text-small text-muted">{t('state.loading')}</p>{/if}
	{:else}
		<p class="text-small text-muted">{t('services.mapScope')}</p>
		{#if galatPeta}<p class="text-small text-muted">{t('home.map.error')}</p>
		{:else if peta}
			<PetaBerlapis
				class="h-peta-kartu"
				tema={peta}
				temaSekolah={kodeTema}
				terpilih={kodeKecamatan}
				onpilih={(kode) => goto(`/${bahasaAktif}/kecamatan/${kode}?bagian=tema&tema=${kodeTema}`)}
				bind:lapisan
				bind:petaDasar
				nama="peta-layanan"
				judulLegenda={teksJenjang[kodeTema].legenda}
			/>
			<LegendaKelas judul={teksJenjang[kodeTema].legenda} mendatar />
		{:else}<div class="skel h-64 rounded-lg" aria-busy="true">
				<p class="sr-only">{t('map.loading')}</p>
			</div>{/if}
	{/if}
</section>
