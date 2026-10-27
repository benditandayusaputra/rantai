# Rantai Design System

Version 1.0 · 25 October 2026

## 1. Purpose and authority

This file defines the shared visual and interaction rules for Rantai, a public transparency app for Jakarta. Use it when creating or changing any page, component, chart, map, or state.

The signed-in homepage in `rantai-beranda.html` is the visual reference. This guide extends that direction across the app and sets reusable tokens. It is a design specification; it does not claim that every production requirement is already implemented in the mockup.

Read this file with the product requirements and the supplied `antislop-ui` principles. If `antislop.md` is available in the project, read its core rules as well. An absent reference must not be treated as a completed review.

Follow this order when making decisions:

1. Preserve data accuracy, neutral wording, and accessibility.
2. Apply the shared tokens and component rules below.
3. Choose the page structure from the task the person needs to complete.
4. Add a visual exception only when its purpose is clear and recorded.

For a monorepo, keep this file beside the web app's root. Share its tokens through one stylesheet or theme package. Changes to shared rules belong here before they spread across pages.

## 2. Product character

Rantai should feel like a clear, modern public-interest news portal: bright, precise, calm, and easy to inspect. Its identity comes from strong sans-serif headings, disciplined alignment, relevant Jakarta imagery, and a small amount of red.

The visual direction combines the reading order of an editorial site, the clear controls of a Material interface, and limited glass effects on floating surfaces.

| Design choice | Shared direction |
| --- | --- |
| Theme | Light by default, suited to reading public information |
| Core palette | Paper, white, charcoal, and muted gray |
| Accent | Rantai red, used for brand, selection, and the main action |
| Typography | Plus Jakarta Sans for headings; Inter for UI and reading text |
| Composition | Clear grid, varied section sizes, strong content hierarchy |
| Surfaces | Mostly flat; borders and spacing define groups |
| Images | Relevant photography and maps with clear geographic context |
| Motion | Level 1 of 5: brief responses to input and changes of state |
| Section rhythm | Level 3 of 5: visible variation that follows content weight |

All product themes, including Pendidikan, Banjir, Kesehatan, and Jalan, inherit this identity. A new topic does not create a new app palette.

## 3. Color tokens

Use tokens instead of new hex values in individual components. Existing names such as `--paper` and `--red` are retained for a simple move from the mockup.

| Token | Value | Use |
| --- | --- | --- |
| `--paper` | `#FAFAF8` | Page background |
| `--white` | `#FFFFFF` | Solid surfaces and text on red or dark backgrounds |
| `--ink` | `#202323` | Headings, body text, important values |
| `--muted` | `#626664` | Supporting text and source notes |
| `--line` | `#E3E5E1` | Dividers and decorative panel edges |
| `--surface-muted` | `#EFF0EA` | Grouped context, such as a responsibility row |
| `--surface-action` | `#EFEFE9` | A shared action area |
| `--red` | `#C62D31` | Brand, selected state, and main action |
| `--red-hover` | `#A92327` | Hover or pressed state of a red button |
| `--soft-red` | `#FBEFED` | Background for a relevant error or warning |
| `--focus` | `#9C2528` | Keyboard focus outline |
| `--control-border` | `#81877D` | Field or control edges that must be easy to identify |
| `--chart-primary` | `#393F3B` | Main data series |
| `--chart-secondary` | `#767D72` | A supporting data series |

### Accent placement

Keep red concentrated in three roles: brand, the current selection, and the main action. Repeated source links, secondary buttons, and ordinary icons stay neutral. A red warning must have a real cause and a written label.

Do not use red to imply that a politician, party, provider, or area is bad. A map's data color must follow a stated measure and a legend.

The map may contain muted colors from its base layer. Those colors do not become additional brand accents.

### Color rules

