<script lang="ts">
	import {
		Camera,
		Check,
		Clock,
		EqualApproximately,
		FileText,
		Gavel,
		Hammer,
		Handshake,
		Hourglass,
		MapPin,
		Trophy,
		Unlink2
	} from '@lucide/svelte';
	import type { KodeMata, MataRantai, StatusMata } from '$lib/api/tipe';
	import { kunciTahap, ringkasRantai } from '$lib/features/paket/paket';
	import { t } from '$lib/i18n';

	let { rantai }: { rantai: MataRantai[] } = $props();

	const ikonTahap: Record<KodeMata, typeof Camera> = {
		rencana: FileText,
		pemilihan: Gavel,
		pemenang: Trophy,
		pelaksanaan: Hammer,
		lokasi: MapPin,
		serah_terima: Handshake,
		masa_pemeliharaan: Hourglass,
		bukti_warga: Camera
	};

	const tampilanStatus: Record<StatusMata, { ikon: typeof Camera; mata: string; teks: string }> = {
		ada: {
			ikon: Check,
			mata: 'border-solid border-primary bg-secondary text-primary',
			teks: 'text-primary'
		},
		putus: {
			ikon: Unlink2,
			mata: 'border-dashed border-attn-fill bg-attn-bg text-attn',
			teks: 'text-attn'
		},
		menunggu: {
			ikon: Clock,
			mata: 'border-dotted border-line-strong text-muted',
			teks: 'text-muted'
		},
		estimasi: {
			ikon: EqualApproximately,
			mata: 'border-dashed border-cert-estimasi text-cert-estimasi',
			teks: 'text-cert-estimasi'
		}
	};

	function isPutusDiSekitar(indeks: number): boolean {
		return rantai[indeks]?.status === 'putus' || rantai[indeks + 1]?.status === 'putus';
	}
</script>

<ol aria-label={ringkasRantai(rantai)} class="flex flex-col lg:flex-row">
	{#each rantai as mata, indeks (mata.kode)}
		{@const IkonTahap = ikonTahap[mata.kode]}
		{@const status = tampilanStatus[mata.status]}
		{@const IkonStatus = status.ikon}
		<li
			class="relative flex gap-3 pb-6 lg:flex-1 lg:flex-col lg:items-center lg:pb-0 lg:text-center"
		>
			{#if indeks < rantai.length - 1}
				<span
					aria-hidden="true"
					class="absolute left-4 top-12 h-full border-l-link lg:left-1/2 lg:top-5 lg:h-0 lg:w-full lg:border-l-0 lg:border-t-link {isPutusDiSekitar(
						indeks
					)
						? 'border-dashed border-attn-fill opacity-60'
						: rantai[indeks + 1]?.status === 'menunggu'
							? 'border-dotted border-line-strong'
							: 'border-primary'}"
				></span>
			{/if}
			<span
				class="relative z-raised inline-flex h-11.5 w-8.5 shrink-0 items-center justify-center rounded-pill border-link bg-surface lg:h-10 lg:w-17 {status.mata}"
			>
				<IkonTahap size={20} strokeWidth={1.75} aria-hidden="true" />
				<span
					class="absolute -right-2 -top-2 inline-flex h-5 w-5 items-center justify-center rounded-pill border-ctrl bg-raised {status.mata}"
				>
					<IkonStatus size={12} strokeWidth={2} aria-hidden="true" />
				</span>
			</span>
			<span class="flex flex-col gap-1 pt-1 lg:px-1 lg:pt-3">
				<span class="font-bold">{t(kunciTahap[mata.kode])}</span>
				<span class="text-caption font-bold uppercase {status.teks}"
					>{t(`chain.status.${mata.status}`)}</span
				>
				{#if mata.alasan}<span class="text-small text-muted">{mata.alasan}</span>{/if}
			</span>
		</li>
	{/each}
</ol>
