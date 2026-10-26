<script lang="ts">
	import { Map as IkonPeta } from '@lucide/svelte';
	import { t } from '$lib/i18n';
	import { contohUbin, daftarPetaDasar, isPetaDasarRaster, type PetaDasar } from './gaya';

	let {
		petaDasar = $bindable(),
		onchange,
		nama = 'peta-dasar'
	}: { petaDasar: PetaDasar; onchange: () => void; nama?: string } = $props();
</script>

<fieldset>
	<legend class="flex items-center gap-2 text-small font-bold">
		<IkonPeta size={16} strokeWidth={1.75} aria-hidden="true" class="text-muted" />
		{t('map.basemap')}
	</legend>
	<div class="mt-2 grid grid-cols-3 gap-2">
		{#each daftarPetaDasar as dasar (dasar)}
			<label
				class="flex cursor-pointer flex-col gap-1 rounded-md border-ctrl border-line p-1 text-caption transition-colors duration-fast ease-standard hover:border-line-strong has-[:checked]:border-primary has-[:checked]:font-semibold has-[:focus-visible]:ring-2 has-[:focus-visible]:ring-focus"
			>
				<input
					type="radio"
					name={nama}
					value={dasar}
					bind:group={petaDasar}
					{onchange}
					class="sr-only"
				/>
				{#if isPetaDasarRaster(dasar)}
					<img
						src={contohUbin(dasar)}
						alt=""
						loading="lazy"
						width="256"
						height="256"
						class="h-10 w-full rounded-sm object-cover"
					/>
				{:else}
					<span class="contoh-peta-standar h-10 w-full rounded-sm" aria-hidden="true"></span>
				{/if}
				<span class="text-center">{t(`map.basemap.${dasar}`)}</span>
			</label>
		{/each}
	</div>
</fieldset>