- Use one main accent across the app. Do not give every section a different color.
- Avoid blue-purple gradients, rainbow palettes, neon, background orbs, and page-wide glows.
- A gradient is allowed over a real image when it keeps text readable. Derive its dark tones from the image, as in the homepage's charcoal overlay.
- Use text and shape as well as color for status, selection, and errors.
- The light divider token is not sufficient when a control boundary carries essential information. Use `--control-border` for that case.
- Dark mode is outside this baseline. Add a toggle only with a complete, tested theme and a saved user choice.

## 4. Typography

Use sans-serif type throughout the product. Serif headings are outside Rantai's identity.

**Plus Jakarta Sans** gives headings a local identity and clear, open letter shapes. **Inter** supports compact navigation, tables, labels, and longer reading text. Use actual installed font weights and avoid synthetic intermediate weights.

| Role | Desktop | Small screen | Weight | Line height |
| --- | --- | --- | --- | --- |
| Feature heading over an image | 38–40 px | 28–36 px, fluid | 700 | 1.2 |
| Page title | 28–32 px | 26–28 px | 700 | 1.35 |
| Section heading | 24 px | 22–24 px | 700 | 1.35 |
| Story or panel heading | 20–24 px | 20–24 px | 600 or 700 | 1.4 |
| Reading text | 16 px | 16 px | 400 | 1.65 |
| Compact UI text | 14 px | 14 px | 400–600 | 1.5 |
| Navigation and buttons | 13–14 px | 13–14 px | 500 or 600 | 1.4 |
| Source notes and metadata | 12 px | 12 px | 400 or 500 | 1.5 |
| Key comparison value | 32–36 px | 32–36 px | 700 | 1.1 |

Use sentence case. Keep letter spacing close to normal. Slight negative spacing is allowed on large headings, within about `-0.03em`. Avoid widely spaced uppercase labels and monospace headings.

Use tabular numbers in aligned numeric columns. Keep the number and its unit together. For Indonesian UI, use Indonesian number formatting, such as `1.248` and `Rp2,4 miliar`, with the period or scope shown beside the figure. A shortened currency value must have its exact amount available in details.

Use a consistent date format. Formatting examples are `25 Okt 2026` in compact UI and `25 Oktober 2026` in reading text. Source timestamps include the time zone, such as `20.30 WIB`. Show the reporting period separately from the time the source was retrieved; only display actual retrieval times.

Source notes, data warnings, and instructions must remain at least 12 px. The original mockup has smaller captions; raise those when building the app. Tiny map attribution may follow the map provider's layout, with a readable expanded source view.

## 5. Spacing, grid, and page structure

Use the spacing scale `4, 8, 12, 16, 20, 24, 32, 40, 48, 64` px. Choose a smaller gap to join related content and a larger gap to separate decisions or sections.

### Shared page frame

Align the header and main content to the same frame. The large desktop wrapper is 1316 px including 40 px gutters, leaving 1236 px for content. Keep reading columns around 60–72 characters wide.

| Viewport width | Side gutter | Layout behavior |
| --- | --- | --- |
| Up to 650 px | 20 px | One main column; compact header and mobile menu |
| 651–850 px | 26 px | Compact navigation; stack lower editorial columns |
| 851–1100 px | 26 px | Two columns where content fits; reduce gaps |
| 1101–1439 px | 30 px | Standard desktop frame, up to 1280 px including gutters |
| 1440 px and wider | 40 px | Large frame, up to 1316 px including gutters |

The homepage has one full-row hero followed by region search and a short explanation. The photograph keeps its natural aspect ratio at every viewport; do not crop it with a fixed height or object-cover. Regional data belongs after the visitor selects an area and an information category.

### Page hierarchy

Begin with a useful title, the selected area, and the reporting period or source state where relevant. Show the content that supports the main decision next. Place the next action close to the evidence it uses.

Feature and area headers may use relevant photography or a real map. Data and workflow headers can be filled by their title, context, and controls. Do not add an empty hero placeholder or a large decorative banner to every screen.

Vary section composition when its content changes. Use rows, lists, reading columns, maps, and tables where they serve the task. Cards are for content that needs its own boundary or action.

## 6. Surfaces, corners, glass, and elevation

