<script lang="ts">
	import { toast } from 'svelte-sonner';
	import { Button } from '$lib/components/shadcn';
	import PageWrapper from '$lib/components/PageWrapper.svelte';
	import AGGridTest from '$lib/components/AGGridTest.svelte';
</script>

<PageWrapper>
	<h1 class="text-3xl font-bold mb-4">Activity</h1>
	<p class="text-muted-foreground mb-8">
		Recently updated content will appear here. (Phase 2)
	</p>

	<!-- AG Grid Validation Section -->
	<section class="mb-12">
		<AGGridTest />
	</section>

	<!-- Toast and Theme Tests -->
	<section class="border-t border-border pt-8">
		<h2 class="text-xl font-semibold mb-4">Theme & Toast Tests</h2>

		<div class="space-x-4 mb-4">
			<Button>Primary (Navy)</Button>
			<Button variant="secondary">Secondary</Button>
			<Button variant="outline">Outline</Button>
		</div>

		<div class="space-x-4">
			<Button variant="outline" size="sm" onclick={() => toast.success('Success!')}>
				Success Toast
			</Button>
			<Button variant="outline" size="sm" onclick={() => toast.error('Error!')}>
				Error Toast
			</Button>
		</div>
	</section>
</PageWrapper>
