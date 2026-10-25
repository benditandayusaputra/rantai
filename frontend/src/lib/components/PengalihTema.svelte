<script lang="ts">
	import { Moon, Sun } from '@lucide/svelte';
	import { pilihTema, temaLawan, temaTerpasang } from '$lib/features/tema/tema';
	import { t } from '$lib/i18n';

	let { class: kelas = '', denganLabel = false }: { class?: string; denganLabel?: boolean } =
		$props();
</script>

<button
	type="button"
	class="btn btn-ghost {denganLabel ? 'btn-sm' : 'btn-icon'} {kelas}"
	title={t('theme.toggle')}
	onclick={() => pilihTema(temaLawan(temaTerpasang()))}
>
	<Moon size={20} strokeWidth={1.75} aria-hidden="true" class="dark:hidden" />
	<Sun size={20} strokeWidth={1.75} aria-hidden="true" class="hidden dark:block" />
	<span class={denganLabel ? '' : 'sr-only'}>
		<span class="dark:hidden">{t('theme.dark')}</span>
		<span class="hidden dark:inline">{t('theme.light')}</span>
	</span>
</button>
