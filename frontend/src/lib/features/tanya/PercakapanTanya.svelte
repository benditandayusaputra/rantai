<script lang="ts">
	import { Sparkles } from '@lucide/svelte';
	import { tanyaRantai, type JawabanTanya } from '$lib/api/tanya';
	import type { Provenance } from '$lib/api/tipe';
	import ChatComposer from '$lib/components/ChatComposer.svelte';
	import ChatMessage from '$lib/components/ChatMessage.svelte';
	import Citation from '$lib/components/Citation.svelte';
	import SourcePanel from '$lib/components/SourcePanel.svelte';
	import { muatProvenance, namaSumber } from '$lib/features/provenance/provenance';
	import { bahasaBawaan, t } from '$lib/i18n';
	import { pecahRujukan, panjangMaksPertanyaan, saranPertanyaan } from './tanya';

	interface Props {
		konteks?: { kode: string; nama: string } | null;
	}

	let { konteks = $bindable(null) }: Props = $props();

	interface Pesan {
		id: number;
		peran: 'user' | 'assistant';
		teks: string;
		status: 'streaming' | 'done' | 'error';
		jumlahAlat: number;
		rujukan: Map<number, Provenance>;
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
		pesan.push({ id, peran, teks, status, jumlahAlat: 0, rujukan: new Map() });
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
				{ pertanyaan, kode_kecamatan_konteks: konteks?.kode ?? '', bahasa: bahasaBawaan },
				() => ubahPesan(id, (p) => ({ jumlahAlat: p.jumlahAlat + 1 }))
			);
			const rujukan = await muatRujukan(jawaban);
			ubahPesan(id, () => ({ teks: jawaban.teks, rujukan, status: 'done' }));
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
						{#each pecahRujukan(p.teks, [...p.rujukan.keys()]) as segmen, i (i)}
							{#if segmen.jenis === 'teks'}{segmen.isi}{:else}
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
		context={konteks ? { label: konteks.nama, onremove: () => (konteks = null) } : undefined}
		suggestions={pesan.length === 0 ? saranPertanyaan() : []}
		disabled={isMenjawab}
		maxLength={panjangMaksPertanyaan}
		onsubmit={tanya}
	/>
</div>
<SourcePanel bind:open={panelTerbuka} provenance={provenanceTerpilih} isGagal={false} />
