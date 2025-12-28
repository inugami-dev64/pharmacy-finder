<script lang="ts">
    import SourceCodeIcon from "../SourceCodeIcon.svelte";
    import IconButton from "./IconButton.svelte";

    export let size: number;
    export let title: string = "View source code";
    export let href: string = "https://github.com/inugami-dev64/pharmacy-finder"
</script>

<IconButton on:click title="{title}" href={href}>
    <SourceCodeIcon size={size} globalClass="svg-btn"/>
</IconButton>