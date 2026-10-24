<script lang="ts">
	import type { WakilKecamatan } from '$lib/api/tipe';
	import KeyNumberCard from '$lib/components/KeyNumberCard.svelte';
	import RepresentativeCard from '$lib/components/RepresentativeCard.svelte';
	import { FileText } from '@lucide/svelte';
	import { bahasaBawaan, t } from '$lib/i18n';

	let { wakilku }: { wakilku: WakilKecamatan } = $props();

	const dapil = $derived(wakilku.dapil);
</script>

<section aria-labelledby="judul-dapil" class="flex flex-col gap-4">
	<h3 id="judul-dapil" class="t-h3">{dapil.nama}</h3>
	<div class="grid grid-cols-2 gap-3 lg:gap-5">
		<KeyNumberCard label={t('reps.seats')} angka={dapil.jumlah_kursi} />
		<KeyNumberCard label={t('reps.count')} angka={dapil.jumlah_anggota} />
	</div>
	<p class="text-small text-muted">{t('reps.noScore')}</p>
	<a
		href="/{bahasaBawaan}/aspirasi?w={encodeURIComponent(wakilku.kecamatan.kode)}"
		class="inline-flex min-h-touch items-center justify-center gap-2 self-start rounded-md bg-primary px-4 font-semibold text-primary-on max-md:w-full"
	>
		<FileText size={20} strokeWidth={1.75} aria-hidden="true" />{t('acc.createPacket')}
	</a>
	<ul class="grid gap-3 md:grid-cols-2 md:gap-4">
		{#each dapil.anggota as wakil (wakil.id)}
			<li><RepresentativeCard {wakil} /></li>
		{/each}
	</ul>
</section>
