<script lang="ts">
	import {
		CircleCheck,
		CircleHelp,
		ExternalLink,
		EyeOff,
		LoaderCircle,
		LogIn,
		MessageSquarePlus,
		ThumbsDown,
		ThumbsUp,
		Users
	} from '@lucide/svelte';
	import { tangkap } from '$lib/api/klien';
	import { ambilTemuan, kirimTemuan, sembunyikanTemuan, tanggapiTemuan } from '$lib/api/temuan';
	import type { SikapTemuan, StatusTemuan, Temuan } from '$lib/api/tipe';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import ErrorState from '$lib/components/ErrorState.svelte';
	import TextArea from '$lib/components/TextArea.svelte';
	import { sesiAktif } from '$lib/features/akun/sesi.svelte';
	import { formatTanggal } from '$lib/format';
	import { bahasaAktif, t } from '$lib/i18n';
	import { isIsiSah, isTautanSah, panjangMaksIsi, panjangMaksTautan } from './temuan';

	let { kodePaket }: { kodePaket: string } = $props();

	const ambangKonfirmasi = 2;
	const awalanTautan = 'https://';
	const kelasStatus: Record<StatusTemuan, string> = {
		menunggu_konfirmasi: 'border-dashed border-st-wait-bd bg-st-wait-bg text-st-wait-fg',
		terkonfirmasi: 'border-st-ver-bd bg-st-ver-bg text-st-ver-fg',
		diperdebatkan: 'border-st-resp-bd bg-st-resp-bg text-st-resp-fg'
	};
	const ikonStatus = {
		menunggu_konfirmasi: CircleHelp,
		terkonfirmasi: CircleCheck,
		diperdebatkan: Users
	} satisfies Record<StatusTemuan, typeof CircleHelp>;

	const sesi = $derived(sesiAktif());
	const isModerator = $derived(
		sesi?.pengguna.peran === 'moderator' || sesi?.pengguna.peran === 'admin'
	);

	let daftar = $state<Temuan[]>([]);
	let status = $state<'memuat' | 'siap' | 'gagal'>('memuat');
	let isi = $state('');
	let tautan = $state('');
	let pesanForm = $state('');
	let isGagalForm = $state(false);
	let isMengirim = $state(false);
	let tautanCek = $state<Record<number, string>>({});
	let pesanCek = $state<Record<number, string>>({});
	let alasanSembunyi = $state<Record<number, string>>({});
	let sibukID = $state<number | null>(null);

	async function muat() {
		const hasil = await tangkap(ambilTemuan(fetch, kodePaket));
		daftar = hasil.ok ? hasil.data : [];
		status = hasil.ok ? 'siap' : 'gagal';
	}

	$effect(() => {
		void muat();
	});

	async function kirim(event: SubmitEvent) {
		event.preventDefault();
		if (!sesi) return;
		isGagalForm = !isIsiSah(isi) || !isTautanSah(tautan.trim());
		if (isGagalForm) {
			pesanForm = t('finding.invalid');
			return;
		}
		isMengirim = true;
		const hasil = await tangkap(
			kirimTemuan(fetch, sesi.token_akses, kodePaket, { isi, tautan_sumber: tautan.trim() })
		);
		isMengirim = false;
		isGagalForm = !hasil.ok;
		pesanForm = hasil.ok ? t('finding.sent') : hasil.galat.message;
		if (!hasil.ok) return;
		isi = '';
		tautan = '';
		await muat();
	}

	async function tanggapi(temuan: Temuan, sikap: SikapTemuan) {
		if (!sesi) return;
		const bukti = (tautanCek[temuan.id] ?? '').trim();
		if (bukti && !isTautanSah(bukti)) {
			pesanCek[temuan.id] = t('finding.linkHelp');
			return;
		}
		sibukID = temuan.id;
		const hasil = await tangkap(
			tanggapiTemuan(fetch, sesi.token_akses, temuan.id, { sikap, tautan_sumber: bukti })
		);
		sibukID = null;
		pesanCek[temuan.id] = hasil.ok ? t('finding.thanks') : hasil.galat.message;
		if (hasil.ok) await muat();
	}

	async function sembunyikan(temuan: Temuan) {
		const alasan = (alasanSembunyi[temuan.id] ?? '').trim();
		if (!sesi || !alasan) return;
		sibukID = temuan.id;
		const hasil = await tangkap(sembunyikanTemuan(fetch, sesi.token_akses, temuan.id, alasan));
		sibukID = null;
		if (hasil.ok) await muat();
		else pesanCek[temuan.id] = hasil.galat.message;
	}
</script>

