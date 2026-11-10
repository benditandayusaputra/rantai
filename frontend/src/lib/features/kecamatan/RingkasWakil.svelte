<script lang="ts">
	import { Landmark } from '@lucide/svelte';
	import type { WakilKecamatan } from '$lib/api/tipe';
	import FotoWakil from '$lib/components/FotoWakil.svelte';
	import { bahasaAktif, t } from '$lib/i18n';
	import { batasWakilRingkas } from './kecamatan';

	let { wakilku }: { wakilku: WakilKecamatan } = $props();

	const tampil = $derived(wakilku.dapil.anggota.slice(0, batasWakilRingkas));
</script>

<p class="text-small text-muted">{t('acc.reps', { dapil: wakilku.dapil.nomor })}</p>
<ul class="mt-4 grid grid-cols-2 gap-3 md:grid-cols-4 md:gap-4">
	{#each tampil as wakil (wakil.id)}
		<li>
			<a
				href="/{bahasaAktif}/anggota/{wakil.id}?w={encodeURIComponent(wakilku.kecamatan.kode)}"
				class="card-link flex h-full flex-col items-center gap-3 p-4 text-center"
			>
				<FotoWakil
					nama={wakil.nama}
					foto={wakil.foto_url}
					class="h-24 w-24 text-h3 md:h-28 md:w-28"
				/>
				<span class="flex flex-col gap-1">
					<span class="font-bold">{wakil.nama}</span>
					<span class="text-small text-muted">
						{wakil.fraksi} · {wakil.komisi
							? t('reps.commission', { kode: wakil.komisi.kode })
							: t('reps.noCommission')}
					</span>
				</span>
				{#if wakil.duduk_di_komisi_pengawas}
					<span
						class="mt-auto inline-flex items-center gap-1 rounded-xs bg-secondary px-2 py-0.5 text-caption font-bold text-secondary-on"
						><Landmark size={14} strokeWidth={1.75} aria-hidden="true" />{t('acc.oversees')}</span
					>
				{/if}
			</a>
		</li>
	{/each}
</ul>
<p class="mt-3 text-caption text-muted">{t('photo.dprd')} {t('reps.noScore')}</p>
