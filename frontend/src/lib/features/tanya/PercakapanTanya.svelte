<script lang="ts">
	import { ExternalLink, Sparkles } from '@lucide/svelte';
	import { tanyaRantai, type JawabanTanya } from '$lib/api/tanya';
	import type { Provenance } from '$lib/api/tipe';
	import ChatComposer from '$lib/components/ChatComposer.svelte';
	import ChatMessage from '$lib/components/ChatMessage.svelte';
	import Citation from '$lib/components/Citation.svelte';
	import SourcePanel from '$lib/components/SourcePanel.svelte';
	import { muatProvenance, namaSumber } from '$lib/features/provenance/provenance';
	import { bahasaAktif, t } from '$lib/i18n';
	import {
		pecahJawaban,
		panjangMaksPertanyaan,
		saranPertanyaan,
		siapkanTautan,
		type TautanSiap
	} from './tanya';

	interface Props {
		konteks?: { kode: string; nama: string } | null;
		onlepas?: () => void;
	}

	let { konteks = null, onlepas }: Props = $props();

	interface Pesan {
		id: number;
		peran: 'user' | 'assistant';
		teks: string;
		status: 'streaming' | 'done' | 'error';
		jumlahAlat: number;
		rujukan: Map<number, Provenance>;
		tautan: TautanSiap[];
	}

	let pesan = $state<Pesan[]>([]);
	let masukan = $state('');
	let panelTerbuka = $state(false);
	let provenanceTerpilih = $state<Provenance | null>(null);
	const isMenjawab = $derived(pesan.at(-1)?.status === 'streaming');
	let urutan = 0;

	async function muatRujukan(jawaban: JawabanTanya): Promise<Map<number, Provenance>> {
		const daftar = await Promise.all(
			jawaban.rujukan.map(async (r) => [r.nomor, await muatProvenance(r.provenance_id)] as const)
		);
		return new Map(daftar);
	}

	function tambahPesan(peran: Pesan['peran'], teks: string, status: Pesan['status']): number {
		const id = urutan++;
		pesan.push({ id, peran, teks, status, jumlahAlat: 0, rujukan: new Map(), tautan: [] });
		return id;
	}

	function ubahPesan(id: number, ubah: (p: Pesan) => Partial<Pesan>) {
		pesan = pesan.map((p) => (p.id === id ? { ...p, ...ubah(p) } : p));
	}

	async function tanya(pertanyaan: string) {
		masukan = '';
		tambahPesan('user', pertanyaan, 'done');
		const id = tambahPesan('assistant', '', 'streaming');
		try {
			const jawaban = await tanyaRantai(
				fetch,
				{ pertanyaan, kode_kecamatan_konteks: konteks?.kode ?? '', bahasa: bahasaAktif },
				() => ubahPesan(id, (p) => ({ jumlahAlat: p.jumlahAlat + 1 }))
			);
			const rujukan = await muatRujukan(jawaban);
			const tautan = siapkanTautan(jawaban.tautan ?? [], bahasaAktif);
			ubahPesan(id, () => ({ teks: jawaban.teks, rujukan, tautan, status: 'done' }));
		} catch {
			ubahPesan(id, () => ({ teks: t('ask.error'), status: 'error' }));
		}
	}

	function bukaSumber(provenance: Provenance) {
		provenanceTerpilih = provenance;
		panelTerbuka = true;
	}
</script>

<div class="flex flex-col gap-6">
	{#if pesan.length === 0}
		<div class="flex items-start gap-3">
			<span class="icon-tile"><Sparkles size={20} strokeWidth={1.75} aria-hidden="true" /></span>
			<div>
				<p class="t-h3">{t('ask.emptyTitle')}</p>
				<p class="mt-1 max-w-prose text-muted">{t('ask.emptyLead')}</p>
			</div>
		</div>
	{/if}
	<div role="log" aria-live="polite" aria-label={t('ask.messages')} class="flex flex-col gap-4">
		{#each pesan as p (p.id)}
			<ChatMessage role={p.peran} state={p.status}>
				{#if p.status === 'streaming'}
					<span class="text-muted">
						{p.jumlahAlat ? t('ask.streaming', { n: p.jumlahAlat }) : t('ask.thinking')}
					</span>
				{:else}
					<p class="whitespace-pre-line">
						{#each pecahJawaban(p.teks, [...p.rujukan.keys()], p.tautan) as segmen, i (i)}
							{#if segmen.jenis === 'teks'}{segmen.isi}{:else if segmen.jenis === 'tautan'}<a
									href={segmen.href}
									class="link"
									target={segmen.isLuar ? '_blank' : undefined}
									rel={segmen.isLuar ? 'noopener' : undefined}
									>{segmen.isi}{#if segmen.isLuar}<ExternalLink
											size={14}
											class="ml-0.5 inline align-baseline"
											aria-hidden="true"
										/><span class="sr-only">{t('a11y.opensNewTab')}</span>{/if}</a
								>{:else}
								{@const provenance = p.rujukan.get(segmen.nomor)}
								{#if provenance}
									<Citation
										n={segmen.nomor}
										level={provenance.label}
										label={namaSumber(provenance.sumber)}
										onopen={() => bukaSumber(provenance)}
									/>
								{/if}
							{/if}
						{/each}
					</p>
				{/if}
			</ChatMessage>
		{/each}
	</div>
	<ChatComposer
		bind:value={masukan}
		context={konteks ? { label: konteks.nama, onremove: onlepas } : undefined}
		suggestions={pesan.length === 0 ? saranPertanyaan() : []}
		disabled={isMenjawab || !konteks}
		placeholder={konteks ? undefined : t('ask.needArea')}
		maxLength={panjangMaksPertanyaan}
		onsubmit={tanya}
	/>
</div>
<SourcePanel bind:open={panelTerbuka} provenance={provenanceTerpilih} isGagal={false} />
