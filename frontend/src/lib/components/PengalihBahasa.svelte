<script lang="ts">
	import { page } from '$app/state';
	import { bahasaAktif, bahasaTersedia, gantiBahasaPath, t } from '$lib/i18n';
</script>

<div
	role="group"
	aria-label={t('lang.label')}
	class="inline-flex shrink-0 rounded-pill border border-line bg-raised p-0.5"
>
	{#each bahasaTersedia as bahasa (bahasa)}
		<a
			href={gantiBahasaPath(page.url.pathname, bahasa) + page.url.search}
			hreflang={bahasa}
			lang={bahasa}
			aria-current={bahasa === bahasaAktif ? 'true' : undefined}
			title={t(`lang.${bahasa}`)}
			class="inline-flex h-10 min-w-10 items-center justify-center rounded-pill px-2 text-caption font-bold uppercase text-muted transition-colors duration-fast ease-standard hover:text-fg aria-[current=true]:bg-secondary aria-[current=true]:text-secondary-on"
		>
			<span aria-hidden="true">{bahasa}</span>
			<span class="sr-only">{t(`lang.${bahasa}`)}</span>
		</a>
	{/each}
</div>
