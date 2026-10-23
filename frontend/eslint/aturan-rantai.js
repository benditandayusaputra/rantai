const polaNilaiMentah = /#[0-9a-fA-F]{3,8}\b|\b(?:rgba?|hsla?)\(|-\[[^\]]*\d(?:px|rem|em)\]/g;
const polaHuruf = /\p{L}/u;
const atributTeks = new Set(['aria-label', 'title', 'alt', 'placeholder', 'label']);

const tanpaKomentar = {
	meta: { type: 'problem', messages: { komentar: 'Kode tidak boleh berisi komentar.' } },
	create(context) {
		const laporkan = (node) => context.report({ loc: node.loc, messageId: 'komentar' });
		return {
			Program() {
				context.sourceCode.getAllComments().forEach(laporkan);
			},
			SvelteHTMLComment: laporkan
		};
	}
};

const tanpaNilaiMentah = {
	meta: {
		type: 'problem',
		messages: { mentah: 'Pakai token desain, bukan nilai warna atau ukuran mentah: {{nilai}}.' }
	},
	create(context) {
		return {
			Program() {
				const teks = context.sourceCode.text;
				for (const cocok of teks.matchAll(polaNilaiMentah)) {
					context.report({
						loc: context.sourceCode.getLocFromIndex(cocok.index),
						messageId: 'mentah',
						data: { nilai: cocok[0] }
					});
				}
			}
		};
	}
};

const tanpaTeksMentah = {
	meta: {
		type: 'problem',
		messages: { teks: 'Teks antarmuka wajib lewat i18n: "{{teks}}".' }
	},
	create(context) {
		const laporkan = (node, teks) =>
			context.report({ loc: node.loc, messageId: 'teks', data: { teks: teks.trim() } });
		return {
			SvelteText(node) {
				if (polaHuruf.test(node.value)) laporkan(node, node.value);
			},
			SvelteAttribute(node) {
				if (!atributTeks.has(node.key.name)) return;
				for (const bagian of node.value) {
					if (bagian.type === 'SvelteLiteral' && polaHuruf.test(bagian.value)) {
						laporkan(bagian, bagian.value);
					}
				}
			}
		};
	}
};

export default {
	rules: {
		'tanpa-komentar': tanpaKomentar,
		'tanpa-nilai-mentah': tanpaNilaiMentah,
		'tanpa-teks-mentah': tanpaTeksMentah
	}
};
