<script lang="ts">
	import type { SaksiLaporan, StatusLaporan } from '$lib/api/tipe';
	import { tampilanStatus } from '$lib/features/lapor/status';
	import { t } from '$lib/i18n';

	interface Props {
		status: StatusLaporan;
		witnesses?: SaksiLaporan;
	}

	let { status, witnesses = { jumlah: 0, dibutuhkan: 3 } }: Props = $props();

	const tampilan = $derived(tampilanStatus[status]);
</script>

<span
	class="inline-flex min-h-8 items-center rounded-pill border-ctrl px-3 text-small font-bold {tampilan.kelas}"
	>{t(tampilan.label, { n: witnesses.jumlah, m: witnesses.dibutuhkan })}</span
>
