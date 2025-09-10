<script lang="ts">
    import Stars from "./Stars.svelte";

    export let title: string | null = null;
    export let scale: number = 1;
    export let value: number;
</script>

{#if value !== 0}
<div class="rating">
    <Stars scale={scale} value={value}/>

    {#if title !== null}
    <span>{title}: </span>
    {/if}
    <span>{value}/5</span>
</div>
{/if}


<style>
    .rating {
        margin-top: 0.5em;
        display: flex;
        justify-content: space-between;
    }
</style>