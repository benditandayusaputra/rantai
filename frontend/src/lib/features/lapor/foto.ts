import { ukuranSkala } from './lapor';

const sisiTerpanjangPiksel = 2048;
const kualitasJPEG = 0.85;

export class GalatFoto extends Error {
	constructor(message: string) {
		super(message);
		this.name = 'GalatFoto';
	}
}

export async function siapkanFoto(berkas: Blob): Promise<Blob> {
	const gambar = await createImageBitmap(berkas, { imageOrientation: 'from-image' });
	const { lebar, tinggi } = ukuranSkala(gambar.width, gambar.height, sisiTerpanjangPiksel);
	const kanvas = document.createElement('canvas');
	kanvas.width = lebar;
	kanvas.height = tinggi;
	const konteks = kanvas.getContext('2d');
	if (!konteks) throw new GalatFoto('kanvas tidak tersedia');
	konteks.drawImage(gambar, 0, 0, lebar, tinggi);
	gambar.close();
	return new Promise((selesai, gagal) =>
		kanvas.toBlob(
			(hasil) => (hasil ? selesai(hasil) : gagal(new GalatFoto('foto tidak dapat diubah ke JPEG'))),
			'image/jpeg',
			kualitasJPEG
		)
	);
}
