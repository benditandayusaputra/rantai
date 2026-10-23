export const kunciTema = 'rantai-theme';
const kueriGelap = '(prefers-color-scheme: dark)';

function isMengikutiSistem(): boolean {
	try {
		return (localStorage.getItem(kunciTema) ?? 'system') === 'system';
	} catch {
		return true;
	}
}

function pasangTemaSistem(gelap: boolean) {
	document.documentElement.setAttribute('data-theme', gelap ? 'dark' : 'light');
}

export function ikutiTemaSistem(): () => void {
	const media = matchMedia(kueriGelap);
	const tangani = (e: MediaQueryListEvent) => {
		if (isMengikutiSistem()) pasangTemaSistem(e.matches);
	};
	media.addEventListener('change', tangani);
	return () => media.removeEventListener('change', tangani);
}
