<script lang="ts">
	import { toast } from 'svelte-sonner';
	import { Button } from '$lib/components/shadcn';
	import UserSelector from '$lib/components/UserSelector.svelte';

	function handleAddVideo() {
		// Placeholder - modal will be implemented in Phase 3
		toast.info('Add Video modal coming in Phase 3');
	}
</script>

<header class="h-16 border-b border-border bg-white sticky top-0 z-50">
	<div
		class="h-full px-4 md:px-6 lg:px-8 max-w-screen-xl mx-auto flex items-center justify-between gap-2 md:gap-4"
	>
		<a href="/" class="font-bold text-base sm:text-lg md:text-xl text-foreground hover:text-primary transition-colors min-w-0 truncate">
			Perspectize
		</a>
		<div class="flex items-center gap-2 md:gap-4 shrink-0">
			<UserSelector />
			<Button size="sm" onclick={handleAddVideo}>Add Video</Button>
		</div>
	</div>
</header>
