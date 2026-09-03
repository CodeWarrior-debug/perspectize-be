<script lang="ts">
	import {
		Dialog,
		DialogContent,
		DialogHeader,
		DialogTitle,
		DialogDescription,
		DialogFooter,
		Button,
	} from '$lib/components/shadcn';
	import { DATA_COLUMNS, INTERNAL_COLUMNS, type TogglableColumn } from '$lib/utils/grid-config';

	interface Props {
		/** Two-way bound open state. */
		open?: boolean;
		/** Whether to show the admin-only "Internal" group. */
		isAdmin?: boolean;
		/** colId → currently-visible, seeded from the live grid when the dialog opens. */
		visibility?: Record<string, boolean>;
		/** True once the user has taken manual control this session. */
		overrideActive?: boolean;
		/** Called with (colId, nextVisible) on each checkbox change. */
		onToggle: (colId: string, next: boolean) => void;
	}

	let { open = $bindable(false), isAdmin = false, visibility = {}, overrideActive = false, onToggle }: Props = $props();

	function handleChange(col: TogglableColumn, e: Event) {
		onToggle(col.colId, (e.currentTarget as HTMLInputElement).checked);
	}
</script>

<Dialog bind:open>
	<DialogContent class="max-w-lg">
		<DialogHeader>
			<DialogTitle>Columns</DialogTitle>
			<DialogDescription>Choose which columns appear in the table for this session.</DialogDescription>
		</DialogHeader>

		<div class="max-h-[60vh] space-y-6 overflow-y-auto py-2">
			<fieldset class="space-y-2">
				<legend class="mb-1 text-xs font-semibold uppercase tracking-wide text-muted-foreground">Columns</legend>
				{#each DATA_COLUMNS as col (col.colId)}
					<label class="flex items-center gap-3 rounded-md px-2 py-1.5 text-sm hover:bg-accent">
						<input
							type="checkbox"
							class="size-4 rounded border-input accent-primary"
							checked={visibility[col.colId] ?? false}
							onchange={(e) => handleChange(col, e)}
						/>
						<span>{col.label}</span>
					</label>
				{/each}
			</fieldset>

			{#if isAdmin}
				<fieldset class="space-y-2">
					<legend class="mb-1 text-xs font-semibold uppercase tracking-wide text-muted-foreground"> Internal </legend>
					{#each INTERNAL_COLUMNS as col (col.colId)}
						<label class="flex items-center gap-3 rounded-md px-2 py-1.5 text-sm hover:bg-accent">
							<input
								type="checkbox"
								class="size-4 rounded border-input accent-primary"
								checked={visibility[col.colId] ?? false}
								onchange={(e) => handleChange(col, e)}
							/>
							<span>{col.label}</span>
						</label>
					{/each}
				</fieldset>
			{/if}
		</div>

		<p
			class="rounded-md px-3 py-2 text-xs {overrideActive
				? 'bg-muted font-medium text-foreground'
				: 'text-muted-foreground'}"
			data-testid="session-hint"
		>
			{#if overrideActive}
				Columns are set manually for this session — refresh the page to return to the standard columns.
			{:else}
				Column choices apply for this session only. Refresh the page to return to the standard columns.
			{/if}
		</p>

		<DialogFooter>
			<Button type="button" onclick={() => (open = false)}>Done</Button>
		</DialogFooter>
	</DialogContent>
</Dialog>
