<script lang="ts">
	import type { SaksiLaporan, StatusLaporan } from '$lib/api/tipe';
	import { tampilanStatus } from '$lib/features/lapor/status';
	import { t } from '$lib/i18n';

	interface Props {
		status: StatusLaporan;
		witnesses?: SaksiLaporan;
	}

	let { status, witnesses }: Props = $props();

	const tampilan = $derived(tampilanStatus[status]);
	const Ikon = $derived(tampilan.ikon);
	const label = $derived(
		witnesses
			? t(tampilan.label, { n: witnesses.jumlah, m: witnesses.dibutuhkan })
			: t(status === 'MENUNGGU_SAKSI' ? 'report.menunggu-saksi.singkat' : tampilan.label)
	);
</script>

<span
	class="inline-flex min-h-8 items-center gap-1.5 rounded-xs border-ctrl px-2.5 text-small font-bold {tampilan.kelas}"
>
	<Ikon size={16} strokeWidth={1.75} aria-hidden="true" />{label}
</span>
