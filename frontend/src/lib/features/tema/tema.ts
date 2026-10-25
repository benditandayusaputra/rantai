export const kunciTema = 'rantai-theme';
const kueriGelap = '(prefers-color-scheme: dark)';

export type Tema = 'light' | 'dark';

function isMengikutiSistem(): boolean {
	try {
		return (localStorage.getItem(kunciTema) ?? 'system') === 'system';
	} catch {
		return true;
	}
}

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

export function ikutiTemaSistem(): () => void {
	const media = matchMedia(kueriGelap);
	const tangani = (e: MediaQueryListEvent) => {
		if (isMengikutiSistem()) pasangTema(e.matches ? 'dark' : 'light');
	};
	media.addEventListener('change', tangani);
	return () => media.removeEventListener('change', tangani);
}
