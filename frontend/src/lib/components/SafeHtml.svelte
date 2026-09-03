<script lang="ts">
	import { sanitizeHtml } from '$lib/utils/sanitize';

	/**
	 * SafeHtml — renders user-generated HTML (e.g., Tiptap review content)
	 * after sanitizing with DOMPurify. Use this instead of raw {@html}.
	 */
	let {
		html,
		class: className = '',
	}: {
		html: string;
		class?: string;
	} = $props();

	const clean = $derived(sanitizeHtml(html));
</script>

<div class="safe-html {className}">
	{@html clean}
</div>

<style>
	.safe-html :global(ul),
	.safe-html :global(ol) {
		padding-left: 22px;
		margin: 4px 0;
	}
	.safe-html :global(li) {
		margin: 2px 0;
	}
	.safe-html :global(p) {
		margin: 0 0 4px;
	}
	.safe-html :global(p:last-child) {
		margin-bottom: 0;
	}
</style>
