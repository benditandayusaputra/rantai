import type { Config } from 'tailwindcss';

const v = (name: string) => `var(--${name})`;

const sectors = ['edu', 'air', 'san', 'jalan', 'kes', 'ruang', 'lain'] as const;
const sectorColors = Object.fromEntries(
	sectors.map((s) => [s, { DEFAULT: v(`sector-${s}`), on: v(`sector-${s}-on`) }])
);

const statuses = ['sent', 'wait', 'ver', 'fwd', 'resp', 'done', 'rej'] as const;
const statusColors = Object.fromEntries(
	statuses.map((s) => [s, { fg: v(`st-${s}-fg`), bg: v(`st-${s}-bg`), bd: v(`st-${s}-bd`) }])
);

export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	darkMode: ['selector', '[data-theme="dark"]'],
	theme: {
		screens: {
			md: '768px',
			lg: '1024px',
			xl: '1360px'
		},
		extend: {
			colors: {
				page: v('bg'),
				surface: v('surface'),
				raised: v('raised'),
				sunken: v('sunken'),
				glass: v('glass'),
				scrim: v('scrim'),
				fg: v('text'),
				muted: v('muted'),
				line: { DEFAULT: v('line'), strong: v('line-strong') },
				focus: v('focus'),
				primary: {
					DEFAULT: v('primary'),
					hover: v('primary-hover'),
					press: v('primary-press'),
					on: v('on-primary')
				},
				secondary: { DEFAULT: v('secondary'), hover: v('secondary-hover'), on: v('on-secondary') },
				dis: { bg: v('dis-bg'), fg: v('dis-text') },
				success: { DEFAULT: v('success'), on: v('on-success'), bg: v('success-bg') },
				info: { DEFAULT: v('info'), bg: v('info-bg') },
				attn: { DEFAULT: v('attn'), fill: v('attn-fill'), bg: v('attn-bg'), on: v('on-attn-fill') },
				danger: { DEFAULT: v('danger'), bg: v('danger-bg') },
				cert: {
					resmi: { DEFAULT: v('cert-resmi'), on: v('cert-resmi-on') },
					proksi: { DEFAULT: v('cert-proksi'), bg: v('cert-proksi-bg') },
					estimasi: v('cert-estimasi'),
					sim: { DEFAULT: v('cert-sim'), bg: v('cert-sim-bg'), stripe: v('cert-sim-stripe') }
				},
				sector: sectorColors,
				st: statusColors,
				gap: {
					1: { DEFAULT: v('gap-1'), on: v('gap-1-on') },
					2: { DEFAULT: v('gap-2'), on: v('gap-2-on') },
					3: { DEFAULT: v('gap-3'), on: v('gap-3-on') },
					4: { DEFAULT: v('gap-4'), on: v('gap-4-on') }
				},
				map: { DEFAULT: v('map'), road: v('map-road'), water: v('map-water'), park: v('map-park') },
				skel: { DEFAULT: v('skel'), hi: v('skel-hi') }
			},
			fontFamily: {
				sans: [v('font-sans')],
				serif: [v('font-serif')],
				mono: [v('font-mono')]
			},
			fontSize: {
				display: ['34px', { lineHeight: '40px', letterSpacing: '-0.015em', fontWeight: '500' }],
				'display-lg': [
					'60px',
					{ lineHeight: '64px', letterSpacing: '-0.015em', fontWeight: '500' }
				],
				h1: ['28px', { lineHeight: '34px', letterSpacing: '-0.01em', fontWeight: '500' }],
				'h1-lg': ['44px', { lineHeight: '50px', letterSpacing: '-0.01em', fontWeight: '500' }],
				h2: ['22px', { lineHeight: '28px', letterSpacing: '-0.01em', fontWeight: '700' }],
				'h2-lg': ['30px', { lineHeight: '36px', letterSpacing: '-0.01em', fontWeight: '700' }],
				h3: ['18px', { lineHeight: '24px', fontWeight: '700' }],
				'h3-lg': ['22px', { lineHeight: '28px', fontWeight: '700' }],
				'body-lg': ['17px', { lineHeight: '26px' }],
				'body-lg-lg': ['18px', { lineHeight: '28px' }],
				body: ['16px', { lineHeight: '24px' }],
				small: ['14px', { lineHeight: '20px' }],
				caption: ['12px', { lineHeight: '16px', fontWeight: '500' }],
				overline: ['12px', { lineHeight: '16px', letterSpacing: '0.08em', fontWeight: '700' }],
				'num-xl': ['44px', { lineHeight: '48px', letterSpacing: '-0.03em', fontWeight: '700' }],
				'num-xl-lg': ['72px', { lineHeight: '72px', letterSpacing: '-0.03em', fontWeight: '700' }],
				'num-lg': ['30px', { lineHeight: '36px', letterSpacing: '-0.02em', fontWeight: '700' }],
				'num-lg-lg': ['44px', { lineHeight: '48px', letterSpacing: '-0.02em', fontWeight: '700' }],
				'num-md': ['20px', { lineHeight: '26px', letterSpacing: '-0.01em', fontWeight: '600' }],
				'num-md-lg': ['24px', { lineHeight: '30px', letterSpacing: '-0.01em', fontWeight: '600' }],
				mono: ['13px', { lineHeight: '20px' }]
			},
			spacing: {
				14: '56px',
				18: '72px'
			},
			maxWidth: {
				content: '1296px',
				prose: '720px'
			},
			borderRadius: {
				xs: v('r-xs'),
				sm: v('r-sm'),
				md: v('r-md'),
				lg: v('r-lg'),
				xl: v('r-xl'),
				pill: v('r-pill')
			},
			borderWidth: {
				ctrl: '1.5px',
				focus: '2px',
				link: '3px'
			},
			boxShadow: {
				1: v('sh-1'),
				2: v('sh-2'),
				3: v('sh-3'),
				ring: v('ring')
			},
			transitionDuration: {
				instant: '80ms',
				fast: '140ms',
				base: '220ms',
				enter: '320ms'
			},
			transitionTimingFunction: {
				standard: v('ease-standard'),
				enter: v('ease-enter'),
				exit: v('ease-exit')
			},
			zIndex: {
				raised: '10',
				sticky: '100',
				'map-ui': '200',
				popover: '300',
				sheet: '400',
				dialog: '500',
				toast: '600',
				tooltip: '700',
				ribbon: '800'
			},
			minHeight: { touch: '44px' },
			minWidth: { touch: '44px' },
			keyframes: {
				shimmer: {
					'0%': { backgroundPosition: '-320px 0' },
					'100%': { backgroundPosition: '320px 0' }
				},
				blink: { '50%': { opacity: '0' } }
			},
			animation: {
				shimmer: 'shimmer 1.4s linear infinite',
				caret: 'blink 1s steps(1) infinite'
			}
		}
	},
	plugins: []
} satisfies Config;
