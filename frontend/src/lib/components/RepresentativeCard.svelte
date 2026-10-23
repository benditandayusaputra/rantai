<script lang="ts">
	import type { Wakil } from '$lib/api/tipe';
	import { inisial } from '$lib/features/kecamatan/kecamatan';
	import { bahasaBawaan, t } from '$lib/i18n';

	let { wakil }: { wakil: Wakil } = $props();
</script>

<a
	href="/{bahasaBawaan}/anggota/{wakil.id}"
	class="flex items-start gap-3 rounded-lg border border-line bg-surface p-4 shadow-1 transition-colors duration-fast ease-standard hover:border-line-strong"
>
	<span
		class="inline-flex h-11 w-11 shrink-0 items-center justify-center rounded-pill bg-secondary font-bold text-secondary-on"
		aria-hidden="true">{inisial(wakil.nama)}</span
	>
	<span class="flex min-w-0 flex-col gap-1">
		<span class="font-bold">{wakil.nama}</span>
		<span class="text-small text-muted">
			{wakil.fraksi} · {wakil.komisi
				? t('reps.commission', { kode: wakil.komisi.kode })
				: t('reps.noCommission')}
		</span>
		{#if wakil.duduk_di_komisi_pengawas}
			<span
				class="mt-1 self-start rounded-pill bg-secondary px-3 py-1 text-caption font-bold text-secondary-on"
				>{t('acc.oversees')}</span
			>
		{/if}
	</span>
</a>
