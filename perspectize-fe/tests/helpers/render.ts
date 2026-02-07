import { render, type RenderResult } from '@testing-library/svelte';
import type { Component } from 'svelte';

export function renderComponent<T extends Record<string, any>>(
	component: Component<T>,
	props?: Partial<T>
): RenderResult<Component<T>> {
	return render(component, { props: props as T });
}

export function expectClasses(element: HTMLElement, ...classes: string[]) {
	for (const cls of classes) {
		expect(element.classList.contains(cls)).toBe(true);
	}
}
