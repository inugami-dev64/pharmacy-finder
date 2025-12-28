<script lang="ts">
    import ShinyStarIcon from "../ShinyStarIcon.svelte";
    import IconButton from "./IconButton.svelte";

    export let size: number;
    export let title: string = "Pharmacy recommender";
</script>

<IconButton on:click title="{title}">
    <ShinyStarIcon size={size} globalClass="svg-btn"/>
</IconButton>