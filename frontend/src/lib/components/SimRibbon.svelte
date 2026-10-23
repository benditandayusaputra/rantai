<script lang="ts">
	import { FlaskConical } from '@lucide/svelte';
	import { t } from '$lib/i18n';
</script>

<p
	class="hatch-sim flex h-8 items-center gap-2 border-b-ctrl border-cert-sim px-3 text-caption font-bold uppercase tracking-wide text-cert-sim"
>
	<FlaskConical size={16} strokeWidth={1.75} aria-hidden="true" />
	<span class="rounded-xs bg-cert-sim-bg px-1">{t('sim.ribbon')}</span>
</p>