| Radius token | Value | Use |
| --- | --- | --- |
| `--radius-small` | 4 px | Small source and status labels |
| `--radius-control` | 8 px | Buttons, inputs, map pins |
| `--radius-group` | 10 px | Composite controls and grouped context |
| `--radius-panel` | 16 px | Map and action panels |
| `--radius-feature` | 20 px | Feature images and dialogs |

Round avatars to a circle. Avoid pill-shaped controls as a general style.

When two rounded edges follow the same corner, calculate their radii:

```text
outer radius = inner radius + inset
```

Examples: a 16 px panel with an 8 px inset uses an 8 px inner corner; a 20 px feature frame with a 12 px inset also uses an 8 px inner corner. Apply this to visually nested edges. A button far inside a panel does not need to trace the panel's corner.

Most surfaces have no shadow. Use a small shadow for controls that float over a map and a stronger shadow for an open popover or dialog. A border, spacing, or background should do the work for ordinary sections.

Glass is limited to the sticky navigation and, when present, the floating map legend. Match the homepage with an almost opaque light surface and moderate blur. Use a solid fallback when blur is unavailable. Keep dialog surfaces opaque and dim the background behind them.

No glow, decorative background grid, or endless moving texture belongs in the baseline.

## 7. Navigation and shared controls

### Header

Keep the Rantai wordmark on the left, the main destinations in the middle, and saved items, notifications, and the account control on the right. Use these labels consistently:

`Beranda` · `Wilayah` · `Tanya` · `Lapor` · `Tentang Data`

Use the same five destinations on phones, shortening `Tentang Data` to `Tentang`. `Wilayah` returns to the selected area or to search when no area is selected. Maps, comparative data, and monitoring remain available through contextual links and the footer.

Use the red underline for the current main destination. Set `aria-current="page"` for a page route. Do not mark the homepage as current on every route.

The header is about 83 px high on desktop, 73 px on compact layouts, and 66 px on phones. Provide scroll offset so sticky navigation does not cover a linked heading.

Every destination must exist. A feature outside the current release gets a clear `Segera hadir` explanation or is omitted from navigation.

### Location control

Use one shared location selector. Show the selected kecamatan and city; place the label `Wilayah yang kamu pantau` above the value for signed-in views.

Search supports a visible label, a clear current choice, and a useful no-results message. Changing area updates all dependent content. Keep the old area's values from appearing beneath a new area title while a request is loading.

### Buttons and links

- **Primary:** solid red for the screen's main action, such as `Buat paket aspirasi`.
- **Secondary:** neutral outline or a quiet surface for a supporting action.
- **Text link:** neutral text for a route or disclosure.
- **Icon control:** a relevant icon with an accessible name and visible focus.

Use 44 px minimum hit areas, including icon-only buttons. A compact visible icon can sit inside a larger hit area. Labels name the action: `Lihat sumber`, `Ubah wilayah`, `Salin draf`.

Use arrows only when direction clarifies the action. Use an external-link mark for an external destination. A save button uses a bookmark and an `aria-pressed` state. Do not add an arrow to every button.

Notification dots and counts reflect actual state. Do not invent unread activity to make a screen look busy.

### Brand mark and icons

Keep the lowercase `rantai.` wordmark, charcoal lettering, and red chain mark and period consistent with the homepage. Use the supplied mark rather than redrawing it differently on each screen.

Continue the mockup's SVG icon family: a 20 px default box, about 1.7 px strokes, rounded line ends, and no fill unless a selected state needs one. Smaller 16 px icons can sit beside text; keep the surrounding hit area large enough.

Choose glyphs for their function: a location pin, school, document, bookmark, or external link. Do not mix icon families or import a new default library without a clear visual reason. Avoid generic sparkles, lightning, robots, and decorative arrows.

## 8. Component contracts

### Source label and source details

Every metric has a nearby route to its source. A section-level source label may cover several values when the shared scope is clear. A mixed-source table needs source access at row or value level.

