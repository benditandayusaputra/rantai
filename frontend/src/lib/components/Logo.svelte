<script lang="ts">
	import { bahasaBawaan, t } from '$lib/i18n';
	import LogoMark from './LogoMark.svelte';
</script>

<a
	href="/{bahasaBawaan}"
	class="inline-flex min-h-touch items-center gap-2 rounded-md"
	aria-label={t('brand.home')}
>
	<LogoMark />
	<span class="text-h3 font-extrabold tracking-tight">{t('brand.name')}</span>
</a>
