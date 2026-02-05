<script lang="ts">
	// Navigation slot will be added in Plan 04
</script>

<header class="h-16 border-b border-border bg-background">
	<div
		class="h-full px-4 md:px-6 lg:px-8 max-w-screen-xl mx-auto flex items-center justify-between"
	>
		<div class="font-bold text-xl text-foreground">Perspectize</div>
		<nav class="flex items-center gap-4">
			<!-- Navigation links added in Plan 04 -->
			<slot name="nav" />
		</nav>
	</div>
</header>
