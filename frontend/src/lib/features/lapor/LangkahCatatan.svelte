<script lang="ts">
	import { Info } from '@lucide/svelte';
	import TextArea from '$lib/components/TextArea.svelte';
	import { t } from '$lib/i18n';
	import { draf } from './draf.svelte';
	import { batasCatatan } from './lapor';
</script>

<div class="flex flex-col gap-4">
	<h1 class="text-balance font-serif text-h1">{t('report.s5.title')}</h1>
	<p class="flex gap-2 rounded-md bg-attn-bg p-3 text-small">
		<Info size={18} strokeWidth={1.75} aria-hidden="true" class="mt-0.5 shrink-0 text-attn" />
		{t('report.s5.rules')}
	</p>
	<TextArea
		id="catatan-laporan"
		label={t('report.s5.label')}
		help={t('report.s5.hint')}
		maxlength={batasCatatan}
		bind:value={draf.catatan}
	/>
</div>
