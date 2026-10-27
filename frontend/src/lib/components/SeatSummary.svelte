<script lang="ts">
	import type { IndikatorRingkas } from '$lib/api/tipe';
	import { isTemaPendidikan, teksJenjang } from '$lib/features/pendidikan/jenjang';
	import { formatAngka } from '$lib/format';
	import { bahasaAktif, t } from '$lib/i18n';

	let { indikator, tema }: { indikator: IndikatorRingkas; tema: string } = $props();
	const kursiPerSeratus = $derived(
		indikator.kebutuhan.nilai > 0 && Number.isFinite(indikator.kapasitas.nilai)
			? Math.round((indikator.kapasitas.nilai / indikator.kebutuhan.nilai) * 100)
			: null
	);
</script>

{#if isTemaPendidikan(tema) && kursiPerSeratus !== null}
	<div class="border-l-link border-fg pl-4">
		<p class="text-caption font-semibold text-muted">{t('services.meaning')}</p>
		<p class="mt-1 max-w-prose font-heading text-h3 font-semibold">
			{t('services.perHundred', {
				seats: formatAngka(kursiPerSeratus, bahasaAktif),
				level: t(teksJenjang[tema].jenjang),
				graduates: t(tema === 'pendidikan' ? 'services.primary' : 'services.junior')
			})}
		</p>
		<p class="mt-2 max-w-prose text-small text-muted">{t('services.zoningNote')}</p>
	</div>
{/if}
