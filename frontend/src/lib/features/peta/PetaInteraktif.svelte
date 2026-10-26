<script lang="ts">
	import type { FeatureCollection } from 'geojson';
	import type {
		GeoJSONSource,
		LngLatBoundsLike,
		MapLayerMouseEvent,
		Map as PetaMaplibre,
		Popup as PopupMaplibre
	} from 'maplibre-gl';
	import { mount, unmount, untrack } from 'svelte';
	import { env } from '$env/dynamic/public';
	import { ambilPetaPaket, ambilSekolahTerdekat, tangkap } from '$lib/api/klien';
	import { ambilPetaLaporan } from '$lib/api/laporan';
	import type { FiturLaporan, PetaPaket, PetaTema, SekolahTerdekat, Titik } from '$lib/api/tipe';
	import { bahasaAktif, t } from '$lib/i18n';
	import {
		bacaTokenDokumen,
		buatGayaDasar,
		daftarPetaDasar,
		idLapisanRaster,
		isPetaDasarRaster,
		urlSumberDasar,
		warnaKelas,
		type PetaDasar
	} from './gaya';
	import { gabungSekolah, radiusTampilanM, type StatusLapisan } from './lapisan';
	import {
		ekspresiPenandaLaporan,
		ekspresiPenandaProyek,
		pasangPenanda,
		penandaSekolah,
		ukuranPenandaMenurutZoom
	} from './penanda';
	import { idPenanda, isIdPenanda, isiPopup } from './popup';
	import PopupPenanda from './PopupPenanda.svelte';

	interface Props {
		tema: PetaTema | null;
		temaSekolah: string;
		lapisan: StatusLapisan;
		petaDasar: PetaDasar;
		pusat: Titik;
		zoom: number;
		batas?: LngLatBoundsLike;
		isKooperatif?: boolean;
		terpilih: string | null;
		onpilih: (kode: string) => void;
		paket?: PetaPaket['features'];
		sekolah?: SekolahTerdekat[];
		laporan?: FiturLaporan[];
		status?: 'memuat' | 'siap' | 'gagal';
	}

	let {
		tema,
		temaSekolah,
		lapisan,
		petaDasar,
		pusat,
		zoom,
		batas,
		isKooperatif = false,
		terpilih,
		onpilih,
		paket = $bindable([]),
		sekolah = $bindable([]),
		laporan = $bindable([]),
		status = $bindable('memuat')
	}: Props = $props();

	const koleksiKosong: FeatureCollection = { type: 'FeatureCollection', features: [] };
	const luasBboxMaksimalDerajat = 1;
	const opasitasTema = 0.85;
	const opasitasTemaDiAtasRaster = 0.55;
	const jarakPopupPx = 28;
	const paddingBatasPx = 12;
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
		m.addLayer({
			id: 'tema-isi',
			type: 'fill',
			source: 'tema',
			paint: { 'fill-opacity': opasitasTema }
		});
		m.addLayer({ id: 'tema-garis', type: 'line', source: 'tema', paint: { 'line-width': 2 } });
		m.addLayer({ id: 'tema-pilih', type: 'line', source: 'tema', paint: { 'line-width': 3.5 } });
		m.addLayer({
			id: 'proyek-area',
			type: 'fill',
			source: 'proyek',
			filter: ['==', ['get', 'jenis_lokasi'], 'area'],
			paint: { 'fill-opacity': 0.35 }
		});
		tambahLapisanPenanda(m);
	}

	function tambahLapisanPenanda(m: PetaMaplibre) {
		const tataLetak = { 'icon-size': ukuranPenandaMenurutZoom, 'icon-allow-overlap': true };
		const tataLetakPin = { ...tataLetak, 'icon-anchor': 'bottom' as const };
		m.addLayer({
			id: 'sekolah-titik',
			type: 'symbol',
			source: 'sekolah',
			layout: { ...tataLetak, 'icon-image': penandaSekolah }
		});
		m.addLayer({
			id: 'proyek-titik',
			type: 'symbol',
			source: 'proyek',
			filter: ['==', ['get', 'jenis_lokasi'], 'titik'],
			layout: { ...tataLetakPin, 'icon-image': ekspresiPenandaProyek }
		});
		m.addLayer({
			id: 'laporan-titik',
			type: 'symbol',
			source: 'laporan',
			layout: { ...tataLetakPin, 'icon-image': ekspresiPenandaLaporan }
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
			ambilSekolahTerdekat(fetch, { lat: tengah.lat, lng: tengah.lng }, radius, temaSekolah)
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
								properties: { npsn: s.npsn, nama: s.nama, bentuk: s.bentuk, status: s.status }
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

	function tampilkanPopup(Popup: typeof PopupMaplibre, m: PetaMaplibre, e: MapLayerMouseEvent) {
		const fitur = e.features?.[0];
		if (!fitur || !isIdPenanda(fitur.layer.id)) return;
		const isi = document.createElement('div');
		const komponen = mount(PopupPenanda, {
			target: isi,
			props: isiPopup(fitur.layer.id, fitur.properties, bahasaAktif)
		});
		const titik = fitur.geometry.type === 'Point' ? fitur.geometry.coordinates : e.lngLat.toArray();
		new Popup({ offset: jarakPopupPx, className: 'popup-rantai' })
			.setLngLat(titik as [number, number])
			.setDOMContent(isi)
			.addTo(m)
			.on('close', () => unmount(komponen));
	}

	function isKenaPenanda(m: PetaMaplibre, e: MapLayerMouseEvent): boolean {
		const aktif = idPenanda.filter((id) => m.getLayer(id));
		return m.queryRenderedFeatures(e.point, { layers: aktif }).length > 0;
	}

	function pasangInteraksi(m: PetaMaplibre, Popup: typeof PopupMaplibre) {
		m.on('click', 'tema-isi', (e) => {
			const kode = e.features?.[0]?.properties?.kode;
			if (typeof kode === 'string' && !isKenaPenanda(m, e)) onpilih(kode);
		});
		for (const id of ['tema-isi', ...idPenanda]) {
			m.on('mouseenter', id, () => (m.getCanvas().style.cursor = 'pointer'));
			m.on('mouseleave', id, () => (m.getCanvas().style.cursor = ''));
		}
		for (const id of idPenanda) m.on('click', id, (e) => tampilkanPopup(Popup, m, e));
	}

	function teksGesturKooperatif() {
		return {
			'CooperativeGesturesHandler.WindowsHelpText': t('map.cooperative.windows'),
			'CooperativeGesturesHandler.MacHelpText': t('map.cooperative.mac'),
			'CooperativeGesturesHandler.MobileHelpText': t('map.cooperative.mobile')
		};
	}

	async function buatPeta(): Promise<PetaMaplibre> {
		const [maplibre, { Protocol }] = await Promise.all([
			import('maplibre-gl'),
			import('pmtiles'),
			import('maplibre-gl/dist/maplibre-gl.css')
		]);
		maplibre.addProtocol('pmtiles', new Protocol().tile);
		const awal = untrack(() => ({ pusat, zoom, batas, isKooperatif }));
		const m = new maplibre.Map({
			container: wadah,
			style: buatGayaDasar(urlSumberDasar(env.PUBLIC_PETA_DASAR_URL), bacaTokenDokumen),
			center: [awal.pusat.lng, awal.pusat.lat],
			zoom: awal.zoom,
			...(awal.batas ? { bounds: awal.batas, fitBoundsOptions: { padding: paddingBatasPx } } : {}),
			cooperativeGestures: awal.isKooperatif,
			locale: teksGesturKooperatif(),
			attributionControl: { compact: true }
		});
		m.addControl(new maplibre.NavigationControl({ showCompass: false }), 'bottom-right');
		pasangInteraksi(m, maplibre.Popup);
		m.on('moveend', () => muatDataTampilan(m));
		return m;
	}

	async function siapkan(m: PetaMaplibre) {
		await pasangPenanda(m, bacaTokenDokumen);
		tambahLapisan(m);
		warnai(m);
		peta = m;
		status = 'siap';
		muatDataTampilan(m);
	}

	function warnaiUlang(m: PetaMaplibre) {
		warnai(m);
		void pasangPenanda(m, bacaTokenDokumen);
	}

	$effect(() => {
		let instansi: PetaMaplibre | null = null;
		const pengamatTema = new MutationObserver(
			() => instansi?.isStyleLoaded() && warnaiUlang(instansi)
		);
		buatPeta()
			.then((m) => {
				instansi = m;
				m.on('load', () => siapkan(m).catch(() => (status = 'gagal')));
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
		if (!peta) return;
		for (const dasar of daftarPetaDasar.filter(isPetaDasarRaster)) {
			peta.setLayoutProperty(
				idLapisanRaster(dasar),
				'visibility',
				dasar === petaDasar ? 'visible' : 'none'
			);
		}
		const opasitas = isPetaDasarRaster(petaDasar) ? opasitasTemaDiAtasRaster : opasitasTema;
		peta.setPaintProperty('tema-isi', 'fill-opacity', opasitas);
	});

	$effect(() => {
		if (!peta || !temaSekolah) return;
		untrack(() => peta && lapisan.sekolah && muatSekolah(peta));
	});

	$effect(() => {
		peta?.setFilter('tema-pilih', ['==', ['get', 'kode'], terpilih ?? '']);
	});

	$effect(() => {
		if (!peta) return;
		if (isKooperatif) peta.cooperativeGestures.enable();
		else peta.cooperativeGestures.disable();
	});

	export function sesuaikanUkuran() {
		peta?.resize();
	}
</script>

<div bind:this={wadah} role="region" aria-label={t('map.region')} class="h-full w-full"></div>
