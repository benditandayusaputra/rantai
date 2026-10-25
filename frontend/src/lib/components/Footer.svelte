<script lang="ts">
	import { kelompokKaki } from '$lib/features/navigasi/navigasi';
	import { bahasaBawaan, t } from '$lib/i18n';
	import Logo from './Logo.svelte';
</script>

<footer class="mt-16 border-t border-line bg-surface">
	<div class="wrap grid gap-10 py-12 md:grid-cols-2 lg:grid-cols-[2fr_1fr_1fr_1fr]">
		<div class="flex max-w-sm flex-col gap-3">
			<Logo />
			<p class="text-small text-muted">{t('footer.tagline')}</p>
			<p class="text-small font-semibold">{t('footer.independent')}</p>
		</div>
		{#each kelompokKaki as kelompok (kelompok.judul)}
			<nav aria-label={t(kelompok.judul)}>
				<p class="t-over">{t(kelompok.judul)}</p>
				<ul class="mt-3 flex flex-col">
					{#each kelompok.tautan as tautan (tautan.path)}
						{@const Ikon = tautan.ikon}
						<li>
							<a
								href="/{bahasaBawaan}{tautan.path}"
								class="inline-flex min-h-touch items-center gap-2 text-small font-semibold text-muted hover:text-fg"
							>
								<Ikon size={16} strokeWidth={1.75} aria-hidden="true" />
								{t(tautan.label)}
							</a>
						</li>
					{/each}
				</ul>
			</nav>
		{/each}
	</div>
	<div class="border-t border-line">
		<p class="wrap py-5 text-caption text-muted">{t('footer.event')}</p>
	</div>
</footer>
