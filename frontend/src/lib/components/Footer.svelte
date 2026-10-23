<script lang="ts">
	import { bahasaBawaan, t } from '$lib/i18n';
	import Logo from './Logo.svelte';
</script>

<footer class="mt-14 border-t border-line bg-surface">
	<div class="wrap flex flex-col gap-4 py-10 md:flex-row md:items-center md:justify-between">
		<div class="flex flex-col gap-2">
			<Logo />
			<p class="text-small text-muted">{t('footer.independent')}</p>
		</div>
		<a
			href="/{bahasaBawaan}/tentang-data"
			class="inline-flex min-h-touch items-center font-bold text-secondary-on underline"
			>{t('nav.about')}</a
		>
	</div>
</footer>
