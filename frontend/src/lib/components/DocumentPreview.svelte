<script lang="ts">
	import { Check, Copy, Download } from '@lucide/svelte';
	import { urlServer } from '$lib/api/klien';
	import type { Dokumen } from '$lib/api/tipe';
	import { t } from '$lib/i18n';

	interface Props {
		id: string;
		title: string;
		dokumen: Dokumen;
	}

	let { id, title, dokumen }: Props = $props();

	const durasiTandaSalin = 3000;
	let isDisalin = $state(false);

	async function salin() {
		await navigator.clipboard.writeText(dokumen.teks);
		isDisalin = true;
		setTimeout(() => (isDisalin = false), durasiTandaSalin);
	}

	const kelasTombol =
		'inline-flex min-h-touch items-center justify-center gap-2 rounded-md px-4 font-semibold';
</script>

<div class="flex flex-wrap items-center justify-between gap-3">
	<h2 {id} class="t-h3">{title}</h2>
	<div class="flex flex-wrap gap-2">
		<button
			type="button"
			class="{kelasTombol} border-ctrl border-line-strong hover:bg-sunken"
			onclick={salin}
		>
			{#if isDisalin}
				<Check size={20} strokeWidth={1.75} aria-hidden="true" />{t('action.copied')}
			{:else}
				<Copy size={20} strokeWidth={1.75} aria-hidden="true" />{t('action.copy')}
			{/if}
		</button>
		<a href={urlServer(dokumen.pdf_url)} download class="{kelasTombol} bg-primary text-primary-on">
			<Download size={20} strokeWidth={1.75} aria-hidden="true" />{t('action.downloadPdf')}
		</a>
	</div>
</div>
<p class="sr-only" aria-live="polite">{isDisalin ? t('action.copied') : ''}</p>
<pre
	class="mt-4 overflow-x-auto whitespace-pre-wrap rounded-lg border border-line bg-raised p-5 font-sans text-body shadow-1">{dokumen.teks}</pre>
