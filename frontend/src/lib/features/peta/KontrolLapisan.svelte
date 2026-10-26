<script lang="ts">
	import { Camera, HardHat, Layers, School } from '@lucide/svelte';
	import { t } from '$lib/i18n';
	import { daftarLapisan, type Lapisan, type StatusLapisan } from './lapisan';

	let { lapisan = $bindable(), onchange }: { lapisan: StatusLapisan; onchange: () => void } =
		$props();

	const ikonLapisan: Record<Exclude<Lapisan, 'tema'>, typeof School> = {
		proyek: HardHat,
		sekolah: School,
		laporan: Camera
	};
	const skalaKesenjangan = ['bg-gap-1', 'bg-gap-2', 'bg-gap-3', 'bg-gap-4'];
</script>

<fieldset>
	<legend class="flex items-center gap-2 text-small font-bold">
		<Layers size={16} strokeWidth={1.75} aria-hidden="true" class="text-muted" />
		{t('map.layers')}
	</legend>
	<div class="mt-2 flex flex-col">
		{#each daftarLapisan as nama (nama)}
			<label
				class="-mx-2 flex min-h-touch cursor-pointer items-center gap-3 rounded-md px-2 text-small hover:bg-sunken"
			>
				<input
					type="checkbox"
					bind:checked={lapisan[nama]}
					{onchange}
					class="h-5 w-5 shrink-0 accent-primary"
				/>
				{#if nama === 'tema'}
					<span class="flex h-3.5 w-5 shrink-0 overflow-hidden rounded-xs" aria-hidden="true">
						{#each skalaKesenjangan as kelas (kelas)}<span class="flex-1 {kelas}"></span>{/each}
					</span>
				{:else}
					{@const Ikon = ikonLapisan[nama]}
					<Ikon size={18} strokeWidth={1.75} aria-hidden="true" class="shrink-0 text-muted" />
				{/if}
				{t(`map.layer.${nama}`)}
			</label>
		{/each}
	</div>
</fieldset>
