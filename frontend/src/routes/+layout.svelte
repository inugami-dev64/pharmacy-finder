<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';

	let { children } = $props();

	import '../lib/i18n'
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

{@render children?.()}
