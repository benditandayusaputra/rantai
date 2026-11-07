<script lang="ts">
	import { ArrowRight, X } from '@lucide/svelte';
	import type { Jaringan, JenisSimpul, SimpulJaringan } from '$lib/api/tipe';
	import ProvenanceBadge from '$lib/components/ProvenanceBadge.svelte';
	import TableScroll from '$lib/components/TableScroll.svelte';
	import TableToggle from '$lib/components/TableToggle.svelte';
	import { tampilkanAngka } from '$lib/features/kecamatan/kecamatan';
	import { bahasaAktif, t } from '$lib/i18n';
	import {
		cariTetangga,
		hitungJari,
		nilaiMaksPerJenis,
		ringkasLabel,
		susunBarisRelasi,
		tautanSimpul
	} from './jaringan';
	import { susunTataLetak } from './tataLetak';

	interface Props {
		jaringan: Jaringan;
		idIsi: string;
	}

	let { jaringan, idIsi }: Props = $props();

	const bidang = { lebar: 960, tinggi: 600, tepi: 32 };
	const batasLabelTetangga = 6;
	const jenisLegenda: JenisSimpul[] = ['instansi', 'paket', 'penyedia'];

	let isTabel = $state(false);
	let terpilih = $state<string | null>(null);

	const posisi = $derived(
		susunTataLetak(
			jaringan.simpul.map((s) => s.id),
			jaringan.sisi,
			bidang
		)
	);
	const nilaiMaks = $derived(nilaiMaksPerJenis(jaringan.simpul));
	const menurutId = $derived(new Map(jaringan.simpul.map((s) => [s.id, s])));
	const tetangga = $derived(terpilih ? cariTetangga(terpilih, jaringan.sisi) : new Set<string>());
	const simpulTerpilih = $derived(terpilih ? (menurutId.get(terpilih) ?? null) : null);
	const daftarTetangga = $derived(
		[...tetangga].flatMap((id) => {
			const s = menurutId.get(id);
			return s ? [s] : [];
		})
	);
	const baris = $derived(susunBarisRelasi(jaringan.simpul, jaringan.sisi));
	const urutanGambar = $derived(
		[...jaringan.simpul].sort((a, b) => Number(a.jenis === 'paket') - Number(b.jenis === 'paket'))
	);

	function isRedup(id: string): boolean {
		return terpilih !== null && id !== terpilih && !tetangga.has(id);
	}

	function isSisiAktif(dari: string, ke: string): boolean {
		return terpilih !== null && (dari === terpilih || ke === terpilih);
	}

	function isLabelTampil(s: SimpulJaringan): boolean {
		if (terpilih === null) return s.jenis === 'instansi';
		return s.id === terpilih || (tetangga.size <= batasLabelTetangga && tetangga.has(s.id));
	}

	function pilih(id: string) {
		terpilih = terpilih === id ? null : id;
	}

	function tanganiTombol(event: KeyboardEvent, id: string) {
		if (event.key !== 'Enter' && event.key !== ' ') return;
		event.preventDefault();
		pilih(id);
	}

	function labelSimpul(s: SimpulJaringan): string {
		return `${t(`network.legend.${s.jenis}`)}: ${s.label}, ${tampilkanAngka(s.nilai).nilai}`;
	}
</script>

