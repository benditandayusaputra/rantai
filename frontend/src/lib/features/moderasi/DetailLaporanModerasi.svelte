<script lang="ts">
	import { CircleCheck, CircleX, EyeOff, SkipForward } from '@lucide/svelte';
	import type { LaporanModerasi } from '$lib/api/tipe';
	import SimRibbon from '$lib/components/SimRibbon.svelte';
	import { formatJam, formatTanggal } from '$lib/format';
	import { bahasaBawaan, t, type KunciPesan } from '$lib/i18n';

	interface Props {
		laporan: LaporanModerasi;
		alasanTolak: string[];
		busy: boolean;
		onverifikasi: () => void;
		ontolak: (alasan: string, pesan: string) => void;
		onsembunyikan: (fotoId: number, alasan: string) => void;
		onlewati: () => void;
	}

	let { laporan, alasanTolak, busy, onverifikasi, ontolak, onsembunyikan, onlewati }: Props =
		$props();

	let alasan = $state('');
	let pesan = $state('');
	let alasanSembunyi = $state('');

	export function fokusTolak() {
		document.getElementById('alasan-tolak')?.focus();
	}
</script>

<article class="flex flex-col gap-5" aria-labelledby="judul-laporan-moderasi">
	{#if laporan.is_simulation}<SimRibbon />{/if}
	<header class="flex flex-col gap-1">
		<h2 id="judul-laporan-moderasi" class="t-h2">{t(`report.cond.${laporan.kondisi}`)}</h2>
		<p class="text-small text-muted">
			{formatTanggal(new Date(laporan.dibuat_pada), bahasaBawaan, 'long')} · {formatJam(
				new Date(laporan.dibuat_pada)
			)} · {laporan.nama_wilayah}
		</p>
	</header>

	<section aria-label={t('mod.photos')} class="flex flex-col gap-3">
		<label class="flex flex-col gap-1 text-small font-semibold">
			{t('mod.hideReason')}
			<input bind:value={alasanSembunyi} class="field font-normal" />
		</label>
		<ul class="grid grid-cols-2 gap-3">
			{#each laporan.foto as foto, i (foto.id)}
				<li class="flex flex-col gap-2">
					<a href={foto.url} class="block overflow-hidden rounded-lg border border-line bg-sunken">
						<img
							src={foto.url_thumbnail}
							alt={t('report.s3.alt', { n: i + 1, jam: formatJam(new Date(laporan.dibuat_pada)) })}
							class="aspect-square w-full object-cover {foto.disembunyikan ? 'opacity-40' : ''}"
						/>
					</a>
					{#if foto.disembunyikan}
						<p class="text-small text-muted">{t('mod.photoHidden', { n: i + 1 })}</p>
					{:else}
						<button
							type="button"
							class="btn btn-outline btn-sm"
							disabled={busy || !alasanSembunyi.trim()}
							onclick={() => onsembunyikan(foto.id, alasanSembunyi.trim())}
						>
							<EyeOff size={18} strokeWidth={1.75} aria-hidden="true" />
							{t('mod.hidePhoto', { n: i + 1 })}
						</button>
					{/if}
				</li>
			{/each}
		</ul>
	</section>

	<section class="card flex flex-col gap-2 p-4">
		<h3 class="font-bold">{t('mod.context')}</h3>
		<p>
			{laporan.paket_nama ? t('mod.project', { nama: laporan.paket_nama }) : t('mod.noProject')}
		</p>
		{#if laporan.catatan}
			<blockquote class="border-l-link border-line-strong pl-3 font-serif">
				{laporan.catatan}
			</blockquote>
		{/if}
	</section>

	<form
		class="card flex flex-col gap-3 p-4"
		onsubmit={(event) => {
			event.preventDefault();
			ontolak(alasan, pesan);
		}}
	>
		<label class="flex flex-col gap-1 text-small font-semibold">
			{t('mod.rejectReason')}
			<select id="alasan-tolak" bind:value={alasan} required class="field font-normal">
				<option value="" disabled></option>
				{#each alasanTolak as kunci (kunci)}
					<option value={kunci}>{t(`mod.reason.${kunci}` as KunciPesan)}</option>
				{/each}
			</select>
		</label>
		<label class="flex flex-col gap-1 text-small font-semibold">
			{t('mod.rejectMessage')}
			<textarea bind:value={pesan} rows="2" maxlength="280" class="field h-auto py-3 font-normal"
			></textarea>
		</label>
		<div class="flex flex-wrap gap-3">
			<button type="button" class="btn btn-primary btn-lg" disabled={busy} onclick={onverifikasi}>
				<CircleCheck size={20} strokeWidth={1.75} aria-hidden="true" />
				{t('mod.verify')}
			</button>
			<button type="submit" class="btn btn-outline btn-lg" disabled={busy || !alasan}>
				<CircleX size={20} strokeWidth={1.75} aria-hidden="true" />
				{t('mod.reject')}
			</button>
			<button type="button" class="btn btn-text btn-lg" onclick={onlewati}>
				<SkipForward size={20} strokeWidth={1.75} aria-hidden="true" />
				{t('mod.skip')}
			</button>
		</div>
	</form>
</article>
