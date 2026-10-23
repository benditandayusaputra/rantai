<script lang="ts">
	import { t } from '$lib/i18n';
</script>

<svelte:head>
	<title>{t('home.title')}</title>
</svelte:head>

<main id="konten" class="wrap py-10">
	<p class="t-over">{t('home.overline')}</p>
	<h1 class="t-display mt-3">{t('home.title')}</h1>
	<p class="t-lead mt-4 max-w-prose">{t('home.lead')}</p>
</main>
