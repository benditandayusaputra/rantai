import { describe, expect, it } from 'vitest';
import type { SimpulJaringan, SisiJaringan } from '$lib/api/tipe';
import {
	cariTetangga,
	hitungJari,
	nilaiMaksPerJenis,
	kotakLabel,
	saringLabelBertumpuk,
	susunBarisRelasi
} from './jaringan';
import { susunTataLetak } from './tataLetak';

function simpul(id: string, jenis: SimpulJaringan['jenis'], nilai = 100): SimpulJaringan {
	return {
		id,
		jenis,
		kunci: id.split(':')[1] ?? id,
		label: id,
		jumlah_paket: 1,
		nilai: { nilai, satuan: 'rupiah', label: 'resmi', provenance_id: 'p' }
	};
}

const paketA = simpul('paket:A', 'paket', 200);
const paketB = simpul('paket:B', 'paket', 100);
const daftarSimpul = [
	simpul('instansi:1', 'instansi', 300),
	paketA,
	paketB,
	simpul('penyedia:9', 'penyedia', 200)
];
const daftarSisi: SisiJaringan[] = [
	{ dari: 'instansi:1', ke: 'paket:A' },
	{ dari: 'instansi:1', ke: 'paket:B' },
	{ dari: 'paket:A', ke: 'penyedia:9' }
];

describe('jaringan', () => {
	it('mencari tetangga dua arah', () => {
		expect([...cariTetangga('paket:A', daftarSisi)].sort()).toEqual(['instansi:1', 'penyedia:9']);
	});

	it('memberi jari terbesar pada nilai tertinggi per jenis', () => {
		const maks = nilaiMaksPerJenis(daftarSimpul);
		expect(hitungJari(paketA, maks)).toBeGreaterThan(hitungJari(paketB, maks));
	});

	it('menyusun baris relasi paket dengan instansi dan penyedia', () => {
		const baris = susunBarisRelasi(daftarSimpul, daftarSisi);
		expect(baris.map((b) => [b.paket.id, b.instansi?.id, b.penyedia?.id ?? null])).toEqual([
			['paket:A', 'instansi:1', 'penyedia:9'],
			['paket:B', 'instansi:1', null]
		]);
	});

	it('membuang label yang bertumpuk dengan label berprioritas lebih tinggi atau penghalang', () => {
		const tampil = saringLabelBertumpuk(
			[
				{ id: 'bentuk', x: 300, y: 100, lebar: 30, tinggi: 30 },
				{ id: 'terpilih', x: 100, y: 90, lebar: 30, tinggi: 30 }
			],
			[
				kotakLabel('terpilih', 100, 100, 'PT ALMA KARYA SEJATI'),
				kotakLabel('tumpuk', 120, 106, 'Perbaikan Pintu Air'),
				kotakLabel('di-bawah', 110, 130, 'Pembangunan Tangki'),
				kotakLabel('kena-bentuk', 300, 110, 'Turap'),
				kotakLabel('jauh', 500, 100, 'Perbaikan Turap')
			]
		);
		expect([...tampil].sort()).toEqual(['di-bawah', 'jauh', 'terpilih']);
	});

	it('menata simpul di dalam bidang dan mendekatkan simpul yang bertaut', () => {
		const bidang = { lebar: 600, tinggi: 400, tepi: 20 };
		const id = [...daftarSimpul.map((s) => s.id), 'penyedia:lepas'];
		const posisi = susunTataLetak(id, daftarSisi, bidang);
		for (const p of posisi.values()) {
			expect(p.x).toBeGreaterThanOrEqual(20);
			expect(p.x).toBeLessThanOrEqual(580);
			expect(p.y).toBeGreaterThanOrEqual(20);
			expect(p.y).toBeLessThanOrEqual(380);
		}
		const jarak = (a: string, b: string) => {
			const pa = posisi.get(a) ?? { x: Infinity, y: Infinity };
			const pb = posisi.get(b) ?? { x: -Infinity, y: -Infinity };
			return Math.hypot(pa.x - pb.x, pa.y - pb.y);
		};
		expect(jarak('paket:A', 'penyedia:9')).toBeLessThan(jarak('paket:B', 'penyedia:lepas'));
	});
});
