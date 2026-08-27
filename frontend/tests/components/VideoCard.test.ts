import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/svelte';
import VideoCard from '$lib/components/discover/VideoCard.svelte';
import type { VideoItem } from '$lib/services/youtubeApi';

const video: VideoItem = {
	id: 'abc123',
	title: 'A great video',
	channelTitle: 'Great Channel',
	publishedAt: '2024-06-15T12:00:00Z',
	description: 'A description of the great video',
	thumbnails: {
		medium: { url: 'https://img.example/medium.jpg', width: 320, height: 180 },
	},
};

describe('VideoCard', () => {
	it('renders thumbnail, title, channel, and description', () => {
		render(VideoCard, { props: { video, onAdd: vi.fn() } });

		expect(screen.getByText('A great video')).toBeInTheDocument();
		expect(screen.getByText('Great Channel')).toBeInTheDocument();
		expect(screen.getByText('A description of the great video')).toBeInTheDocument();
		const img = screen.getByAltText('A great video') as HTMLImageElement;
		expect(img.src).toBe('https://img.example/medium.jpg');
	});

	it('renders an Add to Library button when not in library', () => {
		render(VideoCard, { props: { video, onAdd: vi.fn() } });
		expect(screen.getByRole('button', { name: 'Add to Library' })).toBeInTheDocument();
	});

	it('calls onAdd with the video id when Add to Library is clicked', async () => {
		const onAdd = vi.fn();
		render(VideoCard, { props: { video, onAdd } });

		await fireEvent.click(screen.getByRole('button', { name: 'Add to Library' }));

		expect(onAdd).toHaveBeenCalledWith('abc123');
	});

	it('shows a disabled checkmark button when already in library', () => {
		render(VideoCard, { props: { video, isInLibrary: true, onAdd: vi.fn() } });

		const button = screen.getByRole('button', { name: /In Library/ });
		expect(button).toBeDisabled();
	});

	it('shows an Adding... pending state while the mutation is in flight', () => {
		render(VideoCard, { props: { video, isPending: true, onAdd: vi.fn() } });

		const button = screen.getByRole('button', { name: 'Adding...' });
		expect(button).toBeDisabled();
	});
});
