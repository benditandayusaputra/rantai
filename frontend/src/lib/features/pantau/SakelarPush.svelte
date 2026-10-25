<script lang="ts">
	import { BellRing } from '@lucide/svelte';
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
	<section class="card flex flex-col gap-2 p-5">
		<h2 class="flex items-center gap-2 font-bold">
			<BellRing size={18} strokeWidth={1.75} aria-hidden="true" class="text-muted" />
			{t('watch.push.title')}
		</h2>
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
