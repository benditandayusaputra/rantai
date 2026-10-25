<script lang="ts">
	import { Building2, ExternalLink, Info, Landmark, MapPin, Users } from '@lucide/svelte';
	import type { ProfilAnggota } from '$lib/api/tipe';
	import FotoWakil from '$lib/components/FotoWakil.svelte';
	import { bahasaAktif, t } from '$lib/i18n';
	import IsuDisuarakan from './IsuDisuarakan.svelte';
	import JejakKegiatan from './JejakKegiatan.svelte';
	import SikapFraksi from './SikapFraksi.svelte';

	let { anggota }: { anggota: ProfilAnggota } = $props();

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

<article class="wrap flex flex-col gap-10 py-8 lg:py-12">
	<header class="flex flex-col gap-4">
		<div
			class="card flex flex-col gap-5 p-5 md:p-6 lg:flex-row lg:items-center lg:justify-between lg:p-8"
		>
			<div class="flex items-start gap-4 md:items-center">
				<FotoWakil
					nama={anggota.nama}
					foto={anggota.foto_url}
					class="h-20 w-20 font-heading text-h3 md:h-24 md:w-24 md:text-h2"
				/>
				<div class="flex min-w-0 flex-col gap-2">
					<p class="t-over">{t('member.overline')}</p>
					<h1 class="t-h1">{anggota.nama}</h1>
					<p class="text-muted">{anggota.fraksi} · {anggota.partai}</p>
					{#if anggota.foto_url}
						<p class="text-caption text-muted">{t('photo.dprdShort')}</p>
					{/if}
					<p
						class="inline-flex items-center gap-2 self-start rounded-xs bg-secondary px-3 py-1 text-small font-bold text-secondary-on"
					>
						<Landmark size={16} strokeWidth={1.75} aria-hidden="true" />{keteranganKomisi}
					</p>
				</div>
			</div>
			<div class="flex flex-col gap-2 md:flex-row md:flex-wrap lg:flex-col lg:items-stretch">
				{#each [{ href: anggota.lhkpn_url, label: t('member.lhkpn') }, { href: anggota.url_profil, label: t('member.official') }] as tautan (tautan.label)}
					{#if tautan.href}
						<a href={tautan.href} target="_blank" rel="noopener" class="btn btn-outline">
							{tautan.label}<ExternalLink size={18} strokeWidth={1.75} aria-hidden="true" /><span
								class="sr-only">{t('a11y.opensNewTab')}</span
							>
						</a>
					{/if}
				{/each}
			</div>
		</div>
		<p class="flex max-w-prose items-start gap-2 text-small text-muted">
			<Info size={16} strokeWidth={1.75} aria-hidden="true" class="mt-0.5 shrink-0" />
			{t('member.noScore')}
		</p>
	</header>

	<div class="grid gap-10 lg:grid-cols-12 lg:gap-8">
		<div class="flex flex-col gap-12 lg:col-span-8">
			<IsuDisuarakan idAnggota={anggota.id} />
			<JejakKegiatan idAnggota={anggota.id} />
			<SikapFraksi fraksi={anggota.fraksi} />
		</div>
		<aside class="flex flex-col gap-4 lg:sticky lg:top-24 lg:col-span-4 lg:self-start">
			<section class="card flex flex-col gap-4 p-5">
				<h2 class="flex items-center gap-3 text-h3">
					<span class="icon-tile h-9 w-9"
						><MapPin size={18} strokeWidth={1.75} aria-hidden="true" /></span
					>
					{t('member.dapil', { nomor: anggota.dapil.nomor, nama: anggota.dapil.nama })}
				</h2>
				<ul class="flex flex-wrap gap-2">
					{#each anggota.dapil.kecamatan as kec (kec.kode)}
						<li>
							<a href="/{bahasaAktif}/kecamatan/{kec.kode}" class="chip">{kec.nama}</a>
						</li>
					{/each}
				</ul>
			</section>
			{#if anggota.komisi?.bidang?.length}
				<section class="card flex flex-col gap-3 p-5">
					<h2 class="flex items-center gap-3 text-h3">
						<span class="icon-tile h-9 w-9"
							><Landmark size={18} strokeWidth={1.75} aria-hidden="true" /></span
						>
						{t('member.commission')}
					</h2>
					<p class="font-semibold">{anggota.komisi.nama}</p>
					<ul class="flex flex-col gap-2 text-small">
						{#each anggota.komisi.bidang as bidang (bidang)}
							<li class="flex items-start gap-2">
								<Building2
									size={16}
									strokeWidth={1.75}
									aria-hidden="true"
									class="mt-0.5 shrink-0 text-muted"
								/>{bidang}
							</li>
						{/each}
					</ul>
				</section>
			{/if}
			{#if anggota.alat_kelengkapan.length}
				<section class="card flex flex-col gap-3 p-5">
					<h2 class="flex items-center gap-3 text-h3">
						<span class="icon-tile h-9 w-9"
							><Users size={18} strokeWidth={1.75} aria-hidden="true" /></span
						>
						{t('member.bodies')}
					</h2>
					<ul class="flex flex-col divide-y divide-line text-small">
						{#each anggota.alat_kelengkapan as akd (akd.alat_kelengkapan)}
							<li class="flex flex-col py-2 first:pt-0 last:pb-0">
								<span class="font-semibold">{akd.alat_kelengkapan}</span>
								<span class="text-muted">{akd.jabatan}</span>
							</li>
						{/each}
					</ul>
				</section>
			{/if}
		</aside>
	</div>
</article>
