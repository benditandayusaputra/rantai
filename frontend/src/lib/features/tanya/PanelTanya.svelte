<script lang="ts">
	import { Maximize2, MessageCircle, X } from '@lucide/svelte';
	import { bahasaBawaan, t } from '$lib/i18n';
	import PercakapanTanya from './PercakapanTanya.svelte';

	interface Props {
		konteks: { kode: string; nama: string };
		open: boolean;
	}

	let { konteks, open = $bindable() }: Props = $props();

	const idJudul = $props.id();
	let dialog: HTMLDialogElement;

	$effect(() => {
		if (open && !dialog.open) dialog.showModal();
		else if (!open && dialog.open) dialog.close();
	});
</script>

<dialog
	bind:this={dialog}
	onclose={() => (open = false)}
	aria-labelledby={idJudul}
	class="ml-auto mr-0 h-dvh max-h-none w-panel-tanya max-w-full overflow-y-auto bg-surface p-5 text-fg shadow-3 backdrop:bg-scrim"
>
	<div class="flex items-center justify-between gap-2">
		<h2 id={idJudul} class="t-h3 flex items-center gap-2">
			<MessageCircle size={20} strokeWidth={1.75} aria-hidden="true" />{t('ask.title')}
		</h2>
		<div class="flex items-center">
			<a
				href="/{bahasaBawaan}/tanya?w={encodeURIComponent(konteks.kode)}"
				class="inline-flex min-h-touch min-w-touch items-center justify-center rounded-md hover:bg-sunken"
				aria-label={t('ask.fullScreen')}
			>
				<Maximize2 size={20} strokeWidth={1.75} aria-hidden="true" />
			</a>
			<button
				type="button"
				class="inline-flex min-h-touch min-w-touch items-center justify-center rounded-md hover:bg-sunken"
				aria-label={t('action.close')}
				onclick={() => (open = false)}
			>
				<X size={20} strokeWidth={1.75} aria-hidden="true" />
			</button>
		</div>
	</div>
	<div class="mt-4">
		<PercakapanTanya {konteks} />
	</div>
</dialog>