Use a small rectangular label with a 4 px radius, neutral background, and clear wording. Typical labels are `Data contoh`, `Belum diverifikasi`, and `Rilis resmi DPRD`. Avoid decorative dots, glows, or an `AI Powered` capsule.

A source view includes:

- Publisher and the specific document or dataset link.
- Geographic scope and reporting period.
- Publication date and actual retrieval time, when available.
- Verification status and missing fields.
- Calculation, unit, and any estimate used.
- Limits that affect how the result can be read.

For simulated values, say `Data contoh` beside the content and identify the simulation as the source. Retrieval time is `Tidak berlaku · simulasi`. Do not attach invented values to a real official source.

### Maps

Keep the map in a 16 px panel. Use quiet base colors, clear area selection, relevant markers, zoom controls, and a readable legend. Keep map attribution visible.

Area selection and a data measure are different visual states. Label them separately. A boundary must come from geographic data; an illustration must be marked as an illustration.

Provide `Lihat sebagai daftar` with the same records and source status. Keep selected records and filters in sync between map and list. Label simulated markers as examples and avoid linking them to real providers or schools.

### Charts and comparisons

Write the question the chart answers before choosing a chart. Put scope, period, units, and source access near it. Use direct labels and exact values where they help reading.

For school capacity, use two aligned horizontal bars with named measures. Show the arithmetic and explain limits. Missing or unverified student data must remain visible in the interpretation.

The neutral chart tokens suit separate bars with direct labels. If series touch, choose fills, outlines, or patterns that keep adjacent parts distinct. Do not assume that every pair of neutral tokens has enough contrast for stacked segments.

Avoid charts that only fill an empty space, fabricated trends, unlabelled estimates, and implied admission chances for individual children.

### Tables and data rows

Put the field that determines the next action early. For procurement, useful fields include package, location confidence, contract value, stage, and source. Include provider and maintenance terms when verified and relevant.

Right-align numeric columns, use tabular numbers, and keep headings plain. Show only actions that work. Do not add a three-dot menu to every row by default.

Keep the page free of horizontal overflow. On small screens, use a labelled scrolling table or a row layout that preserves field labels, source links, and record identity.

### News and representative profiles

Use the same profile layout across parties and members. Show verified name, dapil, fraksi, komisi, and source dates when available. Present official releases and independent media with different text labels and equal reading access.

Each news item shows its title, neutral one-sentence summary, publisher, date, source link, and access to a right of reply. Keep old releases visibly marked as archives.

Do not show politician scores, rankings, invented activity, or personal voting records inferred from a faction's position. Explain that coverage frequency does not measure performance.

### Forms and dialogs

Use visible field labels, clear units, honest placeholders, and nearby validation. Empty fields stay empty. Do not prefill made-up names, emails, or messages.

A standard dialog is up to 600 px wide, within the viewport, with a 20 px radius and about 24–28 px content padding. On phones, retain safe outer margins and a usable close control. Put long reading or complex forms on a page when a dialog becomes cramped.

Dialogs need a title, a close action, keyboard focus containment, Escape support where appropriate, and focus return to the opener. Announce save and error messages without moving focus unexpectedly.

### Aspirations and information requests

Keep the sequence clear: inspect evidence, choose the topic and recipient, review the draft, then copy or export it. State whether the draft has been sent. The current product flow lets the resident send it themselves.

Carry source details and simulation labels into every export. Do not add personal information the resident has not entered. Keep export controls labelled with the file type they actually produce.

## 9. Required states

| State | What the screen explains | Useful action |
| --- | --- | --- |
| First visit | No saved area has been selected | `Pilih wilayah` |
| No saved items | Nothing has been bookmarked yet | Explain the bookmark control |
| No filter results | The active filters found no matching records | `Hapus filter` or edit search |
| Loading | Which area's records or measure are being loaded | Preserve context and allow navigation |
| Source unavailable | Which source failed and whether a saved snapshot exists | Retry or open the snapshot source |
| Stale snapshot | The actual last retrieval time | Inspect source details or retry |
| Missing data | Which field or dataset is missing and how it limits reading | View the data gap or prepare a request |
| Permission required | Which content needs sign-in or another role | The appropriate sign-in or access action |
| Action failed | Which action failed and whether input was saved | Retry without re-entering valid input |
| Success | What changed and where the result is stored | Open the result when useful |

