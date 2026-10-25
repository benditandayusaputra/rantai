<script lang="ts">
	import { Bell, BellRing, LoaderCircle } from '@lucide/svelte';
	import { tangkap } from '$lib/api/klien';
	import { daftarLangganan, langgani } from '$lib/api/pantau';
	import { sesiAktif } from '$lib/features/akun/sesi.svelte';
	import { bahasaBawaan, t } from '$lib/i18n';

	let { kodeKecamatan }: { kodeKecamatan: string } = $props();

	let isDipantau = $state(false);
	let isSibuk = $state(false);
	const sesi = $derived(sesiAktif());
	const kelas = 'btn max-md:px-3 max-md:text-small';

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
		class="{kelas} btn-outline"
	>
		<Bell size={18} strokeWidth={1.75} aria-hidden="true" />{t('action.watch')}
	</a>
{:else if isDipantau}
	<a
		href="/{bahasaBawaan}/pantau"
		class="{kelas} border-primary bg-secondary text-secondary-on hover:bg-secondary-hover"
	>
		<BellRing size={18} strokeWidth={1.75} aria-hidden="true" />{t('action.watching')}
	</a>
{:else}
	<button
		type="button"
		class="{kelas} btn-outline"
		disabled={isSibuk}
		aria-busy={isSibuk}
		onclick={pantau}
	>
		{#if isSibuk}
			<LoaderCircle size={18} strokeWidth={1.75} aria-hidden="true" class="animate-spin" />
		{:else}
			<Bell size={18} strokeWidth={1.75} aria-hidden="true" />
		{/if}
		{t('action.watch')}
	</button>
{/if}
