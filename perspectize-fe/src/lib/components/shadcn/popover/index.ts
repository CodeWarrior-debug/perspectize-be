import { Popover as PopoverPrimitive } from "bits-ui";
import Root from "./popover.svelte";
import Trigger from "./popover-trigger.svelte";
import Content from "./popover-content.svelte";

const Popover = {
	Root,
	Trigger,
	Content,
};

export {
	Root,
	Trigger,
	Content,
	//
	Root as Popover,
	Trigger as PopoverTrigger,
	Content as PopoverContent,
};
