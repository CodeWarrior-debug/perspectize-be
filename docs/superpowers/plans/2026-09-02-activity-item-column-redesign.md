# Activity Table — Item Column Redesign Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Redesign the Activity table's Item column (thumbnail + title, with split click behavior and a details modal) and add a mobile stacked-card layout below 860px, per the design handoff.

**Architecture:** A new imperative AG Grid cell renderer (`activityItemCellRenderer`) replaces `itemCellRenderer` on the `item` column only. It reads an `onOpenDetails` callback off `params.context` (same pattern already used for `perspectivesByContentId`). `ActivityTable.svelte` grows a `cardMode` boolean (window width < 860px, via `matchMedia`, mirroring the existing `responsiveTier` effect) that swaps the whole AG Grid for a new `ActivityCardList.svelte` component. Both the cell renderer and the card list open the same new `ActivityDetailsModal.svelte`, a shadcn `Dialog` styled to match the mock. No backend changes — the two stats without a backing aggregate (Perspectives count, Avg. Rating) render as `—` placeholders.

**Tech Stack:** SvelteKit, Svelte 5 runes, `ag-grid-svelte5` (`@ag-grid-community/*`), Tailwind CSS v4, shadcn-svelte (`bits-ui` `Dialog`), `@lucide/svelte` icons, Vitest + `@testing-library/svelte`.

**Spec:** `/Users/jamesjordan/Downloads/design_handoff_activity_item_column/README.md` (plus `ActivityTable.dc.html` prototype and the three screenshots in that folder).

## Global Constraints

