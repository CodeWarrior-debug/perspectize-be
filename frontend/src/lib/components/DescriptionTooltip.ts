import type { ITooltipComp, ITooltipParams } from '@ag-grid-community/core';

export class DescriptionTooltip implements ITooltipComp {
	private el!: HTMLDivElement;

	init(params: ITooltipParams) {
		this.el = document.createElement('div');
		this.el.className = 'description-tooltip';

		const description: string = params.data?.description ?? '';
		if (!description) {
			this.el.textContent = 'No description';
			return;
		}

		this.el.textContent = description;
	}

	getGui() {
		return this.el;
	}
}
