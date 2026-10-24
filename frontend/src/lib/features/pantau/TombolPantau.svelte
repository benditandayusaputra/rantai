<script lang="ts">
	import { Bell, BellRing } from '@lucide/svelte';
	import { tangkap } from '$lib/api/klien';
	import { daftarLangganan, langgani } from '$lib/api/pantau';
	import { sesiAktif } from '$lib/features/akun/sesi.svelte';
	import { bahasaBawaan, t } from '$lib/i18n';

	let { kodeKecamatan }: { kodeKecamatan: string } = $props();

	let isDipantau = $state(false);
	let isSibuk = $state(false);
	const sesi = $derived(sesiAktif());
	const kelas =
		'inline-flex min-h-touch items-center gap-2 rounded-md border-ctrl px-4 font-semibold hover:bg-sunken';

	$effect(() => {
		if (!sesi) return;
		void tangkap(daftarLangganan(fetch, sesi.token_akses)).then((hasil) => {
			isDipantau = hasil.ok && hasil.data.some((l) => l.kode_wilayah === kodeKecamatan);
		});
	});

	async function pantau() {
		if (!sesi) return;
		isSibuk = true;
		const hasil = await tangkap(langgani(fetch, sesi.token_akses, kodeKecamatan));
		isSibuk = false;
		isDipantau = hasil.ok;
	}
</script>

{#if !sesi}
	<a
		href="/{bahasaBawaan}/masuk?lanjut={encodeURIComponent(`/${bahasaBawaan}/pantau`)}"
		class="{kelas} border-line-strong"
	>
		<Bell size={20} strokeWidth={1.75} aria-hidden="true" />{t('action.watch')}
	</a>
{:else if isDipantau}
	<a href="/{bahasaBawaan}/pantau" class="{kelas} border-primary bg-secondary text-secondary-on">
		<BellRing size={20} strokeWidth={1.75} aria-hidden="true" />{t('action.watching')}
	</a>
{:else}
	<button
		type="button"
		class="{kelas} border-line-strong hover:border-fg"
		disabled={isSibuk}
		aria-busy={isSibuk}
		onclick={pantau}
	>
		<Bell size={20} strokeWidth={1.75} aria-hidden="true" />{t('action.watch')}
	</button>
{/if}
