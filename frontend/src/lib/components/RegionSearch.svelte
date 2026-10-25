<script lang="ts">
	import { LocateFixed, RefreshCw, Search } from '@lucide/svelte';
	import type { WilayahDitemukan } from '$lib/api/tipe';
	import {
		bagiSorotan,
		induk,
		jedaKetikMs,
		panjangKataMinimal,
		saringHasilPencarian,
		type WilayahHalaman
	} from '$lib/features/wilayah/wilayah';
	import { t } from '$lib/i18n';

	interface Props {
		cari: (q: string) => Promise<WilayahDitemukan[]>;
		onselect: (wilayah: WilayahDitemukan) => void;
		onuselocation?: () => void;
		value?: string;
		size?: 'md' | 'lg';
	}

	let { cari, onselect, onuselocation, value = $bindable(''), size = 'md' }: Props = $props();

	type Status = 'idle' | 'loading' | 'ready' | 'empty' | 'error';

	const id = $props.id();
	const idDaftar = `${id}-daftar`;
	let status = $state<Status>('idle');
	let hasil = $state<WilayahHalaman[]>([]);
	let terbuka = $state(false);
	let aktif = $state(-1);
	let urutanPermintaan = 0;
	let masukan: HTMLInputElement;

	const jumlahOpsi = $derived(hasil.length + (onuselocation ? 1 : 0));
	const indeksLokasi = $derived(onuselocation ? hasil.length : -1);
	const kelompok = $derived([
		{ judul: t('search.group.kelurahan'), isi: hasil.filter((w) => w.tingkat === 'kelurahan') },
		{ judul: t('search.group.kecamatan'), isi: hasil.filter((w) => w.tingkat === 'kecamatan') }
	]);
	const pengumuman = $derived(
		status === 'ready'
			? t('search.results', { n: hasil.length })
			: status === 'empty'
				? t('search.empty', { q: value })
				: ''
	);

	async function muat(kata: string) {
		const nomor = ++urutanPermintaan;
		status = 'loading';
		try {
			const ditemukan = saringHasilPencarian(await cari(kata));
			if (nomor !== urutanPermintaan) return;
			hasil = ditemukan;
			status = ditemukan.length ? 'ready' : 'empty';
		} catch {
			if (nomor === urutanPermintaan) status = 'error';
		}
	}

	$effect(() => {
		const kata = value.trim();
		aktif = -1;
		if (kata.length < panjangKataMinimal) {
			urutanPermintaan++;
			hasil = [];
			status = 'idle';
			return;
		}
		bukaSaran();
		const jeda = setTimeout(() => muat(kata), jedaKetikMs);
		return () => clearTimeout(jeda);
	});

	function bukaSaran() {
		terbuka = true;
		if (document.activeElement !== masukan) masukan.focus();
	}

	function pilih(indeks: number) {
		terbuka = false;
		const wilayah = hasil[indeks];
		if (wilayah) onselect(wilayah);
		else if (indeks === indeksLokasi) onuselocation?.();
	}

	function geser(arah: 1 | -1) {
		terbuka = true;
		if (jumlahOpsi) aktif = (aktif + arah + jumlahOpsi) % jumlahOpsi;
	}

	function tanganiTombol(e: KeyboardEvent) {
		if (e.key === 'ArrowDown' || e.key === 'ArrowUp') {
			e.preventDefault();
			geser(e.key === 'ArrowDown' ? 1 : -1);
		} else if (e.key === 'Enter' && terbuka && aktif >= 0) {
			e.preventDefault();
			pilih(aktif);
		} else if (e.key === 'Escape') {
			terbuka = false;
		}
	}

	const idOpsi = (indeks: number) => `${idDaftar}-${indeks}`;
</script>

