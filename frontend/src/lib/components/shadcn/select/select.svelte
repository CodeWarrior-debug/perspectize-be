<script lang="ts">
	import { Select as SelectPrimitive } from 'bits-ui';
	import type { Snippet } from 'svelte';

	let {
		open = $bindable(false),
		value = $bindable(undefined),
		onValueChange,
		type = 'single',
		children,
		...restProps
	}: {
		open?: boolean;
		value?: string;
		onValueChange?: (value: string | undefined) => void;
		type?: 'single';
		children: Snippet;
		disabled?: boolean;
		required?: boolean;
		name?: string;
	} = $props();
</script>

<SelectPrimitive.Root bind:open bind:value {onValueChange} {type} {...restProps}>
	{@render children?.()}
</SelectPrimitive.Root>
