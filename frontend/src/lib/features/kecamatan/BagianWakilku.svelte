<script lang="ts">
	import type { WakilKecamatan } from '$lib/api/tipe';
	import KeyNumberCard from '$lib/components/KeyNumberCard.svelte';
	import RepresentativeCard from '$lib/components/RepresentativeCard.svelte';
	import AspirasiReses from './AspirasiReses.svelte';
	import { ArrowRight, FileLock, FileText, Landmark, Users } from '@lucide/svelte';
	import { bahasaAktif, t } from '$lib/i18n';

	let { wakilku }: { wakilku: WakilKecamatan } = $props();

	const dapil = $derived(wakilku.dapil);
</script>

<section aria-labelledby="judul-dapil" class="flex flex-col gap-6">
	<div class="card flex flex-col gap-5 p-4 md:p-6">
		<div class="flex items-center gap-3">
			<span class="icon-tile"><Landmark size={20} strokeWidth={1.75} aria-hidden="true" /></span>
			<h3 id="judul-dapil" class="t-h3">{dapil.nama}</h3>
		</div>
		<div class="grid grid-cols-2 gap-3 lg:gap-4">
			<KeyNumberCard label={t('reps.seats')} angka={dapil.jumlah_kursi} />
			<KeyNumberCard label={t('reps.count')} angka={dapil.jumlah_anggota} />
		</div>
		<p class="text-small text-muted">
			{t('reps.noScore')}
			<a
				href="/{bahasaAktif}/data-belum-dibuka"
				class="link inline-flex items-center gap-1 align-bottom"
			>
				<FileLock size={14} strokeWidth={1.75} aria-hidden="true" />{t('unopened.link')}</a
			>
		</p>
		<a
			href="/{bahasaAktif}/aspirasi?w={encodeURIComponent(wakilku.kecamatan.kode)}"
			class="btn btn-primary self-start max-md:btn-block"
		>
			<FileText size={18} strokeWidth={1.75} aria-hidden="true" />
			{t('acc.createPacket')}
			<ArrowRight size={18} strokeWidth={1.75} aria-hidden="true" />
		</a>
	</div>
	<div class="flex flex-col gap-3">
		<h3 class="flex items-center gap-2 text-h3 font-bold">
			<Users size={20} strokeWidth={1.75} aria-hidden="true" class="text-primary" />
			{t('acc.reps', { dapil: dapil.nomor })}
		</h3>
		<ul class="grid gap-3 md:grid-cols-2 md:gap-4 xl:grid-cols-3">
			{#each dapil.anggota as wakil (wakil.id)}
				<li><RepresentativeCard {wakil} /></li>
			{/each}
		</ul>
	</div>
	<AspirasiReses idDapil={dapil.id} />
</section>