- **Scope:** Item column cell + its click behavior + the new details modal + the mobile card layout that replaces the table below 860px, and the row-height change to 112px for every row. Every other column (Type, Length, Views, Likes, Date, Channel, Tags), the top nav, and the page header are unchanged — do not touch `Header.svelte` or `+page.svelte`.
- **Fidelity:** Recreate pixel-for-pixel using the app's real design tokens from `frontend/src/app.css` (`--color-primary: #1a365d`, `--color-border: #d4d4d4`, `--color-muted-foreground: #525252`, `--color-accent: #f7fafc`, `--font-family-serif: 'Charter', 'Georgia', 'Times New Roman', serif`, `--default-font-family: 'Geist', ...`), not the mock's literal `#1a365d` / Georgia / `#f7fafc` values.
- **Thumbnail:** 104×58px, `border-radius: 6px` (`rounded-md`), `object-fit: cover`.
- **Title:** 4px gap below thumbnail, centered, max 2 lines (`line-clamp-2`), `font-size: 13px`, `line-height: 1.4`, serif, `max-width: 184px`.
- **Row height:** 112px for the whole table (was 44px).
- **Breakpoint:** 860px window width — below it, cards; at/above it, the existing AG Grid table (with its existing column-tier responsive behavior untouched).
- **Thumbnail click:** `stopPropagation()`, then `window.open(videoUrl, '_blank', 'noopener,noreferrer')`. Must not also open the modal.
- **Cell click (elsewhere):** opens the details modal for that row.
- **Modal:** backdrop click closes; click inside the modal card does not close it; × closes it. "Update metadata" button is a non-functional placeholder (no click handler).
- **Stat grid placeholders:** Perspectives count and Avg. Rating render as `—` (no backend aggregate exists yet — tracked in issue #286 as a follow-up, not built here).

---

## File Structure

- `frontend/src/lib/utils/activityItemCellRenderer.ts` — **new.** AG Grid `ICellRendererFunc` for the redesigned Item cell (thumbnail + title, two click zones).
- `frontend/tests/unit/activityItemCellRenderer.test.ts` — **new.** Unit tests for the renderer function.
- `frontend/src/lib/components/ActivityDetailsModal.svelte` — **new.** The details modal (shadcn `Dialog`).
- `frontend/tests/components/ActivityDetailsModal.test.ts` — **new.**
- `frontend/src/lib/components/ActivityCardList.svelte` — **new.** Mobile stacked-card list, one row per card.
- `frontend/tests/components/ActivityCardList.test.ts` — **new.**
- `frontend/src/lib/components/ActivityTable.svelte` — **modify.** Wire in the new renderer, row height, `cardMode` breakpoint, and the modal.
- `frontend/tests/components/ActivityTable.test.ts` — **modify.** Add card/grid breakpoint-switch coverage.

---

### Task 1: Item cell renderer (thumbnail + title, two click zones)

**Files:**
- Create: `frontend/src/lib/utils/activityItemCellRenderer.ts`
- Test: `frontend/tests/unit/activityItemCellRenderer.test.ts`

**Interfaces:**
- Consumes: `extractVideoIdFromUrl` from `frontend/src/lib/utils/formatting.ts` (existing, signature `(url: string | null) => string | null`).
- Produces: `activityItemCellRenderer(params: ActivityItemCellRendererParams): HTMLElement | string` — default export is a **named** export `activityItemCellRenderer`. `ActivityItemCellRendererParams` shape:
  ```ts
  interface ActivityItemCellRendererParams {
  	data?: { id: string | number; name: string; url: string | null };
  	context?: { onOpenDetails?: (contentId: string) => void };
  }
  ```
  Later tasks (Task 4) rely on this exact param shape and on the `context.onOpenDetails` contract.

- [ ] **Step 1: Write the failing tests**

```ts
// frontend/tests/unit/activityItemCellRenderer.test.ts
import { describe, it, expect, vi } from 'vitest';
import { activityItemCellRenderer } from '$lib/utils/activityItemCellRenderer';

describe('activityItemCellRenderer', () => {
	it('returns empty string when data is missing', () => {
		expect(activityItemCellRenderer({})).toBe('');
	});

	it('renders a thumbnail and a 2-line-clamp title', () => {
		const result = activityItemCellRenderer({
			data: { id: '1', name: 'My Video', url: 'https://youtube.com/watch?v=abc123' },
		}) as HTMLElement;

		expect(result).toBeInstanceOf(HTMLDivElement);

		const img = result.querySelector('img');
		expect(img).toBeTruthy();
		expect(img?.src).toBe('https://i.ytimg.com/vi/abc123/hqdefault.jpg');

		const title = result.querySelector('[data-testid="item-title"]');
		expect(title?.textContent).toBe('My Video');
		expect(title?.className).toContain('line-clamp-2');
	});

	it('omits the thumbnail image when the URL has no extractable video id', () => {
		const result = activityItemCellRenderer({
			data: { id: '1', name: 'No Thumb', url: 'https://example.com/not-youtube' },
		}) as HTMLElement;

		expect(result.querySelector('img')).toBeNull();
		expect(result.querySelector('[data-testid="item-title"]')?.textContent).toBe('No Thumb');
	});

	it('clicking the thumbnail opens the video in a new tab and does not open the modal', () => {
		const onOpenDetails = vi.fn();
		const openSpy = vi.spyOn(window, 'open').mockImplementation(() => null);

		const result = activityItemCellRenderer({
			data: { id: '42', name: 'My Video', url: 'https://youtube.com/watch?v=abc123' },
			context: { onOpenDetails },
		}) as HTMLElement;

		const thumbWrap = result.querySelector('[data-testid="item-thumb"]') as HTMLElement;
		thumbWrap.click();

		expect(openSpy).toHaveBeenCalledWith(
			'https://youtube.com/watch?v=abc123',
			'_blank',
			'noopener,noreferrer',
		);
		expect(onOpenDetails).not.toHaveBeenCalled();

		openSpy.mockRestore();
	});

	it('clicking anywhere else in the cell opens the details modal with the row id', () => {
		const onOpenDetails = vi.fn();

		const result = activityItemCellRenderer({
			data: { id: '42', name: 'My Video', url: 'https://youtube.com/watch?v=abc123' },
			context: { onOpenDetails },
		}) as HTMLElement;

		result.click();

		expect(onOpenDetails).toHaveBeenCalledWith('42');
	});

	it('sets accessible titles for both click zones', () => {
		const result = activityItemCellRenderer({
			data: { id: '1', name: 'My Video', url: 'https://youtube.com/watch?v=abc123' },
		}) as HTMLElement;

		expect(result.title).toBe('View content data + details');
		expect(result.querySelector('[data-testid="item-thumb"]')?.getAttribute('title')).toBe(
			'Open original content in new tab',
		);
	});
});
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd frontend && pnpm exec vitest run --project unit tests/unit/activityItemCellRenderer.test.ts`
Expected: FAIL — `Cannot find module '$lib/utils/activityItemCellRenderer'`

- [ ] **Step 3: Implement the renderer**

```ts
// frontend/src/lib/utils/activityItemCellRenderer.ts
import { extractVideoIdFromUrl } from './formatting';

export interface ActivityItemCellRendererParams {
	data?: { id: string | number; name: string; url: string | null };
	context?: { onOpenDetails?: (contentId: string) => void };
}

const PLAY_ICON_SVG = `<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#ffffff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="5 3 19 12 5 21 5 3"></polygon></svg>`;

/**
 * AG Grid cell renderer for the redesigned Item column: a 104x58 thumbnail
 * (click -> open video in a new tab) stacked above a 2-line title
 * (click -> open the details modal). Two independent click zones.
 */
export function activityItemCellRenderer(
	params: ActivityItemCellRendererParams,
): HTMLElement | string {
	if (!params.data) return '';

	const { id, name, url } = params.data;
	const onOpenDetails = params.context?.onOpenDetails;

	const cell = document.createElement('div');
	cell.className =
		'group/cell flex h-full w-full flex-col items-center justify-center gap-1 px-2.5 py-1.5 cursor-pointer';
	cell.title = 'View content data + details';
	cell.addEventListener('click', () => {
		onOpenDetails?.(String(id));
	});

	const thumbWrap = document.createElement('div');
	thumbWrap.dataset.testid = 'item-thumb';
	thumbWrap.className =
		'group/thumb relative h-[58px] w-[104px] flex-none overflow-hidden rounded-md bg-muted';
	thumbWrap.title = 'Open original content in new tab';
	thumbWrap.addEventListener('click', (e) => {
		e.stopPropagation();
		if (url) window.open(url, '_blank', 'noopener,noreferrer');
	});

	const videoId = extractVideoIdFromUrl(url);
	if (videoId) {
		const img = document.createElement('img');
		img.src = `https://i.ytimg.com/vi/${videoId}/hqdefault.jpg`;
		img.alt = '';
		img.className = 'h-full w-full object-cover';
		thumbWrap.appendChild(img);
	}

	const overlay = document.createElement('div');
	overlay.className =
		'pointer-events-none absolute inset-0 flex flex-col items-center justify-center gap-1 bg-[rgba(23,23,23,0.55)] opacity-0 transition-opacity group-hover/thumb:opacity-100';
	overlay.innerHTML = PLAY_ICON_SVG;
	const watchLabel = document.createElement('span');
	watchLabel.className = 'text-[11px] font-semibold text-white';
	watchLabel.textContent = 'Watch';
	overlay.appendChild(watchLabel);
	thumbWrap.appendChild(overlay);

	const title = document.createElement('div');
	title.dataset.testid = 'item-title';
	title.className =
		'mt-1 line-clamp-2 max-w-[184px] text-center font-[family-name:var(--font-family-serif)] text-[13px] leading-[1.4] text-foreground decoration-primary/30 group-hover/cell:underline';
	title.textContent = name;

	cell.appendChild(thumbWrap);
	cell.appendChild(title);

	return cell;
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd frontend && pnpm exec vitest run --project unit tests/unit/activityItemCellRenderer.test.ts`
Expected: PASS (6 tests)

- [ ] **Step 5: Commit**

```bash
git add frontend/src/lib/utils/activityItemCellRenderer.ts frontend/tests/unit/activityItemCellRenderer.test.ts
git commit -m "feat: add redesigned Item column cell renderer"
```

---

### Task 2: Details modal

**Files:**
- Create: `frontend/src/lib/components/ActivityDetailsModal.svelte`
- Test: `frontend/tests/components/ActivityDetailsModal.test.ts`

**Interfaces:**
- Consumes: `Dialog`, `DialogContent`, `DialogTitle`, `DialogClose` from `$lib/components/shadcn` (existing barrel); `formatCount`, `formatDate`, `formatDuration`, `formatTags`, `extractVideoIdFromUrl` from `$lib/utils/formatting` (all existing, already used elsewhere in this file).
- Produces: default-exported Svelte component with props:
  ```ts
  {
  	content: {
  		id: string;
  		name: string;
  		url: string | null;
  		channelTitle: string | null;
  		viewCount: number | null;
  		likeCount: number | null;
  		length: number | null;
  		lengthUnits: string | null;
  		publishedAt: string | null;
  		updatedAt: string;
  		tags: string[] | null;
  	} | null;
  	open: boolean;
  	onClose: () => void;
  }
  ```
  Task 4 renders this with a `ContentItem` (superset of the fields above — extra fields are ignored) and wires `open`/`onClose` to its own state.

- [ ] **Step 1: Write the failing tests**

```ts
// frontend/tests/components/ActivityDetailsModal.test.ts
import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/svelte';
import ActivityDetailsModal from '$lib/components/ActivityDetailsModal.svelte';

const content = {
	id: '42',
	name: 'Stephen Paea breaking bench',
	url: 'https://youtube.com/watch?v=abc123',
	channelTitle: 'TBD tribute',
	viewCount: 1300000,
	likeCount: 26500,
	length: 59,
	lengthUnits: 'seconds',
	publishedAt: '2026-02-24T00:00:00Z',
	updatedAt: '2026-03-01T00:00:00Z',
	tags: ['tom brady', 'tom brady goat'],
};

describe('ActivityDetailsModal', () => {
	it('renders nothing when closed', () => {
		render(ActivityDetailsModal, { props: { content, open: false, onClose: vi.fn() } });
		expect(screen.queryByText(content.name)).not.toBeInTheDocument();
	});

	it('renders the video title, channel, link, and stats when open', () => {
		render(ActivityDetailsModal, { props: { content, open: true, onClose: vi.fn() } });

		expect(screen.getByText(content.name)).toBeInTheDocument();
		expect(screen.getByText('TBD tribute')).toBeInTheDocument();
		expect(screen.getByText(content.url)).toBeInTheDocument();
		expect(screen.getByText('1.3 M')).toBeInTheDocument(); // views
		expect(screen.getByText('26.5 K')).toBeInTheDocument(); // likes
		expect(screen.getByText('0:59')).toBeInTheDocument(); // duration
	});

	it('shows placeholder stats for perspectives and avg rating', () => {
		render(ActivityDetailsModal, { props: { content, open: true, onClose: vi.fn() } });

		expect(screen.getByText('Perspectives')).toBeInTheDocument();
		expect(screen.getByText('Avg. Rating')).toBeInTheDocument();
		expect(screen.getAllByText('—').length).toBeGreaterThanOrEqual(2);
	});

	it('renders tags when present', () => {
		render(ActivityDetailsModal, { props: { content, open: true, onClose: vi.fn() } });
		expect(screen.getByText('tom brady, tom brady goat')).toBeInTheDocument();
	});

	it('omits the tags section when there are none', () => {
		render(ActivityDetailsModal, {
			props: { content: { ...content, tags: null }, open: true, onClose: vi.fn() },
		});
		expect(screen.queryByText('Tags')).not.toBeInTheDocument();
	});

	it('has a non-functional Update metadata button', async () => {
		render(ActivityDetailsModal, { props: { content, open: true, onClose: vi.fn() } });
		const button = screen.getByRole('button', { name: 'Update metadata' });
		expect(button).toBeInTheDocument();
		await fireEvent.click(button); // should not throw
	});

	it('calls onClose when the close button is clicked', async () => {
		const onClose = vi.fn();
		render(ActivityDetailsModal, { props: { content, open: true, onClose } });

		await fireEvent.click(screen.getByRole('button', { name: /close/i }));
		expect(onClose).toHaveBeenCalled();
	});
});
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd frontend && pnpm exec vitest run --project unit tests/components/ActivityDetailsModal.test.ts`
Expected: FAIL — `Failed to resolve import "$lib/components/ActivityDetailsModal.svelte"`

- [ ] **Step 3: Implement the component**

```svelte
<!-- frontend/src/lib/components/ActivityDetailsModal.svelte -->
<script lang="ts">
	import ExternalLinkIcon from '@lucide/svelte/icons/external-link';
	import XIcon from '@lucide/svelte/icons/x';
	import { Dialog, DialogContent, DialogTitle, DialogClose } from '$lib/components/shadcn';
	import {
		formatCount,
		formatDate,
		formatDuration,
		formatTags,
		extractVideoIdFromUrl,
	} from '$lib/utils/formatting';

	interface ModalContent {
		id: string;
		name: string;
		url: string | null;
		channelTitle: string | null;
		viewCount: number | null;
		likeCount: number | null;
		length: number | null;
		lengthUnits: string | null;
		publishedAt: string | null;
		updatedAt: string;
		tags: string[] | null;
	}

	let {
		content,
		open = false,
		onClose,
	}: {
		content: ModalContent | null;
		open?: boolean;
		onClose: () => void;
	} = $props();

	const videoId = $derived(content ? extractVideoIdFromUrl(content.url) : null);
	const thumbSrc = $derived(videoId ? `https://i.ytimg.com/vi/${videoId}/hqdefault.jpg` : null);
	const hasTags = $derived(!!content?.tags && content.tags.length > 0);

	function handleOpenChange(next: boolean) {
		if (!next) onClose();
	}
</script>

{#if content}
	<Dialog {open} onOpenChange={handleOpenChange}>
		<DialogContent
			showCloseButton={false}
			class="max-w-[560px] gap-0 overflow-hidden rounded-xl p-0"
		>
			<div class="flex items-start justify-between gap-3 bg-primary px-[22px] py-[18px]">
				<DialogTitle
					class="text-xs font-semibold tracking-wide text-primary-foreground/70 uppercase"
				>
					YouTube Video
				</DialogTitle>
				<DialogClose class="text-primary-foreground/80 hover:text-primary-foreground">
					<XIcon class="size-[18px]" />
					<span class="sr-only">Close</span>
				</DialogClose>
			</div>

			<div class="max-h-[calc(88vh-56px)] overflow-y-auto p-[22px]">
				<div class="flex items-start gap-3.5">
					<div class="h-[68px] w-[120px] flex-none overflow-hidden rounded-md bg-muted">
						{#if thumbSrc}
							<img src={thumbSrc} alt="" class="h-full w-full object-cover" />
						{/if}
					</div>
					<div class="min-w-0">
						<div class="font-[family-name:var(--font-family-serif)] text-[17px] font-bold leading-tight text-foreground">
							{content.name}
						</div>
						<div class="mt-1 text-[13px] text-muted-foreground">{content.channelTitle}</div>
					</div>
				</div>

				{#if content.url}
					<a
						href={content.url}
						target="_blank"
						rel="noopener noreferrer"
						class="mt-3.5 flex items-center gap-1.5 text-[13px] break-all text-primary no-underline hover:underline"
					>
						<ExternalLinkIcon class="size-3.5 flex-none" />
						<span>{content.url}</span>
					</a>
				{/if}

				<div class="mt-4.5 grid grid-cols-2 gap-2.5">
					<div class="rounded-lg border border-border bg-accent px-3 py-2.5">
						<div class="text-[11px] tracking-wide text-muted-foreground uppercase">
							Perspectives
						</div>
						<div class="mt-0.5 font-[family-name:var(--font-family-serif)] text-lg font-bold text-foreground">
							—
						</div>
					</div>
					<div class="rounded-lg border border-border bg-accent px-3 py-2.5">
						<div class="text-[11px] tracking-wide text-muted-foreground uppercase">
							Avg. Rating
						</div>
						<div class="mt-0.5 font-[family-name:var(--font-family-serif)] text-lg font-bold text-foreground">
							—
						</div>
					</div>
					<div class="rounded-lg border border-border bg-accent px-3 py-2.5">
						<div class="text-[11px] tracking-wide text-muted-foreground uppercase">Views</div>
						<div class="mt-0.5 font-[family-name:var(--font-family-serif)] text-[15px] font-bold text-foreground">
							{formatCount(content.viewCount)}
						</div>
					</div>
					<div class="rounded-lg border border-border bg-accent px-3 py-2.5">
						<div class="text-[11px] tracking-wide text-muted-foreground uppercase">Likes</div>
						<div class="mt-0.5 font-[family-name:var(--font-family-serif)] text-[15px] font-bold text-foreground">
							{formatCount(content.likeCount)}
						</div>
					</div>
					<div class="rounded-lg border border-border bg-accent px-3 py-2.5">
						<div class="text-[11px] tracking-wide text-muted-foreground uppercase">Duration</div>
						<div class="mt-0.5 font-[family-name:var(--font-family-serif)] text-[15px] font-bold text-foreground">
							{formatDuration(content.length, content.lengthUnits)}
						</div>
					</div>
					<div class="rounded-lg border border-border bg-accent px-3 py-2.5">
						<div class="text-[11px] tracking-wide text-muted-foreground uppercase">Published</div>
						<div class="mt-0.5 font-[family-name:var(--font-family-serif)] text-[15px] font-bold text-foreground">
							{content.publishedAt ? formatDate(content.publishedAt) : '—'}
						</div>
					</div>
				</div>

				<div
					class="mt-3.5 flex items-center justify-between border-t border-border pt-3.5"
				>
					<div>
						<div class="text-[11px] tracking-wide text-muted-foreground uppercase">
							Last updated in Perspectize
						</div>
						<div class="mt-0.5 font-[family-name:var(--font-family-serif)] text-sm text-foreground">
							{formatDate(content.updatedAt)}
						</div>
					</div>
					<button
						type="button"
						class="rounded-md border border-primary px-3.5 py-2 text-[13px] font-semibold text-primary hover:bg-primary/5"
					>
						Update metadata
					</button>
				</div>

				{#if hasTags}
					<div class="mt-3.5">
						<div class="mb-1.5 text-[11px] tracking-wide text-muted-foreground uppercase">
							Tags
						</div>
						<div class="font-[family-name:var(--font-family-serif)] text-[13px] text-foreground">
							{formatTags(content.tags)}
						</div>
					</div>
				{/if}
			</div>
		</DialogContent>
	</Dialog>
{/if}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd frontend && pnpm exec vitest run --project unit tests/components/ActivityDetailsModal.test.ts`
Expected: PASS (7 tests). If `DialogClose`'s accessible name isn't matched by `/close/i`, check the rendered button — bits-ui's `Dialog.Close` accepts arbitrary children; the `<span class="sr-only">Close</span>` inside it should already satisfy `getByRole('button', { name: /close/i })`.

- [ ] **Step 5: Commit**

```bash
git add frontend/src/lib/components/ActivityDetailsModal.svelte frontend/tests/components/ActivityDetailsModal.test.ts
git commit -m "feat: add Activity item details modal"
```

---

### Task 3: Mobile card list

**Files:**
- Create: `frontend/src/lib/components/ActivityCardList.svelte`
- Test: `frontend/tests/components/ActivityCardList.test.ts`

**Interfaces:**
- Consumes: `extractVideoIdFromUrl`, `formatDuration` from `$lib/utils/formatting`.
- Produces: default-exported Svelte component with props:
  ```ts
  {
  	rowData: Array<{
  		id: string | number;
  		name: string;
  		url: string | null;
  		channelTitle: string | null;
  		length: number | null;
  		lengthUnits: string | null;
  	}>;
  	onOpenDetails: (contentId: string) => void;
  }
  ```
  Root element carries `data-testid="activity-card-list"` — Task 4 asserts on this to confirm the card layout (vs. the grid) is showing.

- [ ] **Step 1: Write the failing tests**

```ts
// frontend/tests/components/ActivityCardList.test.ts
import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/svelte';
import ActivityCardList from '$lib/components/ActivityCardList.svelte';

const rowData = [
	{
		id: '1',
		name: 'Jordan Peterson: "Why Some People Never Change"',
		url: 'https://youtube.com/watch?v=abc123',
		channelTitle: 'Jordan Peterson',
		length: 2955,
		lengthUnits: 'seconds',
	},
	{
		id: '2',
		name: 'Stephen Paea breaking bench',
		url: null,
		channelTitle: 'TBD tribute',
		length: 59,
		lengthUnits: 'seconds',
	},
];

describe('ActivityCardList', () => {
	it('renders one card per row with title, channel, and duration', () => {
		render(ActivityCardList, { props: { rowData, onOpenDetails: vi.fn() } });

		expect(screen.getByText(rowData[0].name)).toBeInTheDocument();
		expect(screen.getByText('Jordan Peterson')).toBeInTheDocument();
		expect(screen.getByText('49:15')).toBeInTheDocument();
		expect(screen.getByText(rowData[1].name)).toBeInTheDocument();
	});

	it('has the activity-card-list test id at its root', () => {
		render(ActivityCardList, { props: { rowData, onOpenDetails: vi.fn() } });
		expect(screen.getByTestId('activity-card-list')).toBeInTheDocument();
	});

	it('opens the video in a new tab when the thumbnail is clicked, without opening details', () => {
		const onOpenDetails = vi.fn();
		const openSpy = vi.spyOn(window, 'open').mockImplementation(() => null);

		render(ActivityCardList, { props: { rowData, onOpenDetails } });
		fireEvent.click(screen.getByTestId('card-thumb-1'));

		expect(openSpy).toHaveBeenCalledWith(
			'https://youtube.com/watch?v=abc123',
			'_blank',
			'noopener,noreferrer',
		);
		expect(onOpenDetails).not.toHaveBeenCalled();
		openSpy.mockRestore();
	});

	it('opens details when the title/body area is clicked', async () => {
		const onOpenDetails = vi.fn();
		render(ActivityCardList, { props: { rowData, onOpenDetails } });

		await fireEvent.click(screen.getByText(rowData[0].name));
		expect(onOpenDetails).toHaveBeenCalledWith('1');
	});
});
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd frontend && pnpm exec vitest run --project unit tests/components/ActivityCardList.test.ts`
Expected: FAIL — `Failed to resolve import "$lib/components/ActivityCardList.svelte"`

- [ ] **Step 3: Implement the component**

```svelte
<!-- frontend/src/lib/components/ActivityCardList.svelte -->
<script lang="ts">
	import PlayIcon from '@lucide/svelte/icons/play';
	import { extractVideoIdFromUrl, formatDuration } from '$lib/utils/formatting';

	interface CardRow {
		id: string | number;
		name: string;
		url: string | null;
		channelTitle: string | null;
		length: number | null;
		lengthUnits: string | null;
	}

	let {
		rowData,
		onOpenDetails,
	}: {
		rowData: CardRow[];
		onOpenDetails: (contentId: string) => void;
	} = $props();

	function thumbSrc(row: CardRow): string | null {
		const videoId = extractVideoIdFromUrl(row.url);
		return videoId ? `https://i.ytimg.com/vi/${videoId}/hqdefault.jpg` : null;
	}

	function handleThumbClick(row: CardRow, e: MouseEvent) {
		e.stopPropagation();
		if (row.url) window.open(row.url, '_blank', 'noopener,noreferrer');
	}
</script>

<div data-testid="activity-card-list" class="flex flex-col gap-2.5 px-2 py-2">
	{#each rowData as row (row.id)}
		<div
			class="flex items-center gap-3 rounded-lg border border-border bg-card p-2.5 hover:bg-primary/[0.06]"
		>
			<button
				type="button"
				data-testid={`card-thumb-${row.id}`}
				title="Open original content in new tab"
				class="relative h-16 w-24 flex-none overflow-hidden rounded-md bg-muted"
				onclick={(e) => handleThumbClick(row, e)}
			>
				{#if thumbSrc(row)}
					<img src={thumbSrc(row)} alt="" class="h-full w-full object-cover" />
				{/if}
				<span
					class="absolute right-1 bottom-1 flex items-center justify-center rounded bg-[rgba(23,23,23,0.65)] p-1"
				>
					<PlayIcon class="size-2.5 fill-white text-white" />
				</span>
			</button>

			<button
				type="button"
				title="View content data + details"
				class="min-w-0 flex-1 text-left"
				onclick={() => onOpenDetails(String(row.id))}
			>
				<div
					class="line-clamp-2 font-[family-name:var(--font-family-serif)] text-sm font-semibold leading-tight text-foreground"
				>
					{row.name}
				</div>
				<div class="mt-1.5 flex items-center gap-2 text-xs text-muted-foreground">
					{#if row.channelTitle}
						<span>{row.channelTitle}</span>
						<span>&middot;</span>
					{/if}
					<span>{formatDuration(row.length, row.lengthUnits)}</span>
				</div>
			</button>
		</div>
	{/each}
</div>
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd frontend && pnpm exec vitest run --project unit tests/components/ActivityCardList.test.ts`
Expected: PASS (4 tests)

- [ ] **Step 5: Commit**

```bash
git add frontend/src/lib/components/ActivityCardList.svelte frontend/tests/components/ActivityCardList.test.ts
git commit -m "feat: add Activity table mobile card list"
```

---

### Task 4: Wire it all into ActivityTable.svelte

**Files:**
- Modify: `frontend/src/lib/components/ActivityTable.svelte`
- Modify: `frontend/tests/components/ActivityTable.test.ts`

**Interfaces:**
- Consumes: `activityItemCellRenderer` (Task 1), `ActivityDetailsModal` (Task 2), `ActivityCardList` (Task 3) — exact exports as defined above.
- Produces: no new exports; this is the integration point.

- [ ] **Step 1: Write the failing tests**

Add to the bottom of `frontend/tests/components/ActivityTable.test.ts` (inside the existing `describe('ActivityTable', ...)` block, after the last `it(...)`):

```ts
	it('shows the mobile card list below the 860px breakpoint', async () => {
		mockRequest.mockResolvedValue(mockDataResponse);
		const matchMediaSpy = vi.spyOn(window, 'matchMedia').mockImplementation((query: string) => ({
			matches: query === '(max-width: 859px)',
			media: query,
			onchange: null,
			addListener: vi.fn(),
			removeListener: vi.fn(),
			addEventListener: vi.fn(),
			removeEventListener: vi.fn(),
			dispatchEvent: vi.fn(),
		})) as unknown as typeof window.matchMedia;

		const { container } = renderWithQuery();

		await waitFor(() => {
			expect(container.querySelector('[data-testid="activity-card-list"]')).toBeTruthy();
			expect(container.querySelector('[data-testid="ag-grid-container"]')).toBeFalsy();
		});

		matchMediaSpy.mockRestore();
	});

	it('shows the AG Grid table at/above the 860px breakpoint', async () => {
		mockRequest.mockResolvedValue(mockDataResponse);
		// Default window.matchMedia mock from tests/setup.ts always returns matches: false.
		const { container } = renderWithQuery();

		await waitFor(() => {
			expect(container.querySelector('[data-testid="ag-grid-container"]')).toBeTruthy();
			expect(container.querySelector('[data-testid="activity-card-list"]')).toBeFalsy();
		});
	});
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd frontend && pnpm exec vitest run --project unit tests/components/ActivityTable.test.ts`
Expected: FAIL — no element matches `[data-testid="activity-card-list"]` or `[data-testid="ag-grid-container"]` (neither exists yet).

- [ ] **Step 3: Modify ActivityTable.svelte**

3a. Add imports (near the existing component imports, e.g. after the `PerspectivePopover` import at `frontend/src/lib/components/ActivityTable.svelte:65`):

```ts
	import ActivityDetailsModal from '$lib/components/ActivityDetailsModal.svelte';
	import ActivityCardList from '$lib/components/ActivityCardList.svelte';
	import { activityItemCellRenderer } from '$lib/utils/activityItemCellRenderer';
```

3b. Add card-mode + details-modal state (near the existing `popoverOpen` state declarations at the top of `<script>`):

```ts
	// Details modal state (Item cell click -> metadata modal)
	let detailsModalContentId = $state<string | null>(null);
	const detailsModalContent = $derived(
		rowData.find((item) => String(item.id) === detailsModalContentId) ?? null,
	);

	function handleOpenDetails(contentId: string) {
		detailsModalContentId = contentId;
	}
	function handleCloseDetails() {
		detailsModalContentId = null;
	}

	// Mobile card-list breakpoint (< 860px) — replaces the AG Grid entirely, per design handoff.
	let cardMode = $state(false);
```

Note: `rowData` is declared later in the file (`const rowData = $derived(...)` inside the "Derived values from query" section) — Svelte's `$derived` closures resolve at call time, so declaring `detailsModalContent` above that point works the same way `popoverContentId`/`perspectivesByContentId` already do in this file (forward references to `$derived`/`$state` are fine within a single `<script>` block). Keep this block right after the existing `popoverExistingPerspective` declaration for consistency.

3c. Add the `cardMode` breakpoint effect — place it right after the existing `responsiveTier` effect (the one with `mqSm`/`mqMd`/`mqLg`, ending around `frontend/src/lib/components/ActivityTable.svelte:585`):

```ts
	// Card-mode breakpoint: below 860px, replace the grid with stacked cards entirely.
	$effect(() => {
		if (typeof window === 'undefined') return;
		const mq = window.matchMedia('(max-width: 859px)');
		const update = () => {
			cardMode = mq.matches;
		};
		update();
		mq.addEventListener('change', update);
		return () => mq.removeEventListener('change', update);
	});
```

3d. Update the `item` column def (around `frontend/src/lib/components/ActivityTable.svelte:277-286`) from:

```ts
			{
				colId: 'item',
				headerName: 'Item',
				flex: 2,
				minWidth: 200,

				filter: 'agTextColumnFilter',
				filterValueGetter: (params) => params.data?.name ?? '',
				cellRenderer: itemCellRenderer,
				tooltipValueGetter: (params) => params.data?.name ?? '',
				headerTooltip: 'Video title and thumbnail from YouTube API',
			},
```

to:

```ts
			{
				colId: 'item',
				headerName: 'Item',
				width: 200,
				minWidth: 200,
				maxWidth: 200,
				flex: 0,

				filter: 'agTextColumnFilter',
				filterValueGetter: (params) => params.data?.name ?? '',
				cellRenderer: activityItemCellRenderer,
				cellStyle: { padding: 0 },
				tooltipValueGetter: (params) => params.data?.name ?? '',
				headerTooltip: 'Video title and thumbnail from YouTube API',
			},
```

(Leave the `itemCellRenderer` import in place — it's still used by `nameCellRenderer`'s neighbors and the browser-test fixture; do not delete it from `formatting.ts` or its import here even if this file no longer calls it directly. If your editor flags the now-unused `itemCellRenderer` import in `ActivityTable.svelte` specifically, remove *only* that one import line, not the function itself.)

3e. Update `rowHeight` in the `theme` (around `frontend/src/lib/components/ActivityTable.svelte:254`):

```ts
		rowHeight: 112,
```

3f. Update the `gridOptions.context` initial value (around `frontend/src/lib/components/ActivityTable.svelte:337`) from:

```ts
		context: { perspectivesByContentId: new Map() },
```

to:

```ts
		context: { perspectivesByContentId: new Map(), onOpenDetails: handleOpenDetails },
```

3g. Update the reactive context-sync `$effect` (the one that currently does `gridApi.setGridOption('context', { perspectivesByContentId })`) to preserve `onOpenDetails`:

```ts
	$effect(() => {
		if (gridApi) {
			gridApi.setGridOption('context', { perspectivesByContentId, onOpenDetails: handleOpenDetails });
			gridApi.refreshCells({ columns: ['perspectize'], force: true });
		}
	});
```

3h. Update the row-height inline style and wrap the grid in `cardMode` conditional. Find this block in the template:

```svelte
		<!-- AG Grid -->
		<div
			bind:this={gridContainer}
			class="{isMobile ? 'overflow-y-auto' : 'flex-1'} min-h-0"
			style="--ag-row-height: 44px; --ag-header-height: 40px;"
		>
			<AgGridSvelte5Component {gridOptions} {rowData} {theme} {modules} />
		</div>
```

Replace it with:

```svelte
		{#if cardMode}
			<div class="flex-1 min-h-0 overflow-y-auto">
				<ActivityCardList {rowData} onOpenDetails={handleOpenDetails} />
			</div>
		{:else}
			<!-- AG Grid -->
			<div
				bind:this={gridContainer}
				data-testid="ag-grid-container"
				class="{isMobile ? 'overflow-y-auto' : 'flex-1'} min-h-0"
				style="--ag-row-height: 112px; --ag-header-height: 40px;"
			>
				<AgGridSvelte5Component {gridOptions} {rowData} {theme} {modules} />
			</div>
		{/if}
```

3i. Add the details modal at the bottom of the template, alongside the existing `PerspectivePopover` block:

```svelte
<ActivityDetailsModal
	content={detailsModalContent}
	open={detailsModalContentId !== null}
	onClose={handleCloseDetails}
/>
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd frontend && pnpm exec vitest run --project unit tests/components/ActivityTable.test.ts`
Expected: PASS (all tests, including the 2 new ones)

- [ ] **Step 5: Run the full unit test suite and build**

Run:
```bash
cd frontend && pnpm run test:run
cd frontend && pnpm run build
```
Expected: all tests pass, build succeeds with zero errors. Fix any TypeScript errors surfaced by the `ColDef` changes (e.g. AG Grid's `ColDef` type expects `width`/`minWidth`/`maxWidth`/`flex` as numbers — all values above already are).

- [ ] **Step 6: Commit**

```bash
git add frontend/src/lib/components/ActivityTable.svelte frontend/tests/components/ActivityTable.test.ts
git commit -m "feat: wire redesigned Item cell, details modal, and mobile cards into ActivityTable"
```

---

### Task 5: Visual verification with a browser agent

Not a code task — run the app and drive it with a browser automation agent to confirm look/feel and interactivity match the design handoff screenshots before opening the PR.

- [ ] **Step 1:** Start the frontend dev server (`pnpm --dir frontend run dev`, or use the `run` skill).
- [ ] **Step 2:** Navigate to the Activity page. If auth is required to see data, pause and ask the user to sign in.
- [ ] **Step 3:** At desktop width (>= 860px): verify row height ~112px, thumbnail hover shows the dark overlay + "Watch", clicking the thumbnail opens the video in a new tab, clicking elsewhere in the cell opens the details modal, and the modal's backdrop-click / × both close it.
- [ ] **Step 4:** Resize the viewport to < 860px: verify the table is replaced by the stacked card list (thumbnail + title + channel + duration), and that other columns/nav/header are otherwise unaffected.
- [ ] **Step 5:** Compare against `screenshot-desktop-table.png`, `screenshot-details-modal.png`, and `screenshot-mobile-cards.png` in `/Users/jamesjordan/Downloads/design_handoff_activity_item_column/` for fidelity (adjusted for real theme tokens vs. the mock's placeholder colors/fonts).
- [ ] **Step 6:** Report findings; fix any visual/interaction gaps found before proceeding to PR.

---

## Self-Review Notes

- **Spec coverage:** thumbnail size/radius/object-fit (Task 1), title gap/clamp/size/serif/max-width (Task 1), row height 112px (Task 4.3e/3h), two click zones + `stopPropagation` (Task 1, Task 3), thumbnail hover overlay + Watch label (Task 1), title hover underline (Task 1), row hover tint (pre-existing `rowHoverColor` theme param, untouched), modal layout/fields/backdrop-close/non-functional button (Task 2), mobile card layout + 860px breakpoint (Task 3, Task 4.3c/3h), design tokens pulled from `app.css` not the mock's literals (Tasks 1-3 use `--color-primary`/`border`/`accent`/`muted-foreground`/`font-family-serif` via Tailwind theme classes and one `font-[family-name:...]` arbitrary value), scope guard on other columns/nav/header (no edits to those files anywhere in this plan).
- **Deferred by user decision:** Perspectives count / Avg. Rating render as `—` placeholders (Task 2) — no backend task in this plan; tracked via GitHub issue #286.
