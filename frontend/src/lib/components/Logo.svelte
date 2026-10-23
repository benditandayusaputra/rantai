<script lang="ts">
	import { bahasaBawaan, t } from '$lib/i18n';
</script>

<a
	href="/{bahasaBawaan}"
	class="inline-flex min-h-touch items-center gap-2 rounded-md"
	aria-label={t('brand.home')}
>
	<svg viewBox="0 0 32 32" class="h-7 w-7" fill="none" stroke-width="3.5" aria-hidden="true">
		<rect x="3" y="10" width="16" height="12" rx="6" class="stroke-primary" />
		<rect x="13" y="10" width="16" height="12" rx="6" class="stroke-fg" />
	</svg>
	<span class="text-h3 font-extrabold tracking-tight">{t('brand.name')}</span>
</a>
