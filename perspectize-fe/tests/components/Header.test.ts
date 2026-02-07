import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/svelte';
import Header from '$lib/components/Header.svelte';
import * as sonner from 'svelte-sonner';

vi.mock('svelte-sonner', () => ({
	toast: {
		info: vi.fn()
	},
	Toaster: vi.fn()
}));

function renderHeader() {
	const result = render(Header);
	const header = result.container.querySelector('header');
	const inner = result.container.querySelector('header > div');
	return { ...result, header, inner };
}

describe('Header component', () => {
	it('renders without errors', () => {
		render(Header);
	});

	it('renders the Perspectize brand name', () => {
		render(Header);
		expect(screen.getByText('Perspectize')).toBeInTheDocument();
	});

	it('renders a header element', () => {
		const { header } = renderHeader();
		expect(header).toBeInTheDocument();
	});

	it('header has fixed height class (h-16)', () => {
		const { header } = renderHeader();
		expect(header?.className).toContain('h-16');
	});

	it('header has bottom border', () => {
		const { header } = renderHeader();
		expect(header?.className).toContain('border-b');
	});

	it('has responsive padding classes on inner container', () => {
		const { inner } = renderHeader();
		expect(inner?.className).toContain('px-4');
	});

	it('has max-width constraint for large screens', () => {
		const { inner } = renderHeader();
		expect(inner?.className).toContain('max-w-screen-xl');
	});

	it('renders Add Video button', () => {
		render(Header);
		expect(screen.getByRole('button', { name: /add video/i })).toBeInTheDocument();
	});

	it('calls toast.info when Add Video button is clicked', async () => {
		render(Header);
		const button = screen.getByRole('button', { name: /add video/i });
		await fireEvent.click(button);
		expect(sonner.toast.info).toHaveBeenCalledWith('Add Video modal coming in Phase 3');
	});
});
