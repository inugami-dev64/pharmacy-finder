<script lang="ts">
    import SearchIcon from "../SearchIcon.svelte";
    import IconButton from "./IconButton.svelte";

    export let size: number;
</script>

<IconButton on:click title="Search">
    <SearchIcon size={size} globalClass="svg-btn"/>
</IconButton>