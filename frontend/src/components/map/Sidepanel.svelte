<script lang="ts">
    export let zIndex: number;
    export let rightAlign: boolean = false;

    const style = rightAlign ? "right: 1em" : "left: 1em";
</script>

<div
    class="sidepanel"
    style="{style}; --zIndex: {zIndex}"
>
    <slot></slot>
</div>

<style>
    .sidepanel {
        display: flex;
        flex-direction: column;
        position: fixed;
        top: 64px;
        width: calc(25% - 2em);
        min-width: 360px;
        max-width: 420px;
        max-height: calc(100% - 2em - 64px);
        padding: 1em;
        z-index: var(--zIndex);
        background-color: #ffffff;
        border: 1px solid #aaaaaa;
        border-radius: 1.25em;
        box-sizing: border-box;
    }
</style>