<script lang="ts">
    import type { Snippet } from "svelte";

    let {
        children = undefined
    }: { children?: Snippet } = $props();
</script>

<div class="mod-container">
    <div class="list-container">
        {@render children?.()}
    </div>
</div>

<style>
    .mod-container {
        width: 100%;
        height: calc(100% - 65px);
        overflow: hidden;
        box-sizing: border-box;
    }


    .list-container {
        margin: 30px auto;
        overflow: auto;
        max-height: calc(100% - 60px);
        max-width: 50%;
        min-width: 400px;
        box-sizing: border-box;
        padding: 0.5em;
        border-radius: 0.5em;
        background-color: #ffffff;
    }
</style>