<script lang="ts">
	import { bahasaAktif, t } from '$lib/i18n';
	import LogoMark from './LogoMark.svelte';
</script>

<a
	href="/{bahasaAktif}"
	class="inline-flex min-h-touch shrink-0 items-center gap-2.5 rounded-md"
	aria-label={t('brand.home')}
>
	<LogoMark class="h-6 w-auto fill-primary lg:h-7" />
	<span class="font-heading text-merek lg:text-merek-lg"
		>{t('brand.name')}<span class="text-primary">.</span></span
	>
</a>
