<script lang="ts">
	import { ArrowLeft, ArrowRight, GraduationCap, Landmark, Wallet, X } from '@lucide/svelte';
	import { tick } from 'svelte';
	import LogoMark from '$lib/components/LogoMark.svelte';
	import { t } from '$lib/i18n';
	import {
		elemenTarget,
		isTurSelesai,
		langkahTersedia,
		langkahTur,
		posisiKartu,
		type Kotak
	} from './tur';
	import { mulaiTur, tur, tutupTur } from './tur.svelte';

	const jedaTampilMs = 800;
	const ruangSorot = 6;
	const lensa = [
		{ ikon: GraduationCap, teks: t('tour.lens.services') },
		{ ikon: Wallet, teks: t('tour.lens.money') },
		{ ikon: Landmark, teks: t('tour.lens.accountability') }
	];

	let dialog = $state<HTMLDialogElement>();
	let tombolUtama = $state<HTMLButtonElement>();
	let daftar = $state(langkahTur);
	let indeks = $state(0);
	let kotakTarget = $state<Kotak | null>(null);
	let kartu = $state({ lebar: 0, tinggi: 0 });
	let layar = $state({ lebar: 0, tinggi: 0 });

	const langkah = $derived(daftar[indeks]);
	const isTerakhir = $derived(indeks === daftar.length - 1);
	const posisi = $derived(
		kotakTarget
			? posisiKartu(kotakTarget, kartu, layar)
			: { x: (layar.lebar - kartu.lebar) / 2, y: (layar.tinggi - kartu.tinggi) / 2 }
	);

	$effect(() => {
		if (isTurSelesai()) return;
		const jeda = setTimeout(mulaiTur, jedaTampilMs);
		return () => clearTimeout(jeda);
	});

	$effect(() => {
		if (!dialog) return;
		if (tur.isTerbuka && !dialog.open) bukaDialog(dialog);
		if (!tur.isTerbuka && dialog.open) dialog.close();
	});

	$effect(() => {
		if (!tur.isTerbuka || !langkah) return;
		const el = langkah.target ? elemenTarget(langkah.target) : null;
		el?.scrollIntoView({ block: 'nearest' });
		ukurTarget();
		void tick().then(() => tombolUtama?.focus());
	});

	function bukaDialog(el: HTMLDialogElement) {
		daftar = langkahTersedia(langkahTur, (target) => elemenTarget(target) !== null);
		indeks = 0;
		el.showModal();
	}

	function ukurTarget() {
		layar = { lebar: innerWidth, tinggi: innerHeight };
		const el = langkah?.target ? elemenTarget(langkah.target) : null;
		if (!el) {
			kotakTarget = null;
			return;
		}
		const r = el.getBoundingClientRect();
		kotakTarget = {
			x: r.x - ruangSorot,
			y: r.y - ruangSorot,
			lebar: r.width + 2 * ruangSorot,
			tinggi: r.height + 2 * ruangSorot
		};
	}

	function lanjut() {
		if (isTerakhir) dialog?.close();
		else indeks += 1;
	}
</script>

<svelte:window onresize={ukurTarget} onscroll={ukurTarget} />

<dialog
	bind:this={dialog}
	onclose={tutupTur}
	aria-labelledby="judul-tur"
	aria-describedby="isi-tur"
	class="m-0 h-full max-h-none w-full max-w-none overflow-hidden bg-transparent p-0 text-fg backdrop:bg-transparent"
>
	{#if langkah}
		{#if kotakTarget}
			<div
				class="sorot-tur pointer-events-none fixed rounded-md transition-all duration-enter ease-standard motion-reduce:transition-none"
				style:left="{kotakTarget.x}px"
				style:top="{kotakTarget.y}px"
				style:width="{kotakTarget.lebar}px"
				style:height="{kotakTarget.tinggi}px"
				aria-hidden="true"
			></div>
		{:else}
			<div class="fixed inset-0 bg-scrim" aria-hidden="true"></div>
		{/if}

		{#if langkah.target === null}
			<div class="fixed inset-0 flex items-center justify-center overflow-y-auto p-4">
				<div class="card w-full max-w-form p-6 shadow-3 md:p-8">
					<div class="flex items-center gap-3">
						<LogoMark class="h-8 w-auto fill-primary" />
						<p class="t-over text-primary">{t('tour.welcome.over')}</p>
					</div>
					<h2 id="judul-tur" class="t-h1 mt-4">{t(langkah.judul)}</h2>
					<p id="isi-tur" class="mt-2 text-muted">{t(langkah.isi)}</p>
					<ul class="mt-5 flex flex-col gap-3">
						{#each lensa as item (item.teks)}
							{@const Ikon = item.ikon}
							<li class="flex items-center gap-3 text-small font-semibold">
								<span class="icon-tile h-10 w-10 text-primary">
									<Ikon size={18} strokeWidth={1.75} aria-hidden="true" />
								</span>
								{item.teks}
							</li>
						{/each}
					</ul>
					<div class="mt-6 flex flex-col-reverse gap-2 md:flex-row md:justify-end">
						<button type="button" class="btn btn-text" onclick={() => dialog?.close()}>
							{t('action.skip')}
						</button>
						<button
							bind:this={tombolUtama}
							type="button"
							class="btn btn-primary btn-lg"
							onclick={lanjut}
						>
							{t('tour.start')}
							<ArrowRight size={18} strokeWidth={1.75} aria-hidden="true" />
						</button>
					</div>
				</div>
			</div>
		{:else}
			<div
				bind:offsetWidth={kartu.lebar}
				bind:offsetHeight={kartu.tinggi}
				class="card fixed w-80 p-5 shadow-3"
				style:left="{posisi.x}px"
				style:top="{posisi.y}px"
			>
				<div class="flex items-center justify-between gap-3">
					<p class="num text-caption text-muted">
						{t('tour.progress', { n: indeks, total: daftar.length - 1 })}
					</p>
					<button
						type="button"
						class="btn btn-ghost btn-icon -mr-2 -mt-2 h-9 w-9"
						aria-label={t('action.close')}
						onclick={() => dialog?.close()}
					>
						<X size={18} strokeWidth={1.75} aria-hidden="true" />
					</button>
				</div>
				<ol class="mt-1 flex gap-1" aria-hidden="true">
					{#each daftar.slice(1) as butir, i (butir.judul)}
						<li class="h-1 flex-1 rounded-pill {i < indeks ? 'bg-primary' : 'bg-sunken'}"></li>
					{/each}
				</ol>
				<h2 id="judul-tur" class="t-h3 mt-4">{t(langkah.judul)}</h2>
				<p id="isi-tur" class="mt-1 text-small text-muted">{t(langkah.isi)}</p>
				<div class="mt-5 flex items-center justify-between gap-2">
					<button type="button" class="btn btn-outline btn-sm" onclick={() => (indeks -= 1)}>
						<ArrowLeft size={16} strokeWidth={1.75} aria-hidden="true" />
						{t('action.back')}
					</button>
					<button
						bind:this={tombolUtama}
						type="button"
						class="btn btn-primary btn-sm"
						onclick={lanjut}
					>
						{isTerakhir ? t('tour.done') : t('action.next')}
						{#if !isTerakhir}<ArrowRight size={16} strokeWidth={1.75} aria-hidden="true" />{/if}
					</button>
				</div>
			</div>
		{/if}
	{/if}
</dialog>
