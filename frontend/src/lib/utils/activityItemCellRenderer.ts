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
		img.onerror = () => img.remove(); // thumbnail unavailable — fall back to the plain bg-muted block
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
		'mt-1 line-clamp-2 max-w-[184px] whitespace-normal text-center font-[family-name:var(--font-family-serif)] text-[13px] leading-[1.4] text-foreground decoration-primary/30 group-hover/cell:underline';
	title.textContent = name;

	cell.appendChild(thumbWrap);
	cell.appendChild(title);

	return cell;
}
