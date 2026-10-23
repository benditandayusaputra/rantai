import { describe, expect, it } from 'vitest';
import { namaJob, pemilikJob } from './tentang';

describe('nama job sumber', () => {
	it('menerjemahkan job yang dikenal dan membiarkan job baru apa adanya', () => {
		expect(namaJob('dapodik')).toBe('Sekolah dan rombel');
		expect(pemilikJob('dapodik')).toBe('Dapodik, Kemendikdasmen');
		expect(namaJob('job_baru')).toBe('job_baru');
		expect(pemilikJob('job_baru')).toBe('');
	});
});
