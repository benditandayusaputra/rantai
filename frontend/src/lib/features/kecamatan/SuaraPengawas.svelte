<script lang="ts">
	import { ChevronRight, Landmark } from '@lucide/svelte';
	import { ambilWakilDenganUmpan, tangkap } from '$lib/api/klien';
	import type { WakilDenganUmpan } from '$lib/api/tipe';
	import KartuUmpan from '$lib/features/anggota/KartuUmpan.svelte';
	import { bahasaBawaan, t } from '$lib/i18n';
	import { inisial } from './kecamatan';

	let { kodeKecamatan, kodeTema }: { kodeKecamatan: string; kodeTema: string } = $props();

	let data = $state<WakilDenganUmpan | null>(null);
	let isGagal = $state(false);

	$effect(() => {
		void tangkap(ambilWakilDenganUmpan(fetch, kodeKecamatan, kodeTema)).then((hasil) => {
			if (hasil.ok) data = hasil.data;
			else isGagal = true;
		});
	});

	const pengawas = $derived(data?.dapil.anggota.filter((a) => a.duduk_di_komisi_pengawas) ?? []);
</script>

{#if isGagal}
	<p class="text-small text-muted">{t('voice.failed')}</p>
{:else if data && pengawas.length}
	<section aria-labelledby="judul-suara-{kodeTema}" class="flex flex-col gap-3">
		<h4 id="judul-suara-{kodeTema}" class="flex items-center gap-2 text-h3 font-bold">
			<Landmark size={20} strokeWidth={1.75} aria-hidden="true" class="text-primary" />
			{t('voice.title')}
		</h4>
		<p class="max-w-prose text-small text-muted">{t('voice.lead')}</p>
		<ul class="grid gap-4 md:grid-cols-2">
			{#each pengawas as wakil (wakil.id)}
				{@const umpan = data.umpan_isu_pengawas?.[wakil.id] ?? []}
				<li class="card flex flex-col gap-3 p-4">
					<a
						href="/{bahasaBawaan}/anggota/{wakil.id}"
						class="group flex items-center gap-3 font-bold underline-offset-4 hover:underline"
					>
						<span
							class="inline-flex h-10 w-10 shrink-0 items-center justify-center rounded-pill bg-secondary text-small text-secondary-on"
							aria-hidden="true">{inisial(wakil.nama)}</span
						>
						<span class="min-w-0 flex-1">{wakil.nama}</span>
						<ChevronRight
							size={18}
							strokeWidth={1.75}
							aria-hidden="true"
							class="shrink-0 text-muted transition-transform duration-fast ease-standard group-hover:translate-x-0.5 group-hover:text-primary"
						/>
					</a>
					{#if umpan.length}
						<ul class="flex flex-col gap-2">
							{#each umpan as item (item.id)}<KartuUmpan {item} />{/each}
						</ul>
					{:else}
						<p class="rounded-md bg-sunken px-3 py-2 text-small text-muted">{t('voice.none')}</p>
					{/if}
				</li>
			{/each}
		</ul>
	</section>
{/if}