<section aria-labelledby="judul-temuan" class="flex flex-col gap-4">
	<h2 id="judul-temuan" class="t-h2 flex items-center gap-3">
		<span class="icon-tile"><Users size={20} strokeWidth={1.75} aria-hidden="true" /></span>
		{t('finding.title')}
	</h2>
	<p class="max-w-prose text-muted">{t('finding.lead')}</p>

	{#if status === 'memuat'}
		<div class="skel h-32 rounded-lg" aria-busy="true"></div>
	{:else if status === 'gagal'}
		<ErrorState title={t('finding.failed')} onretry={muat} />
	{:else if !daftar.length}
		<EmptyState title={t('finding.empty')} />
	{:else}
		<ul class="flex flex-col gap-3">
			{#each daftar as temuan (temuan.id)}
				{@const IkonStatus = ikonStatus[temuan.status]}
				<li class="card flex flex-col gap-3 p-4 md:p-5">
					<div class="flex flex-wrap items-center justify-between gap-2">
						<span
							class="inline-flex min-h-8 items-center gap-1.5 rounded-xs border-ctrl px-2.5 text-small font-bold {kelasStatus[
								temuan.status
							]}"
							><IkonStatus size={16} strokeWidth={1.75} aria-hidden="true" />{t(
								`finding.status.${temuan.status}`
							)}</span
						>
						<span class="text-small text-muted num"
							>{t('finding.counts', {
								dukung: temuan.jumlah_mendukung,
								bantah: temuan.jumlah_membantah
							})}</span
						>
					</div>
					<p class="whitespace-pre-line break-words">{temuan.isi}</p>
					<div class="flex flex-wrap items-center gap-x-4 gap-y-2 text-small text-muted">
						<span
							>{t('finding.by', {
								nama: temuan.nama_kontributor || t('finding.anonymous'),
								tanggal: formatTanggal(new Date(temuan.dibuat_pada), bahasaAktif)
							})}</span
						>
						{#if isTautanSah(temuan.tautan_sumber)}
							<a
								href={temuan.tautan_sumber}
								target="_blank"
								rel="nofollow ugc noopener noreferrer"
								class="link inline-flex items-center gap-1"
								><ExternalLink size={14} aria-hidden="true" />{t('finding.source')}<span
									class="sr-only">{t('a11y.opensNewTab')}</span
								></a
							>
						{/if}
					</div>
					{#if sesi}
						<div class="flex flex-col gap-2 border-t border-line pt-3">
							<label class="flex flex-col gap-1 text-small font-semibold">
								{t('finding.checkSource')}
								<input
									type="url"
									inputmode="url"
									maxlength={panjangMaksTautan}
									bind:value={tautanCek[temuan.id]}
									class="field font-normal"
									placeholder={awalanTautan}
								/>
							</label>
							<div class="flex flex-wrap gap-2">
								{#each [{ sikap: 'mendukung', ikon: ThumbsUp, label: t('finding.support') }, { sikap: 'membantah', ikon: ThumbsDown, label: t('finding.dispute') }] as aksi (aksi.sikap)}
									{@const Ikon = aksi.ikon}
									<button
										type="button"
										class="btn btn-outline text-small"
										disabled={sibukID === temuan.id}
										onclick={() => tanggapi(temuan, aksi.sikap as SikapTemuan)}
										><Ikon size={16} strokeWidth={1.75} aria-hidden="true" />{aksi.label}</button
									>
								{/each}
							</div>
							{#if pesanCek[temuan.id]}<p class="text-small" role="status">
									{pesanCek[temuan.id]}
								</p>{/if}
							{#if isModerator}
								<div class="flex flex-wrap items-end gap-2">
									<label class="flex min-w-0 flex-1 flex-col gap-1 text-small font-semibold">
										{t('finding.hideReason')}
										<input bind:value={alasanSembunyi[temuan.id]} class="field font-normal" />
									</label>
									<button
										type="button"
										class="btn btn-ghost text-small"
										disabled={sibukID === temuan.id || !alasanSembunyi[temuan.id]?.trim()}
										onclick={() => sembunyikan(temuan)}
										><EyeOff size={16} strokeWidth={1.75} aria-hidden="true" />{t(
											'finding.hide'
										)}</button
									>
								</div>
							{/if}
						</div>
					{/if}
				</li>
			{/each}
		</ul>
	{/if}

	{#if sesi}
		<form class="card flex flex-col gap-4 p-4 md:p-5" onsubmit={kirim} novalidate>
			<h3 class="t-h3 flex items-center gap-2">
				<MessageSquarePlus size={20} strokeWidth={1.75} aria-hidden="true" />{t('finding.add')}
			</h3>
			<p class="text-small text-muted">{t('finding.rule', { n: ambangKonfirmasi })}</p>
			<TextArea
				id="isi-temuan"
				label={t('finding.body')}
				help={t('finding.bodyHelp')}
				bind:value={isi}
				maxlength={panjangMaksIsi}
				rows={4}
			/>
			<label class="flex flex-col gap-2">
				<span class="font-semibold">{t('finding.link')}</span>
				<span class="text-small text-muted">{t('finding.linkHelp')}</span>
				<input
					type="url"
					inputmode="url"
					required
					maxlength={panjangMaksTautan}
					bind:value={tautan}
					class="field"
					placeholder={awalanTautan}
				/>
			</label>
			{#if pesanForm}
				<p
					class="rounded-md p-3 text-small {isGagalForm
						? 'bg-danger-bg text-danger'
						: 'bg-success-bg'}"
					role={isGagalForm ? 'alert' : 'status'}
				>
					{pesanForm}
				</p>
			{/if}
			<button type="submit" class="btn btn-primary self-start" disabled={isMengirim}>
				{#if isMengirim}<LoaderCircle
						size={18}
						strokeWidth={1.75}
						class="animate-spin"
						aria-hidden="true"
					/>{/if}
				{t('finding.submit')}
			</button>
		</form>
	{:else}
		<a
			href="/{bahasaAktif}/masuk?lanjut={encodeURIComponent(
				`/${bahasaAktif}/paket/${encodeURIComponent(kodePaket)}`
			)}"
			class="btn btn-outline self-start"
			><LogIn size={18} strokeWidth={1.75} aria-hidden="true" />{t('finding.login')}</a
		>
	{/if}
</section>
