<script lang="ts">
	import { Check, Copy, Download, FileText } from '@lucide/svelte';
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
	let isGagalSalin = $state(false);

	async function salin() {
		isGagalSalin = false;
		try {
			await navigator.clipboard.writeText(dokumen.teks);
			isDisalin = true;
			setTimeout(() => (isDisalin = false), durasiTandaSalin);
		} catch {
			isGagalSalin = true;
		}
	}
</script>

<div class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
	<h2 {id} class="t-h3 flex items-center gap-3">
		<span class="icon-tile h-9 w-9"
			><FileText size={18} strokeWidth={1.75} aria-hidden="true" /></span
		>{title}
	</h2>
	<div class="grid grid-cols-2 gap-2 md:flex">
		<button type="button" class="btn btn-outline" onclick={salin}>
			{#if isDisalin}
				<Check size={20} strokeWidth={1.75} aria-hidden="true" />{t('action.copied')}
			{:else}
				<Copy size={20} strokeWidth={1.75} aria-hidden="true" />{t('action.copy')}
			{/if}
		</button>
		<a href={urlServer(dokumen.pdf_url)} download class="btn btn-primary">
			<Download size={20} strokeWidth={1.75} aria-hidden="true" />{t('action.downloadPdf')}
		</a>
	</div>
</div>
<p class="sr-only" aria-live="polite">{isDisalin ? t('action.copied') : ''}</p>
{#if isGagalSalin}<p role="alert" class="mt-3 text-small text-danger">
		{t('packet.copyFailed')}
	</p>{/if}
<pre
	class="mt-4 overflow-x-auto whitespace-pre-wrap rounded-lg border border-line bg-raised p-5 font-sans text-body shadow-2 md:p-8">{dokumen.teks}</pre>
