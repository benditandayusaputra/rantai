<script lang="ts">
	import { CircleCheck, ExternalLink, Landmark } from '@lucide/svelte';
	import { ambilData, tangkap } from '$lib/api/klien';
	import type { Teks } from '$lib/api/tipe';
	import { bahasaAktif, t } from '$lib/i18n';

	interface Jalur {
		kode: string;
		nama: Teks;
		cara: Teks;
		url: string;
	}

	let { kodeWilayah, penerima }: { kodeWilayah: string; penerima: string } = $props();
	let jalur = $state<Jalur[]>([]);
	let status = $state<'memuat' | 'siap' | 'gagal'>('memuat');

	$effect(() => {
		let dibatalkan = false;
		status = 'memuat';
		void tangkap(
			ambilData<{ jalur: Jalur[] }>(fetch, '/jalur-aspirasi', { kode_wilayah: kodeWilayah })
		).then((hasil) => {
			if (dibatalkan) return;
			jalur = hasil.ok ? hasil.data.jalur : [];
			status = hasil.ok ? 'siap' : 'gagal';
		});
		return () => {
			dibatalkan = true;
		};
	});
</script>

<section
	class="mt-6 rounded-lg border border-line bg-surface p-5 md:p-6"
	aria-labelledby="judul-penyampaian"
>
	<div class="flex items-start gap-3">
		<CircleCheck size={24} class="mt-1 shrink-0 text-success" aria-hidden="true" />
		<div>
			<h2 id="judul-penyampaian" class="t-h3">{t('packet.ready')}</h2>
			<p class="mt-1 text-small text-muted">{t('packet.notSent')}</p>
		</div>
	</div>
	<ol class="mt-5 flex flex-col gap-3 text-small">
		<li class="flex gap-3">
			<span class="font-semibold" aria-hidden="true">1.</span>{t('packet.next.review')}
		</li>
		<li class="flex gap-3">
			<span class="font-semibold" aria-hidden="true">2.</span>{t('packet.next.save')}
		</li>
		<li class="flex gap-3">
			<span class="font-semibold" aria-hidden="true">3.</span>{t('packet.next.deliver')}
		</li>
	</ol>
	{#if penerima}
		<p class="mt-5 flex items-start gap-2 border-t border-line pt-4 text-small">
			<Landmark size={18} class="mt-0.5 shrink-0" aria-hidden="true" /><span
				>{t('packet.recipients')}: <strong>{penerima}</strong></span
			>
		</p>
	{/if}
	{#if status === 'memuat'}
		<p class="mt-4 text-small text-muted" aria-live="polite">{t('packet.channelsLoading')}</p>
	{:else if status === 'gagal'}
		<p class="mt-4 text-small text-muted">{t('packet.channelsFailed')}</p>
	{:else if jalur.length}
		<h3 class="mt-5 text-small font-semibold">{t('packet.channels')}</h3>
		<ul class="mt-3 divide-y divide-line">
			{#each jalur as item (item.kode)}
				<li class="py-3">
					<p class="text-small font-semibold">{item.nama[bahasaAktif]}</p>
					<p class="mt-1 text-small text-muted">{item.cara[bahasaAktif]}</p>
					{#if item.url}
						<a
							href={item.url}
							target="_blank"
							rel="noopener"
							class="inline-flex min-h-touch items-center gap-2 text-small font-semibold underline underline-offset-4"
							aria-label="{t('packet.channelInfo')}: {item.nama[bahasaAktif]} {t(
								'a11y.opensNewTab'
							)}">{t('packet.channelInfo')}<ExternalLink size={14} aria-hidden="true" /></a
						>
					{/if}
				</li>
			{/each}
		</ul>
	{/if}
</section>
