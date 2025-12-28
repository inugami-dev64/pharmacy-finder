<script lang="ts">
    import SearchIcon from "../SearchIcon.svelte";
    import IconButton from "./IconButton.svelte";

    export let size: number;
    export let title: string = "Search";
</script>

<IconButton on:click title="{title}">
    <SearchIcon size={size} globalClass="svg-btn"/>
</IconButton>