Use skeletons only when they match the content being loaded. A real loading indicator can run while a request is active. Decorative pulses and perpetual motion are excluded.

## 10. Page patterns across the app

| Area | Main question | Layout priority |
| --- | --- | --- |
| Signed-in Beranda | What should I inspect in my saved area? | Area control, a useful local question, map, related budget and oversight, next action |
| Wilayah | What do the available records show about this kecamatan? | Area and period, theme selector, map/list, source status, connected records |
| Wakilku | Who represents this area and what can be verified about their role? | Verified dapil mapping, profiles, source-labelled coverage, faction position, right of reply |
| Arus Anggaran | What package is connected to this place and what is known about it? | Area filters, decision-based rows, package detail, value and source, known delivery terms |
| Pendidikan | How do school capacity and demand compare on the available evidence? | Matched period, separate measures, source limits, capacity and related spending |
| Aksi | How can I turn this evidence into a clear question? | Evidence summary, recipient or channel, draft review, copy or export |
| Data yang Belum Dibuka | What is missing and how can it be requested? | Specific gap, effect on interpretation, PPID request template |
| Account and saved items | What have I chosen or saved? | Preferences and actual saved records with useful empty states |

Use the shared header and tokens across these patterns. Do not force a sidebar, four metric cards, chart, and table onto each page.

## 11. Copy and trust

Use Indonesian as the default UI language. Keep labels short, neutral, and consistent. English translations must preserve the same scope and uncertainty.

Prefer `menyuarakan`, `membahas`, `berdasarkan sumber`, `belum tersedia`, and `perlu ditanyakan`. Avoid praise, blame, or a conclusion that the evidence does not support.

Distinguish these concepts everywhere:

- A budget plan, a contract value, and actual spending.
- A reported activity and a verified outcome.
- A member's role and the position of their faction.
- An actual count, an estimate, missing data, and simulated data.

Education comparisons require verified student and seat data for the same period and clearly stated coverage. A local difference alone does not establish a child's chance of admission. Missing class-6 data remains a visible limitation until verified.

Provide a right of reply on scrutiny of members, institutions, and providers. Show replies with the item they address, verified attribution, and a date.

UI copy contains no decorative emoji, invented customer names, fake metrics, fabricated testimonials, or false live status.

## 12. Accessibility and motion

Set these as implementation targets:

- At least 4.5:1 contrast for normal text and 3:1 for large text.
- At least 3:1 for essential control boundaries and meaningful non-text marks against their surroundings.
- Clear non-color labels or shapes for state and data meaning.
- Visible keyboard focus, logical focus order, and a skip link.
- At least 44 × 44 px interaction areas for shared controls.
- Semantic headings, buttons, links, labels, and table headers.
- Text alternatives for maps and charts; meaningful image descriptions.
- Usable content at 200% zoom and at the target small-screen widths.
- No content or control hidden by sticky headers, clipping, or a mobile keyboard.

Use about 160 ms for hover, pressed, and state changes. Animate only when it explains a response. Avoid scroll reveals, floating elements, autoplay decoration, bouncing buttons, and stacked animation presets.

Honor `prefers-reduced-motion` by removing nonessential transitions and smooth scrolling. Keep essential loading state clear through text.

## 13. Starter tokens

Move these into the app's shared theme layer. Components should consume them rather than repeat local values.

