<script lang="ts">
	// Optional props for customization
	let { class: className = '', children } = $props();
</script>

<main class="px-4 py-6 md:px-6 md:py-8 lg:px-8 max-w-screen-xl mx-auto {className}">
	{@render children()}
</main>
