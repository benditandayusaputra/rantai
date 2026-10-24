<script lang="ts">
	import PhotoGrid from '$lib/components/PhotoGrid.svelte';
	import { t } from '$lib/i18n';
	import { draf } from './draf.svelte';
	import { siapkanFoto } from './foto';
	import { batasFoto } from './lapor';

	let isMenyiapkan = $state(false);
	let isGagal = $state(false);

	async function tambah(semua: File[]) {
		isMenyiapkan = true;
		isGagal = false;
		for (const berkas of semua) {
			try {
				const hasil = await siapkanFoto(berkas);
				draf.foto.push({
					id: crypto.randomUUID(),
					berkas: hasil,
					pratinjau: URL.createObjectURL(hasil),
					diambilPada: new Date(berkas.lastModified)
				});
			} catch {
				isGagal = true;
			}
		}
		isMenyiapkan = false;
	}

	function hapus(id: string) {
		const foto = draf.foto.find((f) => f.id === id);
		if (foto) URL.revokeObjectURL(foto.pratinjau);
		draf.foto = draf.foto.filter((f) => f.id !== id);
	}
</script>

<div class="flex flex-col gap-4">
	<h1 class="t-h1">{t('report.s3.title')}</h1>
	<PhotoGrid
		id="foto-bukti"
		photos={draf.foto.map((f) => ({ id: f.id, src: f.pratinjau, takenAt: f.diambilPada }))}
		max={batasFoto}
		busy={isMenyiapkan}
		onadd={tambah}
		onremove={hapus}
	/>
	{#if isGagal}<p role="alert" class="text-danger">{t('report.s3.failed')}</p>{/if}
	<p class="text-small text-muted">{t('report.s3.privacy')}</p>
</div>
