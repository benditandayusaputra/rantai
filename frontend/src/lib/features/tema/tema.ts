export const kunciTema = 'rantai-theme';

export type Tema = 'light' | 'dark';

function pasangTema(tema: Tema) {
	document.documentElement.setAttribute('data-theme', tema);
}

export function temaTerpasang(): Tema {
	return document.documentElement.getAttribute('data-theme') === 'dark' ? 'dark' : 'light';
}

export function temaLawan(tema: Tema): Tema {
	return tema === 'dark' ? 'light' : 'dark';
}

export function pilihTema(tema: Tema) {
	pasangTema(tema);
	try {
		localStorage.setItem(kunciTema, tema);
	} catch {
		return;
	}
}
