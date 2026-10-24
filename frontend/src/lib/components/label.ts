import type { Label } from '$lib/api/tipe';

export const kelasLabel: Record<Label, string> = {
	resmi: 'border-cert-resmi bg-cert-resmi text-cert-resmi-on',
	proksi: 'border-cert-proksi bg-cert-proksi-bg text-cert-proksi',
	estimasi: 'border-dashed border-cert-estimasi text-cert-estimasi',
	simulasi: 'hatch-sim border-cert-sim text-cert-sim'
};
