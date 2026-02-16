import type { ITooltipComp, ITooltipParams } from '@ag-grid-community/core';

export class TagsTooltip implements ITooltipComp {
	private el!: HTMLDivElement;

	init(params: ITooltipParams) {
		this.el = document.createElement('div');
		this.el.className = 'tags-tooltip';

		const tags: string[] = params.data?.tags ?? [];
		if (tags.length === 0) {
			this.el.textContent = 'No tags';
			return;
		}

		tags.forEach((tag) => {
			const chip = document.createElement('span');
			chip.className = 'tags-tooltip-chip';
			chip.textContent = tag;
			this.el.appendChild(chip);
		});
	}

	getGui() {
		return this.el;
	}
}
