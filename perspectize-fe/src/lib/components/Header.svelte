<script lang="ts">
	import { toast } from 'svelte-sonner';
	import { Button } from '$lib/components/ui/button';

	function handleAddVideo() {
		// Placeholder - modal will be implemented in Phase 3
		toast.info('Add Video modal coming in Phase 3');
	}
</script>

<header class="h-16 border-b border-border bg-background sticky top-0 z-50">
	<div
		class="h-full px-4 md:px-6 lg:px-8 max-w-screen-xl mx-auto flex items-center justify-between"
	>
		<a href="/" class="font-bold text-xl text-foreground hover:text-primary transition-colors">
			Perspectize
		</a>
		<div class="flex items-center gap-4">
			<Button onclick={handleAddVideo}>Add Video</Button>
		</div>
	</div>
</header>
