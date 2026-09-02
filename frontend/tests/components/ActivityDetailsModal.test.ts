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
