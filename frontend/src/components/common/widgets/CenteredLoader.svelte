<script lang="ts">
    import Loader from "./Loader.svelte";
    export let node: HTMLElement | undefined = undefined;
</script>

<div class="loader-container" bind:this={node}>
    <Loader/>
</div>

<style>
    .loader-container {
        display: flex;
        justify-content: center;
        width: 100%;
        margin-top: 1em;
        height: fit-content;
    }
</style>