{#snippet bentuk(jenis: JenisSimpul, jari: number, kelas: string)}
	{#if jenis === 'instansi'}
		<rect x={-jari} y={-jari} width={jari * 2} height={jari * 2} rx="3" class={kelas} />
	{:else if jenis === 'penyedia'}
		<rect
			x={-jari * 0.85}
			y={-jari * 0.85}
			width={jari * 1.7}
			height={jari * 1.7}
			transform="rotate(45)"
			class={kelas}
		/>
	{:else}
		<circle r={jari} class={kelas} />
	{/if}
{/snippet}

{#snippet legenda()}
	<ul class="flex flex-wrap gap-4 text-small text-muted">
		{#each jenisLegenda as jenis (jenis)}
			<li class="flex items-center gap-2">
				<svg width="20" height="20" viewBox="-10 -10 20 20" aria-hidden="true">
					{@render bentuk(jenis, 6, `simpul-${jenis}`)}
				</svg>
				{t(`network.legend.${jenis}`)}
			</li>
		{/each}
	</ul>
{/snippet}

<div class="flex flex-col gap-4">
	<div class="flex flex-wrap items-center justify-between gap-3">
		{@render legenda()}
		<TableToggle bind:pressed={isTabel} controls={idIsi} labelKembali={t('network.viewAsGraph')} />
	</div>

	<div id={idIsi}>
		{#if isTabel}
			<TableScroll label={t('network.title')}>
				<table class="w-full min-w-max text-small">
					<caption class="sr-only">{t('network.title')}</caption>
					<thead class="bg-sunken text-caption font-semibold">
						<tr>
							<th scope="col" class="p-3 text-left">{t('network.col.agency')}</th>
							<th scope="col" class="p-3 text-left">{t('network.col.package')}</th>
							<th scope="col" class="p-3 text-left">{t('network.col.vendor')}</th>
							<th scope="col" class="p-3 text-right">{t('network.col.value')}</th>
						</tr>
					</thead>
					<tbody>
						{#each baris as b (b.paket.id)}
							<tr class="border-b border-line">
								<td class="p-3">{b.instansi?.label ?? ''}</td>
								<th scope="row" class="max-w-md p-3 text-left font-normal">
									<a href={tautanSimpul(b.paket, bahasaAktif)} class="link">{b.paket.label}</a>
								</th>
								<td class="p-3">
									{#if b.penyedia}<a href={tautanSimpul(b.penyedia, bahasaAktif)} class="link"
											>{b.penyedia.label}</a
										>{:else}<span class="text-muted">{t('pkg.noVendor')}</span>{/if}
								</td>
								<td class="p-3 text-right">
									<ProvenanceBadge
										angka={b.paket.nilai}
										nilai={tampilkanAngka(b.paket.nilai).nilai}
										compact
									/>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</TableScroll>
		{:else}
			<div class="grid gap-4 lg:grid-cols-[1fr_theme(spacing.66)]">
				<div class="card overflow-hidden">
					<svg
						viewBox="0 0 {bidang.lebar} {bidang.tinggi}"
						class="h-auto w-full bg-surface"
						role="group"
						aria-label={t('network.graphLabel', {
							simpul: jaringan.simpul.length,
							sisi: jaringan.sisi.length
						})}
					>
						<g aria-hidden="true">
							{#each jaringan.sisi as s (`${s.dari}>${s.ke}`)}
								{@const a = posisi.get(s.dari)}
								{@const b = posisi.get(s.ke)}
								{#if a && b}
									<line
										x1={a.x}
										y1={a.y}
										x2={b.x}
										y2={b.y}
										class="jaringan-sisi"
										class:jaringan-sisi-aktif={isSisiAktif(s.dari, s.ke)}
										class:jaringan-redup={terpilih !== null && !isSisiAktif(s.dari, s.ke)}
									/>
								{/if}
							{/each}
						</g>
						{#each urutanGambar as s (s.id)}
							{@const p = posisi.get(s.id)}
							{@const jari = hitungJari(s, nilaiMaks)}
							{#if p}
								<g
									transform="translate({p.x} {p.y})"
									role="button"
									tabindex="0"
									aria-pressed={terpilih === s.id}
									aria-label={labelSimpul(s)}
									class="jaringan-simpul cursor-pointer"
									class:jaringan-redup={isRedup(s.id)}
									onclick={() => pilih(s.id)}
									onkeydown={(event) => tanganiTombol(event, s.id)}
								>
									<title>{s.label}</title>
									<circle
										r={jari + 6}
										class="jaringan-cincin"
										class:jaringan-cincin-aktif={terpilih === s.id}
									/>
									{@render bentuk(s.jenis, jari, `simpul-${s.jenis}`)}
									{#if isLabelTampil(s)}
										<text y={jari + 15} text-anchor="middle" class="jaringan-label"
											>{ringkasLabel(s.label)}</text
										>
									{/if}
								</g>
							{/if}
						{/each}
					</svg>
				</div>

				<aside class="card flex flex-col gap-3 p-4" aria-live="polite">
					{#if simpulTerpilih}
						<div class="flex items-start justify-between gap-2">
							<div class="flex min-w-0 flex-col gap-1">
								<p class="t-over">{t(`network.legend.${simpulTerpilih.jenis}`)}</p>
								<p class="break-words font-semibold">{simpulTerpilih.label}</p>
							</div>
							<button
								type="button"
								class="btn btn-ghost min-w-touch shrink-0 px-2"
								aria-label={t('network.reset')}
								onclick={() => (terpilih = null)}
								><X size={18} strokeWidth={1.75} aria-hidden="true" /></button
							>
						</div>
						<ProvenanceBadge
							angka={simpulTerpilih.nilai}
							nilai={tampilkanAngka(simpulTerpilih.nilai).nilai}
						/>
						{#if simpulTerpilih.jenis !== 'paket'}
							<p class="text-small text-muted">
								{t('network.packages', { jumlah: simpulTerpilih.jumlah_paket })}
							</p>
						{/if}
						{#if tautanSimpul(simpulTerpilih, bahasaAktif)}
							<a
								href={tautanSimpul(simpulTerpilih, bahasaAktif)}
								class="btn btn-outline self-start text-small"
								>{t('network.open')}<ArrowRight size={16} aria-hidden="true" /></a
							>
						{/if}
						<p class="mt-2 text-caption font-semibold text-muted">{t('network.linked')}</p>
						<ul class="flex max-h-72 flex-col gap-1 overflow-y-auto">
							{#each daftarTetangga as n (n.id)}
								<li>
									<button
										type="button"
										class="flex min-h-touch w-full items-center gap-2 rounded-md px-2 text-left text-small hover:bg-sunken"
										onclick={() => pilih(n.id)}
									>
										<svg width="16" height="16" viewBox="-8 -8 16 16" aria-hidden="true">
											{@render bentuk(n.jenis, 5, `simpul-${n.jenis}`)}
										</svg>
										<span class="min-w-0 flex-1 truncate">{n.label}</span>
									</button>
								</li>
							{/each}
						</ul>
					{:else}
						<p class="text-small text-muted">{t('network.pick')}</p>
					{/if}
				</aside>
			</div>
		{/if}
	</div>
</div>
