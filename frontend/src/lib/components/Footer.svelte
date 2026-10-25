<script lang="ts">
	import { kelompokKaki } from '$lib/features/navigasi/navigasi';
	import { bahasaAktif, t } from '$lib/i18n';
	import Logo from './Logo.svelte';
	import PengalihTema from './PengalihTema.svelte';
</script>

<footer class="mt-16 border-t border-line bg-sunken lg:mt-20">
	<div class="wrap grid gap-10 py-10 md:grid-cols-2 lg:grid-cols-[2fr_1fr_1fr_1fr] lg:py-12">
		<div class="flex max-w-sm flex-col gap-3">
			<Logo />
			<p class="text-small text-muted">{t('footer.tagline')}</p>
			<p class="text-small font-semibold">{t('footer.independent')}</p>
		</div>
		{#each kelompokKaki as kelompok (kelompok.judul)}
			<nav aria-label={t(kelompok.judul)}>
				<p class="t-over">{t(kelompok.judul)}</p>
				<ul class="mt-2 flex flex-col">
					{#each kelompok.tautan as tautan (tautan.path)}
						<li>
							<a
								href="/{bahasaAktif}{tautan.path}"
								class="inline-flex min-h-touch items-center text-small font-medium text-fg underline-offset-4 hover:underline"
							>
								{t(tautan.label)}
							</a>
						</li>
					{/each}
				</ul>
			</nav>
		{/each}
	</div>
	<div class="border-t border-line">
		<div class="wrap flex flex-wrap items-center justify-between gap-3 py-4">
			<p class="text-caption text-muted">{t('footer.event')}</p>
			<PengalihTema denganLabel />
		</div>
	</div>
</footer>
