<script lang="ts">
	import favicon from '$lib/assets/favicon.ico';

	let { children } = $props();

	import '../lib/i18n'
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

{@render children?.()}
