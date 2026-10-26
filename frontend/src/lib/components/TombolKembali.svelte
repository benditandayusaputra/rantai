<script lang="ts">
	import { ArrowLeft } from '@lucide/svelte';
	import { afterNavigate } from '$app/navigation';
	import { bahasaAktif, t } from '$lib/i18n';

	let { cadangan = `/${bahasaAktif}` }: { cadangan?: string } = $props();

	let isDatangDariHalamanRantai = $state(false);

	afterNavigate(({ from }) => {
		isDatangDariHalamanRantai = from !== null;
	});

	function isKlikBiasa(e: MouseEvent): boolean {
		return e.button === 0 && !e.metaKey && !e.ctrlKey && !e.shiftKey && !e.altKey;
	}

	function kembali(e: MouseEvent) {
		if (!isDatangDariHalamanRantai || !isKlikBiasa(e)) return;
		e.preventDefault();
		history.back();
	}
</script>

<a href={cadangan} onclick={kembali} class="btn btn-ghost btn-sm -ml-3 self-start">
	<ArrowLeft size={18} strokeWidth={1.75} aria-hidden="true" />
	{t('action.back')}
</a>
