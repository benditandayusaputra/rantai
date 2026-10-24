<script lang="ts">
	import { ExternalLink } from '@lucide/svelte';
	import type { ProfilAnggota } from '$lib/api/tipe';
	import { inisial } from '$lib/features/kecamatan/kecamatan';
	import { bahasaBawaan, t } from '$lib/i18n';
	import IsuDisuarakan from './IsuDisuarakan.svelte';
	import JejakKegiatan from './JejakKegiatan.svelte';
	import SikapFraksi from './SikapFraksi.svelte';

	let { anggota }: { anggota: ProfilAnggota } = $props();

	const kelasTautan =
		'inline-flex min-h-touch items-center gap-2 rounded-md border-ctrl border-line-strong px-4 font-semibold hover:border-fg hover:bg-sunken';
	const keteranganKomisi = $derived(
		anggota.komisi
			? `${anggota.komisi.nama}`
			: anggota.keterangan_komisi === 'pimpinan_dprd'
				? t('member.leader')
				: t('reps.noCommission')
	);
</script>

<svelte:head>
	<title>{anggota.nama} · {t('brand.name')}</title>
</svelte:head>

<article class="wrap flex flex-col gap-8 py-6 lg:py-10">
	<header class="flex flex-col gap-4">
		<p class="t-over">{t('member.overline')}</p>
		<div class="flex items-center gap-4">
			<span
				class="inline-flex h-18 w-18 shrink-0 items-center justify-center rounded-pill bg-secondary text-h2 font-bold text-secondary-on"
				aria-hidden="true">{inisial(anggota.nama)}</span
			>
			<div class="flex flex-col gap-1">
				<h1 class="t-h1">{anggota.nama}</h1>
				<p class="text-muted">{anggota.fraksi} · {anggota.partai}</p>
			</div>
		</div>
		<p
			class="self-start rounded-pill bg-secondary px-3 py-1 text-small font-bold text-secondary-on"
		>
			{keteranganKomisi}
		</p>
		<div class="flex flex-wrap gap-3">
			{#each [{ href: anggota.lhkpn_url, label: t('member.lhkpn') }, { href: anggota.url_profil, label: t('member.official') }] as tautan (tautan.label)}
				{#if tautan.href}
					<a href={tautan.href} target="_blank" rel="noopener" class={kelasTautan}>
						{tautan.label}<ExternalLink size={16} strokeWidth={1.75} aria-hidden="true" /><span
							class="sr-only">{t('a11y.opensNewTab')}</span
						>
					</a>
				{/if}
			{/each}
		</div>
		<p class="max-w-prose text-small text-muted">{t('member.noScore')}</p>
	</header>

	<div class="grid gap-8 lg:grid-cols-12">
		<div class="flex flex-col gap-10 lg:col-span-8">
			<IsuDisuarakan idAnggota={anggota.id} />
			<JejakKegiatan idAnggota={anggota.id} />
			<SikapFraksi fraksi={anggota.fraksi} />
		</div>
		<aside class="flex flex-col gap-4 lg:col-span-4">
			<section class="flex flex-col gap-3 rounded-lg border border-line bg-surface p-4">
				<h2 class="t-h3">
					{t('member.dapil', { nomor: anggota.dapil.nomor, nama: anggota.dapil.nama })}
				</h2>
				<ul class="flex flex-wrap gap-2">
					{#each anggota.dapil.kecamatan as kec (kec.kode)}
						<li>
							<a
								href="/{bahasaBawaan}/kecamatan/{kec.kode}"
								class="inline-flex min-h-touch items-center rounded-pill bg-sunken px-3 text-small font-semibold hover:bg-secondary"
								>{kec.nama}</a
							>
						</li>
					{/each}
				</ul>
			</section>
			{#if anggota.komisi?.bidang?.length}
				<section class="flex flex-col gap-2 rounded-lg border border-line bg-surface p-4">
					<h2 class="t-h3">{t('member.commission')}</h2>
					<p class="font-semibold">{anggota.komisi.nama}</p>
					<ul class="list-disc pl-5 text-small">
						{#each anggota.komisi.bidang as bidang (bidang)}<li>{bidang}</li>{/each}
					</ul>
				</section>
			{/if}
			{#if anggota.alat_kelengkapan.length}
				<section class="flex flex-col gap-2 rounded-lg border border-line bg-surface p-4">
					<h2 class="t-h3">{t('member.bodies')}</h2>
					<ul class="flex flex-col gap-1 text-small">
						{#each anggota.alat_kelengkapan as akd (akd.alat_kelengkapan)}
							<li><span class="font-semibold">{akd.alat_kelengkapan}</span> · {akd.jabatan}</li>
						{/each}
					</ul>
				</section>
			{/if}
		</aside>
	</div>
</article>
