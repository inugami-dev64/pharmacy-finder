<script lang="ts">
    import favicon from "$lib/assets/favicon.ico";

    let { children } = $props();

    import '$lib/i18n';
</script>

<svelte:head>
    <meta
        name="description"
        content="Moderator users - Pharmacy Finder">
    <meta
        property="og:description"
        content="Moderator users - Pharmacy Finder">
    <meta
        property="og:title"
        content="Pharmacy Finder">
    <meta
        property="og:image"
        content={favicon}>
    <link rel="icon" href={favicon}/>
</svelte:head>

{@render children?.()}