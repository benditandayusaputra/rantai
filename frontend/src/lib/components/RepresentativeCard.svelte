<script lang="ts">
	import { ChevronRight, Landmark } from '@lucide/svelte';
	import type { Wakil } from '$lib/api/tipe';
	import { inisial } from '$lib/features/kecamatan/kecamatan';
	import { bahasaBawaan, t } from '$lib/i18n';

	let { wakil }: { wakil: Wakil } = $props();
</script>

<a
	href="/{bahasaBawaan}/anggota/{wakil.id}"
	class="card-link group flex h-full items-start gap-3 p-4"
>
	<span
		class="inline-flex h-11 w-11 shrink-0 items-center justify-center rounded-pill bg-primary font-bold text-primary-on"
		aria-hidden="true">{inisial(wakil.nama)}</span
	>
	<span class="flex min-w-0 flex-1 flex-col gap-1">
		<span class="font-bold">{wakil.nama}</span>
		<span class="text-small text-muted">
			{wakil.fraksi} · {wakil.komisi
				? t('reps.commission', { kode: wakil.komisi.kode })
				: t('reps.noCommission')}
		</span>
		{#if wakil.duduk_di_komisi_pengawas}
			<span
				class="mt-1 inline-flex items-center gap-1 self-start rounded-pill bg-secondary px-3 py-1 text-caption font-bold text-secondary-on"
				><Landmark size={14} strokeWidth={1.75} aria-hidden="true" />{t('acc.oversees')}</span
			>
		{/if}
	</span>
	<ChevronRight
		size={20}
		strokeWidth={1.75}
		aria-hidden="true"
		class="mt-3 shrink-0 text-muted transition-transform duration-fast ease-standard group-hover:translate-x-0.5 group-hover:text-primary"
	/>
</a>
