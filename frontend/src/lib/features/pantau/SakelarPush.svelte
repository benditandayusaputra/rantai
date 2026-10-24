<script lang="ts">
	import { onMount } from 'svelte';
	import { tangkap } from '$lib/api/klien';
	import { ambilKunciPublikPush } from '$lib/api/webpush';
	import { t } from '$lib/i18n';
	import { nyalakanPush, padamkanPush, statusPushAwal, type StatusPush } from './push';

	let { token }: { token: string } = $props();

	let kunciPublik = $state('');
	let status = $state<StatusPush>('mati');
	let isSibuk = $state(false);

	onMount(async () => {
		const kunci = await tangkap(ambilKunciPublikPush(fetch));
		if (!kunci.ok) return;
		kunciPublik = kunci.data.kunci_publik;
		status = await statusPushAwal();
	});

	async function ubah() {
		isSibuk = true;
		try {
			status =
				status === 'aktif'
					? await padamkanPush(fetch, token)
					: await nyalakanPush(fetch, token, kunciPublik);
		} catch {
			status = 'gagal';
		}
		isSibuk = false;
	}

	const pesan = $derived(
		{
			tidak_didukung: t('watch.push.unsupported'),
			ditolak: t('watch.push.denied'),
			gagal: t('watch.push.failed')
		}[status as 'tidak_didukung' | 'ditolak' | 'gagal'] ?? ''
	);
</script>

{#if kunciPublik}
	<section class="flex flex-col gap-2 rounded-lg border border-line bg-surface p-4">
		<h2 class="font-bold">{t('watch.push.title')}</h2>
		<label class="flex min-h-touch cursor-pointer items-center justify-between gap-3">
			<span>{t('watch.push.label')}</span>
			<input
				type="checkbox"
				role="switch"
				class="h-6 w-6 accent-primary"
				checked={status === 'aktif'}
				disabled={isSibuk || status === 'tidak_didukung' || status === 'ditolak'}
				aria-describedby="bantuan-push"
				onchange={ubah}
			/>
		</label>
		<p id="bantuan-push" class="text-small text-muted">{pesan || t('watch.push.help')}</p>
	</section>
{/if}
