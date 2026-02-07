import { describe, it, expect } from 'vitest';
import { render } from '@testing-library/svelte';
import PageWrapper from '$lib/components/PageWrapper.svelte';
import { createRawSnippet } from 'svelte';

function createChildrenSnippet(text: string = 'Test content') {
	return createRawSnippet(() => ({
		render: () => `<span>${text}</span>`
	}));
}

describe('PageWrapper component', () => {
	it('renders a main element with children', () => {
		const { container } = render(PageWrapper, {
			props: { children: createChildrenSnippet() }
		});
		const main = container.querySelector('main');
		expect(main).toBeInTheDocument();
	});

	it('renders the provided children content', () => {
		const { container } = render(PageWrapper, {
			props: { children: createChildrenSnippet('Hello World') }
		});
		expect(container.textContent).toContain('Hello World');
	});

	it('has responsive padding classes', () => {
		const { container } = render(PageWrapper, {
			props: { children: createChildrenSnippet() }
		});
		const main = container.querySelector('main');
		expect(main?.className).toContain('px-4');
	});

	it('has max-width constraint', () => {
		const { container } = render(PageWrapper, {
			props: { children: createChildrenSnippet() }
		});
		const main = container.querySelector('main');
		expect(main?.className).toContain('max-w-screen-xl');
	});

	it('renders without custom className (default empty string)', () => {
		const { container } = render(PageWrapper, {
			props: { children: createChildrenSnippet(), class: undefined }
		});
		const main = container.querySelector('main');
		// Should have base classes but no extra custom class
		expect(main?.className).toContain('px-4');
		expect(main?.className).not.toContain('custom-class');
	});

	it('applies custom className when provided', () => {
		const { container } = render(PageWrapper, {
			props: { class: 'custom-class', children: createChildrenSnippet() }
		});
		const main = container.querySelector('main');
		expect(main?.className).toContain('custom-class');
	});
});
