<script lang="ts">
	import favicon from '$lib/assets/favicon.ico';
    import 'leaflet/dist/leaflet.css'

	let { children } = $props();

	import '$lib/i18n'
</script>

<svelte:head>
	<meta
		name="description"
		content="Pharmacy finder - a community driven review platform
for trans and non-binary people in Estonia">
	<meta
		property="og:description"
		content="A community driven pharmacy review platform
for trans and non-binary people in Estonia">
	<meta
		property="og:title"
		content="Pharmacy Finder">
	<meta
		property="og:image"
		content="{favicon}">
	<link rel="icon" href={favicon} />
</svelte:head>

{@render children?.()}