```css
:root {
  --paper: #fafaf8;
  --white: #ffffff;
  --ink: #202323;
  --muted: #626664;
  --line: #e3e5e1;
  --surface-muted: #eff0ea;
  --surface-action: #efefe9;
  --red: #c62d31;
  --red-hover: #a92327;
  --soft-red: #fbefed;
  --focus: #9c2528;
  --control-border: #81877d;
  --chart-primary: #393f3b;
  --chart-secondary: #767d72;

  --heading: "Plus Jakarta Sans", Arial, sans-serif;
  --ui: "Inter", Arial, sans-serif;

  --space-1: 4px;
  --space-2: 8px;
  --space-3: 12px;
  --space-4: 16px;
  --space-5: 20px;
  --space-6: 24px;
  --space-8: 32px;
  --space-10: 40px;
  --space-12: 48px;
  --space-16: 64px;

  --radius-small: 4px;
  --radius-control: 8px;
  --radius-group: 10px;
  --radius-panel: 16px;
  --radius-feature: 20px;

  --page-max: 1316px;
  --content-max: 1236px;
  --dialog-max: 600px;
  --reading-max: 68ch;
  --target-size: 44px;
  --motion-fast: 160ms;

  --shadow-float: 0 2px 10px rgb(32 35 35 / 8%);
  --shadow-dialog: 0 20px 80px rgb(32 35 35 / 19%);
}

:where(button, a, input, select, textarea):focus-visible {
  outline: 3px solid var(--focus);
  outline-offset: 4px;
}

@media (prefers-reduced-motion: reduce) {
  html { scroll-behavior: auto; }
  *, *::before, *::after {
    animation: none !important;
    transition: none !important;
  }
}
```

When porting the mockup, preserve its composition while normalizing small text, hit areas, dialog corners, chart contrast, and field boundaries to this shared specification. These rules apply to new components and planned production cleanup; they are not a claim that the exported HTML has already been changed.

## 14. Design review before delivery

- [ ] The page has one clear job and a useful first action.
- [ ] Shared colors, fonts, spacing, and corner tokens are used.
- [ ] The area, period, and data state stay consistent across the page.
- [ ] Red is limited to meaningful brand, selection, action, or warning roles.
- [ ] Images identify their actual place; overlays keep text readable.
- [ ] Glass appears only where allowed; ordinary surfaces stay flat.
- [ ] Every metric has source access, scope, unit, and any needed warning.
- [ ] Simulated values and exports remain visibly labelled.
- [ ] Profiles, media labels, and replies follow the neutrality rules.
- [ ] Every visible control and navigation destination works.
- [ ] Empty, loading, missing, stale, error, and success states have useful copy.
- [ ] Layouts work at 360, 768, 1024, and 1440 px and at intermediate widths.
- [ ] Long titles, larger numbers, missing images, and no-results searches fit.
- [ ] Keyboard use, focus return, contrast, hit areas, and reduced motion are checked.
- [ ] The page avoids decorative emoji, status dots, fake activity, default bento grids, and repeated filler cards.

Record a new shared pattern under the relevant section of this file. For a local exception, record its purpose, where it applies, why existing tokens do not fit, and the accessibility checks completed. Reuse the pattern only after it has become part of the shared guide.

## 15. Implementation in this repository

The mockup `rantai-beranda.html` is kept outside the repository. Its production port is the homepage at `frontend/src/routes/[lang=lang]/+page.svelte`. Shared tokens live in `frontend/src/lib/styles/tokens.css` and are exposed to Tailwind in `frontend/tailwind.config.ts`.

| Spec token | Code token | Tailwind class |
| --- | --- | --- |
| `--paper` | `--bg` | `bg-page` |
| `--white` | `--surface`, `--raised` | `bg-surface`, `bg-raised` |
| `--ink` | `--text` | `text-fg` |
| `--muted` | `--muted` | `text-muted` |
| `--line` | `--line` | `border-line` |
| `--surface-muted` | `--sunken` | `bg-sunken` |
| `--surface-action` | `--secondary` | `bg-secondary` |
| `--red` | `--primary` | `bg-primary`, `text-primary` |
| `--red-hover` | `--primary-hover` | `hover:bg-primary-hover` |
| `--soft-red` | `--danger-bg` | `bg-danger-bg` |
| `--focus` | `--focus` | focus ring on `:focus-visible` |
| `--control-border` | `--line-strong` | `border-line-strong` |
| `--chart-primary` | `--chart` | `bg-chart`, `fill-chart` |
| `--chart-secondary` | `--chart-2` | `bg-chart-2` |
| `--heading` | `--font-heading` | `font-heading` |
| `--ui` | `--font-sans` | `font-sans` |
| `--radius-small`, `-control`, `-group`, `-panel`, `-feature` | `--r-xs`, `--r-md`, `--r-group`, `--r-lg`, `--r-xl` | `rounded-xs`, `rounded-md`, `rounded-group`, `rounded-lg`, `rounded-xl` |

