<script lang="ts">
	import RadioCard from '$lib/components/RadioCard.svelte';
	import { t } from '$lib/i18n';
	import { draf } from './draf.svelte';
	import { daftarKondisi } from './lapor';
</script>

<fieldset class="flex flex-col gap-3">
	<legend class="mb-2 text-balance font-serif text-h1">{t('report.s4.title')}</legend>
	{#each daftarKondisi as kondisi (kondisi)}
		<RadioCard
			name="kondisi"
			value={kondisi}
			bind:group={draf.kondisi}
			title={t(`report.cond.${kondisi}`)}
		/>
	{/each}
</fieldset>