<div class="relative">
	<label for={id} class="sr-only">{t('search.label')}</label>
	<div
		class="flex items-center gap-3 rounded-md border-ctrl border-line-strong bg-raised px-4 transition-colors duration-fast ease-standard focus-within:border-focus focus-within:shadow-ring hover:border-fg {size ===
		'lg'
			? 'h-14 lg:h-15'
			: 'h-14'}"
	>
		<Search size={20} strokeWidth={1.75} aria-hidden="true" class="shrink-0 text-muted" />
		<input
			{id}
			bind:this={masukan}
			bind:value
			type="search"
			role="combobox"
			autocomplete="off"
			aria-autocomplete="list"
			aria-expanded={terbuka && status !== 'idle'}
			aria-controls={idDaftar}
			aria-activedescendant={aktif >= 0 ? idOpsi(aktif) : undefined}
			aria-describedby="{id}-bantuan"
			placeholder={t('search.placeholder')}
			class="h-full w-full bg-transparent text-body outline-none placeholder:text-muted focus-visible:shadow-none"
			oninput={() => (terbuka = true)}
			onfocus={() => (terbuka = true)}
			onblur={() => (terbuka = false)}
			onkeydown={tanganiTombol}
		/>
	</div>
	<p id="{id}-bantuan" class="sr-only">{t('search.hint')}</p>
	<div
		class="absolute inset-x-0 top-full z-popover mt-2 overflow-hidden rounded-lg border border-line bg-raised shadow-2"
		hidden={!terbuka || status === 'idle'}
	>
		{#if status === 'loading'}
			<p class="px-4 py-3 text-small text-muted">{t('search.loading', { q: value.trim() })}</p>
		{:else if status === 'empty'}
			<div class="px-4 py-3">
				<p class="text-small font-semibold">{t('search.empty', { q: value.trim() })}</p>
				<p class="text-small text-muted">{t('search.emptyHint')}</p>
			</div>
		{:else if status === 'error'}
			<div class="flex items-center justify-between gap-3 px-4 py-2">
				<p class="text-small text-danger">{t('search.error')}</p>
				<button
					type="button"
					class="btn btn-ghost btn-sm text-secondary-on"
					onmousedown={(e) => e.preventDefault()}
					onclick={() => muat(value.trim())}
				>
					<RefreshCw size={16} strokeWidth={1.75} aria-hidden="true" />
					{t('action.retry')}
				</button>
			</div>
		{/if}
		<ul
			id={idDaftar}
			role="listbox"
			aria-label={t('search.label')}
			class="max-h-96 overflow-y-auto"
		>
			{#if status === 'ready'}
				{#each kelompok as grup (grup.judul)}
					{#if grup.isi.length}
						<li role="presentation" class="t-over px-4 pb-1 pt-3">{grup.judul}</li>
						{#each grup.isi as wilayah (wilayah.kode)}
							{@const indeks = hasil.indexOf(wilayah)}
							{@const [sebelum, cocok, sesudah] = bagiSorotan(wilayah.nama, value)}
							<li
								id={idOpsi(indeks)}
								role="option"
								aria-selected={aktif === indeks}
								class="flex min-h-13 cursor-pointer flex-col justify-center px-4 py-2 {aktif ===
								indeks
									? 'bg-secondary'
									: 'hover:bg-sunken'}"
								onmousedown={(e) => e.preventDefault()}
								onclick={() => pilih(indeks)}
								onkeydown={tanganiTombol}
							>
								<span>{sebelum}<strong>{cocok}</strong>{sesudah}</span>
								<span class="text-small text-muted">{induk(wilayah)}</span>
							</li>
						{/each}
					{/if}
				{/each}
			{/if}
			{#if onuselocation && status !== 'loading'}
				<li
					id={idOpsi(indeksLokasi)}
					role="option"
					aria-selected={aktif === indeksLokasi}
					class="flex min-h-13 cursor-pointer items-center gap-3 border-t border-line px-4 font-semibold text-secondary-on {aktif ===
					indeksLokasi
						? 'bg-secondary'
						: 'hover:bg-sunken'}"
					onmousedown={(e) => e.preventDefault()}
					onclick={() => pilih(indeksLokasi)}
					onkeydown={tanganiTombol}
				>
					<LocateFixed size={20} strokeWidth={1.75} aria-hidden="true" />
					{t('location.useOnly')}
				</li>
			{/if}
		</ul>
	</div>
	<p class="sr-only" aria-live="polite">{pengumuman}</p>
</div>
