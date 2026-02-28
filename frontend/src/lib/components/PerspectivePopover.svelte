<script lang="ts">
	import { toast } from 'svelte-sonner';
	import ThumbsUpIcon from '@lucide/svelte/icons/thumbs-up';
	import ThumbsDownIcon from '@lucide/svelte/icons/thumbs-down';
	/* import PlusIcon from '@lucide/svelte/icons/plus'; */
	import InfoIcon from '@lucide/svelte/icons/info';
	import {
		Dialog,
		DialogContent,
		DialogHeader,
		DialogTitle,
		Button,
	} from '$lib/components/shadcn';
	import RatingInput from '$lib/components/RatingInput.svelte';
	import { useCreatePerspective } from '$lib/queries/hooks/useCreatePerspective';
	import { useUpdatePerspective } from '$lib/queries/hooks/useUpdatePerspective';
	/* import { useCreateClaim } from '$lib/queries/hooks/useCreateClaim'; */
	import type { PerspectiveItem } from '$lib/queries/perspectives';

	/**
	 * PerspectivePopover — centered modal for creating or editing a perspective.
	 *
	 * Per CONTEXT.md (Figma Make decisions):
	 * - Always rendered as a centered Dialog (not anchored to table cell)
	 * - No text inputs (no Review textarea, no Title input) in Phase 4
	 * - Like field = thumbs up/down toggle buttons (THUMBS_UP / THUMBS_DOWN)
	 * - Ratings: 2x2 grid with hasInteracted tracking, submit null if not interacted
	 * - "+ Add More..." expansion for claim creation (claim creation is Wave 3)
	 */
	let {
		contentId,
		contentName,
		existingPerspective = null,
		userId,
		open = $bindable(true),
		onClose,
	}: {
		contentId: number;
		contentName: string;
		existingPerspective?: PerspectiveItem | null;
		userId: number;
		open?: boolean;
		onClose: () => void;
	} = $props();

	const isEditMode = $derived(existingPerspective !== null);

	// Rating state — null means unset (user has not interacted).
	// Initialize from existingPerspective for edit mode; $effect resets when it changes.
	let quality = $state<number | null>(null);
	let agreement = $state<number | null>(null);
	let importance = $state<number | null>(null);
	let confidence = $state<number | null>(null);

	// Like field — 'THUMBS_UP', 'THUMBS_DOWN', or null
	type LikeValue = 'THUMBS_UP' | 'THUMBS_DOWN' | null;
	function parseLike(l: string | null | undefined): LikeValue {
		if (l === 'THUMBS_UP') return 'THUMBS_UP';
		if (l === 'THUMBS_DOWN') return 'THUMBS_DOWN';
		return null;
	}
	let likeValue = $state<LikeValue>(null);

	/* "+ Add More..." expansion for claim creation — TODO: Re-enable in a future phase */
	/* let showMore = $state(false); */
	/* let claimText = $state(''); */

	// Reset state when existingPerspective changes (e.g., when switching rows)
	$effect(() => {
		quality = existingPerspective?.quality ?? null;
		agreement = existingPerspective?.agreement ?? null;
		importance = existingPerspective?.importance ?? null;
		confidence = existingPerspective?.confidence ?? null;
		const l = existingPerspective?.like;
		likeValue = l === 'THUMBS_UP' ? 'THUMBS_UP' : l === 'THUMBS_DOWN' ? 'THUMBS_DOWN' : null;
		/* showMore = false; */
		/* claimText = ''; */
	});

	const createMutation = useCreatePerspective();
	const updateMutation = useUpdatePerspective();
	/* const createClaimMutation = useCreateClaim(); */

	const isPending = $derived(createMutation.isPending || updateMutation.isPending);
	/* const isClaimPending = $derived(createClaimMutation.isPending); */

	function toggleLike(val: 'THUMBS_UP' | 'THUMBS_DOWN') {
		likeValue = likeValue === val ? null : val;
	}

	function handleSubmit(e: Event) {
		e.preventDefault();

		// Validate: at least one field must be non-empty
		const hasAnyRating = quality !== null || agreement !== null || importance !== null || confidence !== null;
		const hasLike = likeValue !== null;

		if (!hasAnyRating && !hasLike) {
			toast.error('Please fill in at least one field');
			return;
		}

		if (isEditMode && existingPerspective) {
			updateMutation.mutate(
				{
					id: parseInt(existingPerspective.id, 10),
					quality: quality ?? undefined,
					agreement: agreement ?? undefined,
					importance: importance ?? undefined,
					confidence: confidence ?? undefined,
					like: likeValue ?? undefined,
				},
				{
					onSuccess: () => {
						onClose();
					},
				},
			);
		} else {
			createMutation.mutate(
				{
					userID: userId,
					contentID: contentId,
					quality: quality ?? undefined,
					agreement: agreement ?? undefined,
					importance: importance ?? undefined,
					confidence: confidence ?? undefined,
					like: likeValue ?? undefined,
				},
				{
					onSuccess: () => {
						onClose();
					},
				},
			);
		}
	}

	/* function handleCreateClaim() {
		const trimmed = claimText.trim();
		if (!trimmed) {
			toast.error('Please enter claim text');
			return;
		}

		createClaimMutation.mutate(
			{
				text: trimmed,
				userID: userId,
				parentContentID: contentId,
			},
			{
				onSuccess: () => {
					// Clear the claim textarea; popover stays open so user can continue their perspective
					claimText = '';
				},
			},
		);
	} */

	function handleCancel() {
		onClose();
	}
</script>

<Dialog
	bind:open
	onOpenChange={(isOpen) => {
		if (!isOpen) onClose();
	}}
