<script lang="ts">
	import type { Snippet } from 'svelte';
	import { Popover, PopoverContent, PopoverTrigger, buttonVariants, Button } from '$lib/components/shadcn';

	let {
		triggerLabel,
		triggerIcon,
		title,
		description,
		submitLabel,
		pendingLabel,
		isPending = false,
		isSubmitDisabled = false,
		formFields,
		onSubmit,
		open = $bindable(false),
		triggerVariant = 'default',
		triggerSize = 'default',
		align = 'end',
	}: {
		triggerLabel: string;
		triggerIcon: Snippet;
		title: string;
		description: string;
		submitLabel: string;
		pendingLabel: string;
		isPending?: boolean;
		isSubmitDisabled?: boolean;
		formFields: Snippet;
		onSubmit: () => void;
		open?: boolean;
		triggerVariant?: 'default' | 'outline' | 'ghost';
		triggerSize?: 'default' | 'sm' | 'icon';
		align?: 'start' | 'center' | 'end';
	} = $props();

	function handleSubmit(e: Event) {
		e.preventDefault();
		onSubmit();
	}
</script>

<Popover bind:open>
	<PopoverTrigger class={buttonVariants({ variant: triggerVariant, size: triggerSize })}>
		{@render triggerIcon()}
		{triggerLabel}
	</PopoverTrigger>
	<PopoverContent {align} sideOffset={8}>
		<form onsubmit={handleSubmit}>
			<div class="space-y-4">
				<div>
					<h3 class="font-semibold text-base">{title}</h3>
					<p class="text-muted-foreground text-sm mt-1">
						{description}
					</p>
				</div>

				{@render formFields()}

				<div class="flex gap-2 justify-end">
					<Button
						type="button"
						variant="outline"
						size="sm"
						onclick={() => (open = false)}
						disabled={isPending}
					>
						Cancel
					</Button>
					<Button type="submit" size="sm" disabled={isPending || isSubmitDisabled}>
						{isPending ? pendingLabel : submitLabel}
					</Button>
				</div>
			</div>
		</form>
	</PopoverContent>
</Popover>
