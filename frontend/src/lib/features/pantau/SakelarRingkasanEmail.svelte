<script lang="ts">
	import { Mail } from '@lucide/svelte';
	import { onMount } from 'svelte';
	import {
		aktifkanRingkasanEmail,
		ambilRingkasanEmail,
		matikanRingkasanEmail
	} from '$lib/api/pantau';
	import { tangkap } from '$lib/api/klien';
	import { t } from '$lib/i18n';

	let { token }: { token: string } = $props();

	let isAktif = $state(false);
	let isSiap = $state(false);
	let isSibuk = $state(false);
	let isGagal = $state(false);

	onMount(async () => {
		const hasil = await tangkap(ambilRingkasanEmail(fetch, token));
		if (!hasil.ok) return;
		isAktif = hasil.data.aktif;
		isSiap = true;
	});

	async function ubah() {
		isSibuk = true;
		const hasil = await tangkap(
			isAktif ? matikanRingkasanEmail(fetch, token) : aktifkanRingkasanEmail(fetch, token)
		);
		isSibuk = false;
		isGagal = !hasil.ok;
		if (hasil.ok) isAktif = hasil.data.aktif;
	}
</script>

{#if isSiap}
	<section class="card flex flex-col gap-2 p-5">
		<h2 class="flex items-center gap-2 font-bold">
			<Mail size={18} strokeWidth={1.75} aria-hidden="true" class="text-muted" />
			{t('watch.email.title')}
		</h2>
		<label class="flex min-h-touch cursor-pointer items-center justify-between gap-3">
			<span>{t('watch.email.label')}</span>
			<input
				type="checkbox"
				role="switch"
				class="h-6 w-6 accent-primary"
				checked={isAktif}
				disabled={isSibuk}
				aria-describedby="bantuan-ringkasan-email"
				onchange={ubah}
			/>
		</label>
		<p
			id="bantuan-ringkasan-email"
			class="text-small {isGagal ? 'text-danger' : 'text-muted'}"
			aria-live="polite"
		>
			{isGagal ? t('watch.email.failed') : t('watch.email.help')}
		</p>
	</section>
{/if}