Recorded decisions:

- **Brand mark.** `frontend/src/lib/components/LogoMark.svelte` is a vector trace of the supplied `logo-rantai.png`; the favicon is `frontend/static/favicon.svg`. The wordmark uses the `text-merek` type token.
- **Navigation labels.** The header uses the region-first destinations in §7. Area category navigation presents public services, government projects, and representatives/agencies.
- **Dark theme.** The earlier release already shipped a complete dark theme behind the header toggle with a saved choice (`localStorage` key `rantai-theme`). It is kept and retuned to these tokens. Light is the default and the theme no longer follows the system setting.
- **Map data colour.** Seat-gap classes use a slate sequential scale (`--gap-1` to `--gap-4`) with a legend. The selected area is a dashed red outline, a separate state from the data measure (§8 Maps).
- **Hero photo.** `frontend/static/foto/bundaran-hi.webp` (Unsplash) with the charcoal overlay `.bayang-foto`; the caption links to the source.

## UX revision · 27 October 2026

- The homepage explains a concrete benefit in one full-row photo hero, with a link to region search below it. Render the complete photograph at its native ratio, including on phones. Keep the caption in normal flow on phones where an overlay would not fit.
- Search accepts Jakarta kelurahan and district names. Show its scope and three-character minimum. Enter selects a unique exact match; ambiguous results require an explicit choice.
- Onboarding is optional and introduces the region task. Do not interrupt a first visit with an automatic modal.
- The homepage offers an explicit example link without showing sample regional statistics. The area landing page asks the visitor to choose one information category before displaying its data.
- Carry explicit district context through the URL when opening Ask, Check My School, and aspiration drafts. For schools searched from an area's centre, explain the origin of the distance calculation.
- Education starts with seats per 100 prospective graduates, its interpretation limit, then detailed values and sources. Keep related actions near the summary. Representative activity is available through a disclosure after the education summaries.
- Source controls retain certainty labels and accessible source descriptions. Use a complete short label and a 44 px hit area on narrow numeric cards.
- Project chain status describes data availability. Explain the distinction between missing data and the project's physical progress.
- Aspiration documents show that they have not been sent, identify the responsible agency, and expose existing delivery instructions outside the document. Links to channel information must not be labelled as submission buttons.

## Region-first information architecture · 27 October 2026

1. Homepage: purpose, complete hero photograph, area search, and a brief explanation. Do not fetch or render regional statistics, a thematic map, or project lists here.
2. Area landing page: selected area, change-area control, and three information choices: `Layanan publik`, `Proyek pemerintah`, `Wakil rakyat & instansi`. A missing or invalid `bagian` query shows these choices.
3. Category view: explicitly select a `bagian`; preserve area context and provide a return to the information choices. Services display one education level at a time.
4. Explain services in the order: condition, meaning for residents, evidence, responsible agency/oversight, related spending, then map/table and next actions.
5. Label a measured shortage as `Perlu perhatian`. Adequate capacity receives a factual explanation of its benefit without an issue badge. Missing or invalid data is `Belum dapat dinilai`, never inferred to be a good condition.
6. A project overview identifies its recorded agency, contractor, planned dates, recorded procurement stage, available handover date, and location precision. Missing structured outcomes remain unknown. Do not infer work duration or completion from a budget or tender stage.
7. Agency delivery, contractor work, and council oversight are distinct roles. Do not call the current DPRD directory a regional cabinet.
