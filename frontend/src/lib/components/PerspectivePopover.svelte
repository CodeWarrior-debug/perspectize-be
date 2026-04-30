<script lang="ts">
	import { toast } from 'svelte-sonner';
	import ExternalLinkIcon from '@lucide/svelte/icons/external-link';
	import InfoIcon from '@lucide/svelte/icons/info';
	import XIcon from '@lucide/svelte/icons/x';
	import { MediaQuery } from 'svelte/reactivity';
	import {
		Dialog,
		DialogContent,
		DialogHeader,
		DialogTitle,
		Drawer,
		DrawerContent,
		DrawerHeader,
		DrawerTitle,
		Button,
	} from '$lib/components/shadcn';
	import RatingInput from '$lib/components/RatingInput.svelte';
	import Thumbs from '$lib/components/Thumbs.svelte';
	import CommentEditor from '$lib/components/CommentEditor.svelte';
	import CommentFullscreen from '$lib/components/CommentFullscreen.svelte';
	import AddFieldSearch from '$lib/components/AddFieldSearch.svelte';
	import type { FieldDef } from '$lib/components/AddFieldSearch.svelte';
	import { useCreatePerspective } from '$lib/queries/hooks/useCreatePerspective';
	import { useUpdatePerspective } from '$lib/queries/hooks/useUpdatePerspective';
	import type { PerspectiveItem } from '$lib/queries/perspectives';

	/**
	 * PerspectivePopover — centered modal for creating or editing a perspective.
	 *
	 * Redesigned layout (Variation A):
	 * - Header with content title (single line, truncated)
	 * - Overall (thumbs) + Comment row directly under header
	 * - Scrollable body: 2x2 rating grid with dynamic field adding
	 * - Equal-width Save/Cancel buttons at bottom
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

	// --- Core rating fields (backend-mapped) ---
	let quality = $state<number | null>(null);
	let agreement = $state<number | null>(null);
	let importance = $state<number | null>(null);
	let confidence = $state<number | null>(null);

	// Like field — 'THUMBS_UP', 'THUMBS_DOWN', or null
	type LikeValue = 'THUMBS_UP' | 'THUMBS_DOWN' | null;
	let likeValue = $state<LikeValue>(null);

	// Comment (rich text HTML)
	// TODO: Backend integration — comment field not yet in GraphQL schema
	let comment = $state('');
	let commentFullscreenOpen = $state(false);

	// Dynamic fields — tracks which rating fields are shown
	const DEFAULT_FIELDS = ['quality', 'agreement', 'importance', 'confidence'];
	let activeFields = $state<string[]>([...DEFAULT_FIELDS]);

	// Dynamic field values for non-core fields
	// TODO: Backend integration — custom/suggested fields not yet in schema
	let dynamicValues = $state<Record<string, number | null>>({});

	// Mapping from field key to bindable state getter/setter
	function getFieldValue(key: string): number | null {
		switch (key) {
			case 'quality':
				return quality;
			case 'agreement':
				return agreement;
			case 'importance':
				return importance;
			case 'confidence':
				return confidence;
			default:
				return dynamicValues[key] ?? null;
		}
	}

	function setFieldValue(key: string, val: number | null) {
		switch (key) {
			case 'quality':
				quality = val;
				break;
			case 'agreement':
				agreement = val;
				break;
			case 'importance':
				importance = val;
				break;
			case 'confidence':
				confidence = val;
				break;
			default:
				dynamicValues = { ...dynamicValues, [key]: val };
				// TODO: Backend integration — persist custom field value
				console.log('[Perspectize] dynamic field changed', key, val);
				break;
		}
	}

	function getFieldLabel(key: string): string {
		const labels: Record<string, string> = {
			quality: 'Quality',
			agreement: 'Agreement',
			importance: 'Importance',
			confidence: 'Confidence',
		};
		return (
			labels[key] ??
			key
				.replace(/^custom:/, '')
				.replace(/-/g, ' ')
				.replace(/\b\w/g, (c) => c.toUpperCase())
		);
	}

	function addField(field: FieldDef) {
		if (!activeFields.includes(field.key)) {
			activeFields = [...activeFields, field.key];
		}
	}

	function removeField(key: string) {
		activeFields = activeFields.filter((k) => k !== key);
		if (!DEFAULT_FIELDS.includes(key)) {
			const { [key]: _, ...rest } = dynamicValues;
			dynamicValues = rest;
		} else {
			setFieldValue(key, null);
		}
	}

	// Reset state when existingPerspective changes
	$effect(() => {
		quality = existingPerspective?.quality ?? null;
		agreement = existingPerspective?.agreement ?? null;
		importance = existingPerspective?.importance ?? null;
		confidence = existingPerspective?.confidence ?? null;
		const l = existingPerspective?.like;
		likeValue = l === 'THUMBS_UP' ? 'THUMBS_UP' : l === 'THUMBS_DOWN' ? 'THUMBS_DOWN' : null;
		comment = '';
		commentFullscreenOpen = false;
		activeFields = [...DEFAULT_FIELDS];
		dynamicValues = {};
	});

	const isMobile = new MediaQuery('(max-width: 639px)');

	const createMutation = useCreatePerspective();
	const updateMutation = useUpdatePerspective();
	const isPending = $derived(createMutation.isPending || updateMutation.isPending);

	const hasComment = $derived(!!comment.replace(/<[^>]*>/g, '').trim());

	function handleCommentChange(html: string) {
		comment = html;
		// TODO: Backend integration — comment not yet persisted
		console.log('[Perspectize] comment changed', html);
	}

	function handleSubmit(e: Event) {
		e.preventDefault();

		const hasAnyRating = quality !== null || agreement !== null || importance !== null || confidence !== null;
		const hasLike = likeValue !== null;
		const hasDynamic = Object.values(dynamicValues).some((v) => v !== null);

		if (!hasAnyRating && !hasLike && !hasComment && !hasDynamic) {
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
</script>

{#snippet modalBody(mobile: boolean)}
	<!-- Header -->
	<div class="px-5 pt-4 pb-3 border-b border-border relative">
		{#if mobile}
			<!-- Grabber handle -->
			<div class="flex justify-center pb-2">
				<div class="w-9 h-1 rounded-full bg-black/[0.18]"></div>
			</div>
		{/if}
		<div class="flex items-center justify-center gap-2">
			<h2 class="text-lg font-semibold text-center tracking-tight">
				{isEditMode ? 'Edit perspective' : 'Add perspective'}
			</h2>
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
		<p class="text-[13px] font-serif text-muted-foreground text-center mt-0.5 line-clamp-1 max-w-full">
			{contentName}
		</p>
	</div>

	<!-- Overall + Comment row -->
	<div class="flex items-center gap-3.5 px-5 py-3.5 border-b border-border">
		<div class="flex flex-col items-center gap-1.5 flex-shrink-0">
			<span class="font-sans text-[11px] font-medium text-muted-foreground uppercase tracking-[0.06em]"> Overall </span>
			<Thumbs
				value={likeValue}
				onChange={(val) => {
					likeValue = val;
				}}
				size="md"
			/>
		</div>

		{#if mobile}
			<button
				type="button"
				onclick={() => {
					commentFullscreenOpen = true;
				}}
				aria-label="Add comment"
				class="flex flex-1 items-center justify-between gap-2 h-14 px-3 rounded-lg border border-border bg-white cursor-pointer text-left"
			>
				<span
					class="font-serif text-[13px] line-clamp-2 flex-1"
					class:text-foreground={hasComment}
					class:text-muted-foreground={!hasComment}
				>
					{#if hasComment}
						{comment.replace(/<[^>]*>/g, '')}
					{:else}
						Add a comment
					{/if}
				</span>
				<span class="text-muted-foreground flex-shrink-0">
					<ExternalLinkIcon class="size-4" />
				</span>
			</button>
		{:else}
			<div class="flex-1 min-w-0 relative">
				<CommentEditor
					value={comment}
					onChange={handleCommentChange}
					minHeight={68}
					showPopout={true}
					onPopout={() => {
						commentFullscreenOpen = true;
					}}
				/>
			</div>
		{/if}
	</div>

	<!-- Scrollable body — ratings + add field -->
	<form onsubmit={handleSubmit} class="flex flex-col flex-1 overflow-hidden">
		<div class="flex-1 overflow-y-auto px-5 py-4 flex flex-col gap-3.5">
			<div class="grid grid-cols-2 gap-3">
				{#each activeFields as fieldKey, idx (fieldKey)}
					{@const isLeftCol = idx % 2 === 0}
					<div class="relative">
						{#if fieldKey === 'quality'}
							<RatingInput
								label="Quality"
								name="quality"
								bind:value={quality}
								compact
								trackWidth={110}
								onRemove={() => removeField('quality')}
							/>
						{:else if fieldKey === 'agreement'}
							<RatingInput
								label="Agreement"
								name="agreement"
								bind:value={agreement}
								compact
								trackWidth={110}
								onRemove={() => removeField('agreement')}
							/>
						{:else if fieldKey === 'importance'}
							<RatingInput
								label="Importance"
								name="importance"
								bind:value={importance}
								compact
								trackWidth={110}
								onRemove={() => removeField('importance')}
							/>
						{:else if fieldKey === 'confidence'}
							<RatingInput
								label="Confidence"
								name="confidence"
								bind:value={confidence}
								compact
								trackWidth={110}
								onRemove={() => removeField('confidence')}
							/>
						{:else}
							<RatingInput
								label={getFieldLabel(fieldKey)}
								name={fieldKey}
								value={dynamicValues[fieldKey] ?? null}
								compact
								trackWidth={110}
								onRemove={() => removeField(fieldKey)}
							/>
						{/if}
						<button
							type="button"
							onclick={() => removeField(fieldKey)}
							class="absolute top-1/2 -translate-y-1/2 flex items-center justify-center bg-white border border-border rounded-full shadow-sm text-muted-foreground hover:opacity-70 transition-opacity z-10"
							style="width: 16px; height: 16px; {isLeftCol
								? `right: ${mobile ? '4px' : '12px'};`
								: `left: ${mobile ? '4px' : '12px'};`}"
							aria-label="Remove {getFieldLabel(fieldKey)}"
							title="Remove {getFieldLabel(fieldKey)}"
						>
							<XIcon class="size-2" strokeWidth={3} />
						</button>
					</div>
				{/each}
			</div>

			<AddFieldSearch addedKeys={activeFields} onAdd={addField} placeholder="Add a field — e.g. clarity" dense />
		</div>

		<!-- Action buttons — extra bottom padding on mobile for home indicator -->
		<div
			class="flex gap-2.5 px-5 border-t border-border bg-background"
			style="padding-top: 14px; padding-bottom: {mobile ? 'calc(22px + env(safe-area-inset-bottom))' : '14px'};"
		>
			<Button
				type="button"
				variant="outline"
				size="default"
				onclick={() => onClose()}
				disabled={isPending}
				class="flex-1"
			>
				Cancel
			</Button>
			<Button type="submit" size="default" disabled={isPending} class="flex-1">
				{#if isPending}
					{isEditMode ? 'Saving...' : 'Adding...'}
				{:else}
					Save perspective
				{/if}
			</Button>
		</div>
	</form>
{/snippet}

<!-- Desktop: centered dialog -->
{#if !isMobile.current}
	<Dialog
		bind:open
		onOpenChange={(isOpen) => {
			if (!isOpen) onClose();
		}}
	>
		<DialogContent
			class="sm:max-w-[460px] w-[calc(100vw-2rem)] max-h-[90vh] overflow-hidden p-0 flex flex-col"
			overlayClass="bg-black/45"
		>
			{@render modalBody(false)}
		</DialogContent>
	</Dialog>
{:else}
	<!-- Mobile: bottom sheet drawer -->
	<Drawer
		bind:open
		onOpenChange={(isOpen) => {
			if (!isOpen) onClose();
		}}
	>
		<DrawerContent class="max-h-[88vh] flex flex-col p-0">
			{@render modalBody(true)}
		</DrawerContent>
	</Drawer>
{/if}

{#if commentFullscreenOpen}
	<CommentFullscreen
		value={comment}
		onChange={handleCommentChange}
		onClose={() => {
			commentFullscreenOpen = false;
		}}
	/>
{/if}
