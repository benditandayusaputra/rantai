<script lang="ts">
	import { berandaPengguna } from '$lib/features/wilayah/wilayahku.svelte';
	import { bahasaAktif, t } from '$lib/i18n';
	import LogoMark from './LogoMark.svelte';

	let { ukuran = 'normal' }: { ukuran?: 'normal' | 'besar' } = $props();

	const gaya = {
		normal: { tanda: 'h-6 lg:h-7', teks: 'text-merek lg:text-merek-lg' },
		besar: { tanda: 'h-8 lg:h-10', teks: 'text-merek-lg' }
	};
</script>

<a
	href={berandaPengguna(bahasaAktif)}
	class="inline-flex min-h-touch shrink-0 items-center gap-2.5 rounded-md"
	aria-label={t('brand.home')}
>
	<LogoMark class="{gaya[ukuran].tanda} w-auto fill-primary" />
	<span class="font-heading {gaya[ukuran].teks}"
		>{t('brand.name')}<span class="text-primary">.</span></span
	>
</a>
