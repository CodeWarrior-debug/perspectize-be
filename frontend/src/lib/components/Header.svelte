<script lang="ts">
	import { Show, SignInButton, UserButton } from 'svelte-clerk';
	import AddVideoPopover from '$lib/components/AddVideoPopover.svelte';
</script>

<header class="h-16 border-b border-border bg-primary text-primary-foreground sticky top-0 z-50">
	<div class="h-full px-4 md:px-6 lg:px-8 max-w-screen-xl mx-auto flex items-center justify-between gap-2 md:gap-4">
		<a
			href="/"
			class="font-bold text-base sm:text-lg md:text-xl text-primary-foreground hover:text-primary-foreground/80 active:opacity-75 transition-colors min-w-0 truncate"
		>
			Perspectize
		</a>
		<div class="flex items-center gap-2 md:gap-4 shrink-0">
			<Show when="signed-in">
				<AddVideoPopover triggerVariant="outline" />
				<UserButton
					appearance={{
						elements: {
							avatarBox: 'w-8 h-8',
						},
					}}
				/>
			</Show>

			<Show when="signed-out">
				<SignInButton mode="modal">
					<button class="inline-flex items-center justify-center rounded-md text-sm font-medium h-9 px-4 py-2 border border-primary-foreground/20 text-primary-foreground hover:bg-primary-foreground/10 transition-colors">
						Sign In
					</button>
				</SignInButton>
			</Show>
		</div>
	</div>
</header>
