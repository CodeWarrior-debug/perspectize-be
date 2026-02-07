import { render, type RenderResult } from '@testing-library/svelte';
import type { Component } from 'svelte';

export function renderComponent<T extends Record<string, any>>(
	component: Component<T>,
	props?: Partial<T>
): RenderResult<Component<T>> {
	// @ts-expect-error - Testing Library type mismatch with Svelte 5
	return render(component, props ? { props } : {});
}

export function expectClasses(element: HTMLElement, ...classes: string[]) {
	for (const cls of classes) {
		if (!element.classList.contains(cls)) {
			throw new Error(`Expected element to have class "${cls}", but it was not found`);
		}
	}
}