>
	<DialogContent
		class="sm:max-w-sm w-[calc(100vw-2rem)] max-h-[90vh] overflow-y-auto overflow-x-hidden p-0"
		overlayClass="bg-transparent"
	>
		<!-- Header -->
		<div class="px-6 py-4 border-b border-border">
			<DialogHeader>
				<div class="flex items-center justify-center gap-2 mb-2">
					<DialogTitle class="text-xl font-semibold text-center">
						{isEditMode ? 'Edit Perspective' : 'Add Perspective'}
					</DialogTitle>
					<div class="relative group">
						<button
							type="button"
							class="p-1 rounded-full text-muted-foreground hover:opacity-70 transition-opacity"
							aria-label="About perspectives"
						>
							<InfoIcon class="size-4" />
						</button>
						<div
							class="absolute left-1/2 -translate-x-1/2 top-full mt-2 px-3 py-2 rounded-lg shadow-xl text-xs w-56 text-center z-[100] bg-foreground text-background opacity-0 group-hover:opacity-100 pointer-events-none transition-opacity"
						>
							Add as much or as little as you like
							<div class="absolute left-1/2 -translate-x-1/2 -top-1 w-2 h-2 rotate-45 bg-foreground"></div>
						</div>
					</div>
				</div>
				<p class="text-sm text-muted-foreground text-center text-wrap line-clamp-2 max-w-full">
					{contentName}
				</p>
			</DialogHeader>
		</div>

		<!-- Form -->
		<form onsubmit={handleSubmit} class="px-6 py-4 space-y-4">
			<!-- Ratings 2x2 grid -->
			<div class="px-2">
				<div class="grid grid-cols-2 gap-4">
					<RatingInput label="Quality" name="quality" bind:value={quality} />
					<RatingInput label="Agreement" name="agreement" bind:value={agreement} />
					<RatingInput label="Importance" name="importance" bind:value={importance} />
					<RatingInput label="Confidence" name="confidence" bind:value={confidence} />
				</div>
			</div>

			<div class="border-t border-border"></div>

			<!-- Like — thumbs up/down toggle -->
			<div class="flex flex-col items-center gap-2">
				<span class="text-xs font-medium text-muted-foreground">Like</span>
				<div class="flex items-center justify-center gap-4">
					<button
						type="button"
						onclick={() => toggleLike('THUMBS_UP')}
						class="p-3 rounded-lg transition-all border-2"
						class:bg-green-500={likeValue === 'THUMBS_UP'}
						class:text-white={likeValue === 'THUMBS_UP'}
						class:border-green-500={likeValue === 'THUMBS_UP'}
						class:bg-muted={likeValue !== 'THUMBS_UP'}
						class:text-muted-foreground={likeValue !== 'THUMBS_UP'}
						class:border-transparent={likeValue !== 'THUMBS_UP'}
						aria-label="Thumbs up"
						aria-pressed={likeValue === 'THUMBS_UP'}
					>
						<ThumbsUpIcon class="size-6" strokeWidth={2} />
					</button>
					<button
						type="button"
						onclick={() => toggleLike('THUMBS_DOWN')}
						class="p-3 rounded-lg transition-all border-2"
						class:bg-red-500={likeValue === 'THUMBS_DOWN'}
						class:text-white={likeValue === 'THUMBS_DOWN'}
						class:border-red-500={likeValue === 'THUMBS_DOWN'}
						class:bg-muted={likeValue !== 'THUMBS_DOWN'}
						class:text-muted-foreground={likeValue !== 'THUMBS_DOWN'}
						class:border-transparent={likeValue !== 'THUMBS_DOWN'}
						aria-label="Thumbs down"
						aria-pressed={likeValue === 'THUMBS_DOWN'}
					>
						<ThumbsDownIcon class="size-6" strokeWidth={2} />
					</button>
				</div>
			</div>

			<!-- TODO: Re-enable Add More / Claim creation in a future phase -->
			<!-- <div class="flex justify-center pt-1">
				<button
					type="button"
					onclick={() => (showMore = !showMore)}
					class="flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium hover:opacity-80 transition-all border border-border text-primary"
				>
					<PlusIcon class="size-4" />
					Add More...
				</button>
			</div>

			{#if showMore}
				<div class="rounded-lg border border-border bg-muted/30 p-4 space-y-3 animate-in fade-in-0 slide-in-from-top-2">
					<label for="claim-text" class="text-sm font-semibold text-foreground">Add a Claim</label>
					<textarea
						id="claim-text"
						bind:value={claimText}
						placeholder="e.g., @this ran 22.3 mph in the 1987 game"
						rows={2}
						disabled={isClaimPending}
						class="w-full resize-none rounded-md border border-input bg-background px-3 py-2 text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring disabled:opacity-50"
					></textarea>
					<p class="text-xs text-muted-foreground">
						Use <code class="font-mono bg-muted px-1 rounded">@this</code> to reference the current content
					</p>
					<Button
						type="button"
						variant="outline"
						size="sm"
						onclick={handleCreateClaim}
						disabled={isClaimPending || !claimText.trim()}
						class="w-full"
					>
						{isClaimPending ? 'Creating...' : 'Create Claim'}
					</Button>
				</div>
			{/if} -->

			<!-- Action buttons -->
			<div class="flex items-center justify-center gap-3 pt-1">
				<Button
					type="button"
					variant="outline"
					size="default"
					onclick={handleCancel}
					disabled={isPending}
					class="px-6"
				>
					Cancel
				</Button>
				<Button type="submit" size="default" disabled={isPending} class="px-6">
					{#if isPending}
						{isEditMode ? 'Saving...' : 'Adding...'}
					{:else}
						{isEditMode ? 'Save Changes' : 'Add Perspective'}
					{/if}
				</Button>
			</div>
		</form>
	</DialogContent>
</Dialog>
