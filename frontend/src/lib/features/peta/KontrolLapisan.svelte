<script lang="ts">
	import { Layers } from '@lucide/svelte';
	import { t } from '$lib/i18n';
	import { daftarLapisan, type StatusLapisan } from './lapisan';

	let { lapisan = $bindable(), onchange }: { lapisan: StatusLapisan; onchange: () => void } =
		$props();
</script>

<fieldset>
	<legend class="flex items-center gap-2 text-small font-bold">
		<Layers size={16} strokeWidth={1.75} aria-hidden="true" />
		{t('map.layers')}
	</legend>
	<div class="mt-2 flex flex-col">
		{#each daftarLapisan as nama (nama)}
			<label class="flex min-h-touch cursor-pointer items-center gap-3 text-small">
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
