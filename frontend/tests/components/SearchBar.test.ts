import { describe, it, expect } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/svelte';
import SearchBar from '$lib/components/discover/SearchBar.svelte';

// Note: SearchBar exposes `debouncedQuery` as a $bindable prop (per the task
// brief). Programmatic render() here can't observe two-way bound prop writes
// without a dedicated host wrapper component, so this file sticks to
// DOM-observable behavior. Formal dedicated SearchBar tests (including
// debounce-timing coverage via a host component) land in a later plan.
describe('SearchBar', () => {
	it('renders with the expected placeholder text', () => {
		render(SearchBar);
		expect(screen.getByPlaceholderText('Search Content Sources...')).toBeInTheDocument();
	});

	it('reflects typed input in the field', async () => {
		render(SearchBar, { props: { value: '' } });
		const input = screen.getByPlaceholderText('Search Content Sources...') as HTMLInputElement;

		await fireEvent.input(input, { target: { value: 'svelte' } });

		expect(input.value).toBe('svelte');
	});

	it('does not show a clear button when there is no input', () => {
		render(SearchBar, { props: { value: '' } });
		expect(screen.queryByRole('button', { name: 'Clear search' })).not.toBeInTheDocument();
	});

	it('shows a clear button once there is input', () => {
		render(SearchBar, { props: { value: 'svelte' } });
		expect(screen.getByRole('button', { name: 'Clear search' })).toBeInTheDocument();
	});

	it('clears the input value when the clear button is clicked', async () => {
		render(SearchBar, { props: { value: 'svelte' } });
		const input = screen.getByPlaceholderText('Search Content Sources...') as HTMLInputElement;

		await fireEvent.click(screen.getByRole('button', { name: 'Clear search' }));

		expect(input.value).toBe('');
		expect(screen.queryByRole('button', { name: 'Clear search' })).not.toBeInTheDocument();
	});

	it('autofocuses the input on mount', () => {
		render(SearchBar);
		const input = screen.getByPlaceholderText('Search Content Sources...');
		expect(document.activeElement).toBe(input);
	});
});
