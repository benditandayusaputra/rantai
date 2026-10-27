<script lang="ts">
	import { Landmark } from '@lucide/svelte';
	import { ambilTema, tangkap } from '$lib/api/klien';
	import type { TemaDetail } from '$lib/api/tipe';
	import { t } from '$lib/i18n';

	let { kodeTema }: { kodeTema: string } = $props();
	let tema = $state<TemaDetail | null>(null);
	let isGagal = $state(false);
	$effect(() => {
		let dibatalkan = false;
		tema = null;
		isGagal = false;
		void tangkap(ambilTema(fetch, kodeTema)).then((hasil) => {
			if (dibatalkan) return;
			if (hasil.ok) tema = hasil.data;
			else isGagal = true;
		});
		return () => {
			dibatalkan = true;
		};
	});
</script>

<section
	class="rounded-lg border border-line bg-surface p-4 md:p-5"
	aria-label={t('services.responsibility')}
	aria-busy={!tema && !isGagal}
>
	<h3 class="flex items-center gap-2 font-semibold">
		<Landmark size={18} aria-hidden="true" />{t('services.responsibility')}
	</h3>
	{#if tema?.tanggung_jawab.length}
		<dl class="mt-3 grid gap-4 text-small md:grid-cols-2">
			<div>
				<dt class="text-caption text-muted">{t('services.executor')}</dt>
				<dd class="mt-1 font-semibold">{tema.tanggung_jawab[0]?.instansi_pelaksana}</dd>
				<dd class="mt-1 text-muted">{t('services.executorRole')}</dd>
			</div>
			<div>
				<dt class="text-caption text-muted">{t('services.overseer')}</dt>
				<dd class="mt-1 font-semibold">
					{t('services.commission', { code: tema.tanggung_jawab[0]?.komisi ?? '' })}
				</dd>
				<dd class="mt-1 text-muted">{t('services.overseerRole')}</dd>
			</div>
		</dl>
	{:else}<p class="mt-3 text-small text-muted">
			{t(isGagal ? 'services.responsibilityFailed' : 'state.loading')}
		</p>{/if}
</section>
