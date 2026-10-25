<script lang="ts">
	import { Layers } from '@lucide/svelte';
	import { t } from '$lib/i18n';
	import { daftarLapisan, type StatusLapisan } from './lapisan';

	let { lapisan = $bindable(), onchange }: { lapisan: StatusLapisan; onchange: () => void } =
		$props();
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
					class="h-5 w-5 accent-primary"
				/>
				{t(`map.layer.${nama}`)}
			</label>
		{/each}
	</div>
</fieldset>
