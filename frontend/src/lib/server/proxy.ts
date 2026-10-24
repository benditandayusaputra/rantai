const headerDiteruskan = ['accept', 'authorization', 'content-type', 'cookie'];
const headerRespons = [
	'cache-control',
	'content-disposition',
	'content-type',
	'retry-after',
	'set-cookie'
];

export interface TujuanProxy {
	hulu: string;
	rahasia: string;
}

export function susunPermintaanHulu(
	permintaan: Request,
	path: string,
	tujuan: TujuanProxy,
	ipKlien: string
): Request {
	const asal = new URL(permintaan.url);
	const url = `${tujuan.hulu}/api/v1/${path}${asal.search}`;
	const header = new Headers();
	for (const nama of headerDiteruskan) {
		const nilai = permintaan.headers.get(nama);
		if (nilai) header.set(nama, nilai);
	}
	header.set('x-klien-asli', ipKlien);
	header.set('x-rantai-proxy', tujuan.rahasia);
	const isBerbadan = permintaan.method !== 'GET' && permintaan.method !== 'HEAD';
	return new Request(url, {
		method: permintaan.method,
		headers: header,
		body: isBerbadan ? permintaan.body : undefined,
		duplex: 'half'
	} as RequestInit);
}

export function saringRespons(respons: Response): Response {
	const header = new Headers();
	for (const nama of headerRespons) {
		if (nama === 'set-cookie') {
			for (const cookie of respons.headers.getSetCookie()) header.append(nama, cookie);
			continue;
		}
		const nilai = respons.headers.get(nama);
		if (nilai) header.set(nama, nilai);
	}
	return new Response(respons.body, { status: respons.status, headers: header });
}
