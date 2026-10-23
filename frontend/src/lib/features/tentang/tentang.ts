import { isKunciPesan, t } from '$lib/i18n';

export function namaJob(job: string): string {
	const kunci = `job.${job}`;
	return isKunciPesan(kunci) ? t(kunci) : job;
}

export function pemilikJob(job: string): string {
	const kunci = `job.${job}.pemilik`;
	return isKunciPesan(kunci) ? t(kunci) : '';
}
