<script lang="ts">
	import { Camera, X } from '@lucide/svelte';
	import { formatJam } from '$lib/format';
	import { t } from '$lib/i18n';

	interface Foto {
		id: string;
		src: string;
		takenAt: Date;
	}

	interface Props {
		id: string;
		photos: Foto[];
		max?: number;
		busy?: boolean;
		onadd: (berkas: File[]) => void;
		onremove: (id: string) => void;
	}

	let { id, photos, max = 4, busy = false, onadd, onremove }: Props = $props();

	const slotKosong = $derived(Math.max(0, max - photos.length));

	function pilih(event: Event & { currentTarget: HTMLInputElement }) {
		const berkas = Array.from(event.currentTarget.files ?? []).slice(0, slotKosong);
		event.currentTarget.value = '';
		if (berkas.length > 0) onadd(berkas);
	}
</script>

<ul class="grid grid-cols-2 gap-3">
	{#each photos as foto, i (foto.id)}
		<li class="relative aspect-square overflow-hidden rounded-lg bg-sunken">
			<img
				src={foto.src}
				alt={t('report.s3.alt', { n: i + 1, jam: formatJam(foto.takenAt) })}
				class="h-full w-full object-cover"
			/>
			<span class="absolute bottom-2 left-2 rounded-xs bg-surface px-1 text-caption num"
				>{formatJam(foto.takenAt)}</span
			>
			<button
				type="button"
				class="absolute right-1 top-1 inline-flex h-11 w-11 items-center justify-center rounded-pill bg-surface shadow-1"
				aria-label={t('report.s3.remove', { n: i + 1 })}
				onclick={() => onremove(foto.id)}
			>
				<X size={20} strokeWidth={1.75} aria-hidden="true" />
			</button>
		</li>
	{/each}
	{#each Array.from({ length: slotKosong }, (_, i) => i) as urutan (urutan)}
		<li class="aspect-square">
			{#if urutan === 0}
				<label
					for={id}
					class="flex h-full cursor-pointer flex-col items-center justify-center gap-2 rounded-lg border-ctrl border-dashed border-line-strong text-center font-semibold hover:border-fg has-[:focus-visible]:shadow-ring"
					aria-busy={busy}
				>
					<Camera size={28} strokeWidth={1.75} aria-hidden="true" />
					{busy ? t('report.s3.processing') : t('report.s3.take')}
					<input
						{id}
						type="file"
						accept="image/*"
						capture="environment"
						multiple
						class="sr-only"
						disabled={busy}
						onchange={pilih}
					/>
				</label>
			{:else}
				<span
					class="flex h-full items-center justify-center rounded-lg border-ctrl border-dashed border-line text-line-strong"
					aria-hidden="true"><Camera size={24} strokeWidth={1.5} /></span
				>
			{/if}
		</li>
	{/each}
</ul>
