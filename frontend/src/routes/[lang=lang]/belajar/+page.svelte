<script lang="ts">
	import { BookOpen, Info, MapPin, MessageCircle } from '@lucide/svelte';
	import { daftarIstilah, langkahPemerintahan } from '$lib/features/belajar/belajar';
	import { tautanDenganWilayah } from '$lib/features/navigasi/navigasi';
	import { wilayahku } from '$lib/features/wilayah/wilayahku.svelte';
	import { bahasaAktif, t } from '$lib/i18n';
</script>

<svelte:head>
	<title>{t('nav.learn')} · {t('brand.name')}</title>
</svelte:head>

<div class="wrap pt-8 lg:pt-12">
	<header class="border-b border-line pb-6">
		<div class="flex items-center gap-4">
			<span class="icon-tile h-12 w-12"
				><BookOpen size={24} strokeWidth={1.75} aria-hidden="true" /></span
			>
			<h1 class="t-h1">{t('learn.title')}</h1>
		</div>
		<p class="mt-4 max-w-prose text-muted">{t('learn.lead')}</p>
	</header>

	<section class="py-8 lg:py-10" aria-labelledby="judul-cara">
		<h2 id="judul-cara" class="t-h2">{t('learn.how.title')}</h2>
		<p class="mt-2 max-w-prose text-muted">{t('learn.how.lead')}</p>
		<ol class="mt-6 grid gap-4 md:grid-cols-2 lg:grid-cols-3 lg:gap-6">
			{#each langkahPemerintahan as langkah, i (langkah)}
				<li class="card flex flex-col gap-2 p-5">
					<span
						class="num inline-flex h-8 w-8 items-center justify-center rounded-pill bg-secondary text-small font-bold text-secondary-on"
						aria-hidden="true">{i + 1}</span
					>
					<h3 class="t-h3 mt-1">{t(`learn.step.${langkah}.title`)}</h3>
					<p class="text-small text-muted">{t(`learn.step.${langkah}.body`)}</p>
				</li>
			{/each}
		</ol>
		<div class="mt-6 flex gap-3 rounded-lg bg-info-bg p-4">
			<Info size={20} strokeWidth={1.75} class="mt-0.5 shrink-0 text-info" aria-hidden="true" />
			<div>
				<h3 class="font-semibold">{t('learn.jakartaNote.title')}</h3>
				<p class="mt-1 max-w-prose text-small">{t('learn.jakartaNote.body')}</p>
			</div>
		</div>
	</section>

	<section class="border-t border-line py-8 lg:py-10" aria-labelledby="judul-istilah">
		<h2 id="judul-istilah" class="t-h2">{t('learn.terms.title')}</h2>
		<p class="mt-2 max-w-prose text-muted">{t('learn.terms.lead')}</p>
		<dl class="mt-4 grid gap-x-10 md:grid-cols-2">
			{#each daftarIstilah as istilah (istilah)}
				<div class="border-b border-line py-4">
					<dt class="font-semibold">{t(`learn.term.${istilah}`)}</dt>
					<dd class="mt-1 text-small text-muted">{t(`learn.term.${istilah}.def`)}</dd>
				</div>
			{/each}
		</dl>
	</section>

	<div
		class="flex flex-col gap-4 rounded-lg bg-sunken p-5 md:flex-row md:items-center md:justify-between md:p-6"
	>
		<div>
			<h2 class="t-h3">{t('learn.next.title')}</h2>
			<p class="mt-1 text-small text-muted">{t('learn.next.body')}</p>
		</div>
		<div class="flex shrink-0 flex-wrap gap-2">
			<a href={tautanDenganWilayah('/wilayah', bahasaAktif, wilayahku())} class="btn btn-primary"
				><MapPin size={18} aria-hidden="true" />{t('learn.next.area')}</a
			>
			<a href={tautanDenganWilayah('/tanya', bahasaAktif, wilayahku())} class="btn btn-secondary"
				><MessageCircle size={18} aria-hidden="true" />{t('ask.title')}</a
			>
		</div>
	</div>
</div>
