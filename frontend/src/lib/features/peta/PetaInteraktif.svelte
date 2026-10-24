<script lang="ts">
	import type { FeatureCollection } from 'geojson';
	import type { GeoJSONSource, Map as PetaMaplibre } from 'maplibre-gl';
	import { untrack } from 'svelte';
	import { env } from '$env/dynamic/public';
	import { ambilPetaPaket, ambilSekolahTerdekat, tangkap } from '$lib/api/klien';
	import { ambilPetaLaporan } from '$lib/api/laporan';
	import type { FiturLaporan, PetaPaket, PetaTema, SekolahTerdekat, Titik } from '$lib/api/tipe';
	import { t } from '$lib/i18n';
	import {
		bacaTokenDokumen,
		buatGayaDasar,
		sumberDasarBawaan,
		warnaKelas,
		warnaLaporan,
		warnaSektor
	} from './gaya';
	import { gabungSekolah, radiusTampilanM, type StatusLapisan } from './lapisan';

	interface Props {
		tema: PetaTema | null;
		lapisan: StatusLapisan;
		pusat: Titik;
		zoom: number;
		terpilih: string | null;
		onpilih: (kode: string) => void;
		paket: PetaPaket['features'];
		sekolah: SekolahTerdekat[];
		laporan: FiturLaporan[];
		status: 'memuat' | 'siap' | 'gagal';
	}

	let {
		tema,
		lapisan,
		pusat,
		zoom,
		terpilih,
		onpilih,
		paket = $bindable(),
		sekolah = $bindable(),
		laporan = $bindable(),
		status = $bindable()
	}: Props = $props();

	const koleksiKosong: FeatureCollection = { type: 'FeatureCollection', features: [] };
	const luasBboxMaksimalDerajat = 1;
	const idLapisan = {
		tema: ['tema-isi', 'tema-garis', 'tema-pilih'],
		proyek: ['proyek-area', 'proyek-titik'],
		sekolah: ['sekolah-titik'],
		laporan: ['laporan-titik']
	} as const;

	let wadah: HTMLDivElement;
	let peta = $state.raw<PetaMaplibre | null>(null);

	function sumber(m: PetaMaplibre, nama: string): GeoJSONSource | undefined {
		return m.getSource<GeoJSONSource>(nama);
	}

	function tambahLapisan(m: PetaMaplibre) {
		for (const nama of ['tema', 'proyek', 'sekolah', 'laporan']) {
			m.addSource(nama, { type: 'geojson', data: koleksiKosong });
		}
		m.addLayer({ id: 'tema-isi', type: 'fill', source: 'tema', paint: { 'fill-opacity': 0.85 } });
		m.addLayer({ id: 'tema-garis', type: 'line', source: 'tema', paint: { 'line-width': 2 } });
		m.addLayer({ id: 'tema-pilih', type: 'line', source: 'tema', paint: { 'line-width': 3.5 } });
		m.addLayer({
			id: 'proyek-area',
			type: 'fill',
			source: 'proyek',
			filter: ['==', ['get', 'jenis_lokasi'], 'area'],
			paint: { 'fill-opacity': 0.35 }
		});
		m.addLayer({
			id: 'proyek-titik',
			type: 'circle',
			source: 'proyek',
			filter: ['==', ['get', 'jenis_lokasi'], 'titik'],
			paint: { 'circle-radius': 7, 'circle-stroke-width': 2.5 }
		});
		m.addLayer({
			id: 'sekolah-titik',
			type: 'circle',
			source: 'sekolah',
			paint: { 'circle-radius': 6, 'circle-stroke-width': 2 }
		});
		m.addLayer({
			id: 'laporan-titik',
			type: 'circle',
			source: 'laporan',
			paint: { 'circle-radius': 6, 'circle-stroke-width': 2 }
		});
	}

	function warnai(m: PetaMaplibre) {
		const token = bacaTokenDokumen;
		m.setPaintProperty('latar', 'background-color', token('map'));
		m.setPaintProperty('taman', 'fill-color', token('map-park'));
		m.setPaintProperty('air', 'fill-color', token('map-water'));
		m.setPaintProperty('jalan', 'line-color', token('map-road'));
		m.setPaintProperty('tema-isi', 'fill-color', warnaKelas(token));
		m.setPaintProperty('tema-garis', 'line-color', token('map'));
		m.setPaintProperty('tema-pilih', 'line-color', token('text'));
		m.setPaintProperty('proyek-area', 'fill-color', token('sector-lain'));
		m.setPaintProperty('proyek-titik', 'circle-color', warnaSektor(token));
		m.setPaintProperty('proyek-titik', 'circle-stroke-color', token('surface'));
		m.setPaintProperty('sekolah-titik', 'circle-color', token('surface'));
		m.setPaintProperty('sekolah-titik', 'circle-stroke-color', token('text'));
		m.setPaintProperty('laporan-titik', 'circle-color', warnaLaporan(token));
		m.setPaintProperty('laporan-titik', 'circle-stroke-color', token('surface'));
	}

	function bboxTampilan(m: PetaMaplibre): string | null {
		const b = m.getBounds();
		const luas = (b.getEast() - b.getWest()) * (b.getNorth() - b.getSouth());
		if (luas > luasBboxMaksimalDerajat) return null;
		return [b.getWest(), b.getSouth(), b.getEast(), b.getNorth()]
			.map((n) => n.toFixed(5))
			.join(',');
	}

	async function muatProyek(m: PetaMaplibre) {
		const bbox = bboxTampilan(m);
		const hasil = bbox ? await tangkap(ambilPetaPaket(fetch, bbox)) : null;
		paket = hasil?.ok ? hasil.data.features : [];
		sumber(m, 'proyek')?.setData({ type: 'FeatureCollection', features: paket } as never);
	}

	async function muatSekolah(m: PetaMaplibre) {
		const tengah = m.getCenter();
		const b = m.getBounds();
		const radius = radiusTampilanM(b.getNorth(), b.getSouth());
		const hasil = await tangkap(
			ambilSekolahTerdekat(fetch, { lat: tengah.lat, lng: tengah.lng }, radius)
		);
		sekolah = hasil.ok ? gabungSekolah(hasil.data) : [];
		sumber(m, 'sekolah')?.setData({
			type: 'FeatureCollection',
			features: sekolah.flatMap((s) =>
				s.lokasi
					? [
							{
								type: 'Feature',
								geometry: { type: 'Point', coordinates: [s.lokasi.lng, s.lokasi.lat] },
								properties: { npsn: s.npsn, nama: s.nama }
							}
						]
					: []
			)
		} as never);
	}

	async function muatLaporan(m: PetaMaplibre) {
		const bbox = bboxTampilan(m);
		const hasil = bbox ? await tangkap(ambilPetaLaporan(fetch, bbox)) : null;
		laporan = hasil?.ok ? hasil.data.features : [];
		sumber(m, 'laporan')?.setData({ type: 'FeatureCollection', features: laporan } as never);
	}

	function muatDataTampilan(m: PetaMaplibre) {
		if (lapisan.proyek) muatProyek(m);
		if (lapisan.sekolah) muatSekolah(m);
		if (lapisan.laporan) muatLaporan(m);
	}

	async function buatPeta(): Promise<PetaMaplibre> {
		const maplibre = await import('maplibre-gl');
		await import('maplibre-gl/dist/maplibre-gl.css');
		const awal = untrack(() => ({ pusat, zoom }));
		const m = new maplibre.Map({
			container: wadah,
			style: buatGayaDasar(env.PUBLIC_PETA_DASAR_URL || sumberDasarBawaan, bacaTokenDokumen),
			center: [awal.pusat.lng, awal.pusat.lat],
			zoom: awal.zoom,
			attributionControl: { compact: true }
		});
		m.addControl(new maplibre.NavigationControl({ showCompass: false }), 'bottom-right');
		m.on('click', 'tema-isi', (e) => {
			const kode = e.features?.[0]?.properties?.kode;
			if (typeof kode === 'string') onpilih(kode);
		});
		m.on('mouseenter', 'tema-isi', () => (m.getCanvas().style.cursor = 'pointer'));
		m.on('mouseleave', 'tema-isi', () => (m.getCanvas().style.cursor = ''));
		m.on('moveend', () => muatDataTampilan(m));
		return m;
	}

	$effect(() => {
		let instansi: PetaMaplibre | null = null;
		const pengamatTema = new MutationObserver(() => instansi?.isStyleLoaded() && warnai(instansi));
		buatPeta()
			.then((m) => {
				instansi = m;
				m.on('load', () => {
					tambahLapisan(m);
					warnai(m);
					peta = m;
					status = 'siap';
					muatDataTampilan(m);
				});
			})
			.catch(() => (status = 'gagal'));
		pengamatTema.observe(document.documentElement, { attributeFilter: ['data-theme'] });
		return () => {
			pengamatTema.disconnect();
			instansi?.remove();
			peta = null;
		};
	});

	$effect(() => {
		if (peta) sumber(peta, 'tema')?.setData((tema ?? koleksiKosong) as never);
	});

	$effect(() => {
		if (!peta) return;
		for (const nama of ['tema', 'proyek', 'sekolah', 'laporan'] as const) {
			for (const id of idLapisan[nama]) {
				peta.setLayoutProperty(id, 'visibility', lapisan[nama] ? 'visible' : 'none');
			}
		}
		untrack(() => peta && muatDataTampilan(peta));
	});

	$effect(() => {
		peta?.setFilter('tema-pilih', ['==', ['get', 'kode'], terpilih ?? '']);
	});

	export function sesuaikanUkuran() {
		peta?.resize();
	}
</script>

<div bind:this={wadah} role="region" aria-label={t('map.region')} class="h-full w-full"></div>
