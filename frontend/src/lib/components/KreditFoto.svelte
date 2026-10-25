<script lang="ts">
	import type { FotoSekolah } from '$lib/api/tipe';
	import { t } from '$lib/i18n';

	let { foto }: { foto: FotoSekolah } = $props();

	const sumber = $derived(t(`photo.source.${foto.sumber}`));
</script>

<a
	href={foto.halaman}
	target="_blank"
	rel="noopener"
	class="text-caption font-normal text-muted underline-offset-2 hover:text-fg hover:underline"
>
	{foto.kredit
		? t('photo.creditBy', { kredit: foto.kredit, sumber })
		: t('photo.credit', { sumber })}<span class="sr-only"> {t('a11y.opensNewTab')}</span>
</a>
