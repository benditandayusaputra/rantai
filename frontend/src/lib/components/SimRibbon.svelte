<script lang="ts">
	import { FlaskConical } from '@lucide/svelte';
	import { t } from '$lib/i18n';
</script>

<p
	class="flex min-h-8 items-center gap-2 border-b border-line bg-cert-sim-bg px-3 text-caption font-semibold text-cert-sim"
>
	<FlaskConical size={16} strokeWidth={1.75} aria-hidden="true" />
	<span>{t('sim.ribbon')}</span>
</p>
