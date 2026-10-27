<script lang="ts">
	import { Calendar, MapPin } from '@lucide/svelte';
	import type { DetailPaket } from '$lib/api/tipe';
	import { formatTanggal } from '$lib/format';
	import { bahasaAktif, t } from '$lib/i18n';
	let { paket }: { paket: DetailPaket } = $props();
	const tanggal = (nilai: string | null | undefined) =>
		nilai && Number.isFinite(Date.parse(nilai))
			? formatTanggal(new Date(nilai), bahasaAktif)
			: t('project.notPublished');
	const tautanPeta = $derived(
		paket.lokasi.titik
			? `/${bahasaAktif}/peta?lapisan=proyek&lat=${paket.lokasi.titik.lat}&lng=${paket.lokasi.titik.lng}`
			: `/${bahasaAktif}/peta?lapisan=proyek&w=${paket.lokasi.kode_wilayah.split('.').slice(0, 3).join('.')}`
	);
</script>

<section
	class="rounded-lg border border-line bg-surface p-5 md:p-6"
	aria-labelledby="judul-ringkasan-proyek"
>
	<div class="flex flex-wrap items-center justify-between gap-3">
		<h2 id="judul-ringkasan-proyek" class="t-h2">{t('project.summary')}</h2>
		{#if paket.sinyal.length}<span
				class="rounded-xs bg-attn-bg px-3 py-1 text-caption font-semibold text-attn"
				>{t('services.issue')}</span
			>{/if}
	</div>
	<dl class="mt-5 grid gap-5 md:grid-cols-2">
		<div>
			<dt class="text-small text-muted">{t('project.responsible')}</dt>
			<dd class="mt-1 font-semibold">{paket.instansi || t('project.notPublished')}</dd>
			<dd class="mt-1 text-caption text-muted">
				{t(
					paket.instansi === 'Provinsi DKI Jakarta' ? 'project.agencyBroad' : 'project.agencyRole'
				)}
			</dd>
		</div>
		<div>
			<dt class="text-small text-muted">{t('pkg.vendor')}</dt>
			<dd class="mt-1 font-semibold">{paket.penyedia?.nama ?? t('project.notPublished')}</dd>
			<dd class="mt-1 text-caption text-muted">{t('project.vendorRole')}</dd>
		</div>
		<div>
			<dt class="flex items-center gap-2 text-small text-muted">
				<Calendar size={16} aria-hidden="true" />{t('project.plannedStart')}
			</dt>
			<dd class="mt-1 font-semibold">{tanggal(paket.rup?.pelaksanaan_mulai)}</dd>
		</div>
		<div>
			<dt class="text-small text-muted">{t('project.plannedEnd')}</dt>
			<dd class="mt-1 font-semibold">{tanggal(paket.rup?.pelaksanaan_selesai)}</dd>
		</div>
		<div>
			<dt class="text-small text-muted">{t('project.procurementStage')}</dt>
			<dd class="mt-1 font-semibold">{paket.tahap || t('project.notPublished')}</dd>
		</div>
		<div>
			<dt class="text-small text-muted">{t('project.handover')}</dt>
			<dd class="mt-1 font-semibold">{tanggal(paket.tanggal_serah_terima)}</dd>
		</div>
	</dl>
	<p class="mt-5 border-t border-line pt-4 text-small text-muted">{t('project.timelineNote')}</p>
	<div class="mt-4 rounded-md bg-sunken p-4">
		<h3 class="text-small font-semibold">{t('project.impactTitle')}</h3>
		<p class="mt-2 text-small text-muted">{t('project.impactUnknown')}</p>
	</div>
	{#if paket.lokasi.nama_wilayah}
		<div class="mt-5 flex flex-wrap items-center justify-between gap-3">
			<div>
				<p class="flex items-center gap-2 text-small font-semibold">
					<MapPin size={18} aria-hidden="true" />{paket.lokasi.nama_wilayah}
				</p>
				<p class="mt-1 text-caption text-muted">
					{t(paket.lokasi.titik ? 'project.pointAvailable' : 'project.areaOnly')}
				</p>
			</div>
			<a href={tautanPeta} class="btn btn-outline">{t('project.viewMap')}</a>
		</div>
	{/if}
</section>
