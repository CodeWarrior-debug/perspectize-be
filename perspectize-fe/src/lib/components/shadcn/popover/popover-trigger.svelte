<script lang="ts">
	import { Popover as PopoverPrimitive } from "bits-ui";

	let { children, ...restProps }: PopoverPrimitive.TriggerProps = $props();
</script>

<PopoverPrimitive.Trigger {...restProps}>
	{@render children?.()}
</PopoverPrimitive.Trigger>
