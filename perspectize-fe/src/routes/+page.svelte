<script lang="ts">
	import { Button } from '$lib/components/ui/button';
</script>

<main class="p-8 max-w-4xl mx-auto">
	<h1 class="text-4xl font-bold mb-4 text-foreground">Perspectize</h1>
	<p class="mb-6 text-lg text-muted-foreground">
		This text should be in Inter font. Check DevTools to verify font-family.
	</p>

	<div class="space-x-4 mb-8">
		<Button>Primary (Navy)</Button>
		<Button variant="secondary">Secondary</Button>
		<Button variant="outline">Outline</Button>
		<Button variant="ghost">Ghost</Button>
	</div>

	<div class="mt-8 p-6 bg-primary text-primary-foreground rounded-lg shadow-md">
		<h2 class="text-2xl font-semibold mb-2">Theme Verification</h2>
		<p>This box should have a navy background (#1a365d) with white text.</p>
		<p class="mt-2 text-sm opacity-90">
			The primary color is set to: oklch(0.216 0.006 56.043)
		</p>
	</div>

	<div class="mt-8 grid grid-cols-1 md:grid-cols-3 gap-4">
		<div class="p-4 bg-secondary text-secondary-foreground rounded-lg">
			<h3 class="font-semibold mb-2">Secondary</h3>
			<p class="text-sm">Secondary color scheme</p>
		</div>
		<div class="p-4 bg-accent text-accent-foreground rounded-lg">
			<h3 class="font-semibold mb-2">Accent</h3>
			<p class="text-sm">Accent color scheme</p>
		</div>
		<div class="p-4 bg-muted text-muted-foreground rounded-lg">
			<h3 class="font-semibold mb-2">Muted</h3>
			<p class="text-sm">Muted color scheme</p>
		</div>
	</div>

	<div class="mt-8 space-y-2">
		<h2 class="text-2xl font-semibold mb-4">Tailwind Utilities Test</h2>
		<p class="text-sm text-muted-foreground">Testing various Tailwind classes:</p>
		<div class="flex gap-2">
			<div class="w-12 h-12 bg-red-500 rounded"></div>
			<div class="w-12 h-12 bg-blue-500 rounded"></div>
			<div class="w-12 h-12 bg-green-500 rounded"></div>
			<div class="w-12 h-12 bg-yellow-500 rounded"></div>
		</div>
	</div>
</main>
