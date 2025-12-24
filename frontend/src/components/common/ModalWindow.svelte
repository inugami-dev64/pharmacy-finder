<script lang="ts">
    import TitleBar from "./TitleBar.svelte";

    export let zIndex: number;
    export let onClose: () => void;
    export let minWidth: number = 600;
    export let minHeight: number = 600;
</script>

<div class="modal-container" style="--zIndex: {zIndex}">
    <div class="modal-window" style="--minWidth: {minWidth}px; --minHeight: {minHeight}px">
        <TitleBar onClose={onClose}/>
        <slot/>
    </div>
</div>

<style>
    .modal-container {
        display: flex;
        justify-content: center;
        align-items: center;
        position: absolute;
        left: 0;
        top: 0;
        z-index: var(--zIndex);
        width: 100vw;
        height: 100vh;
        background-color: rgb(0 0 0 / 50%);
    }

    .modal-window {
        display: flex;
        flex-direction: column;
        justify-content: flex-start;
        background-color: #ffffff;
        width: 50%;
        height: 50%;
        border-radius: 1.25em;
        padding: 1em;
        border: 1px solid #aaaaaa;
        min-width: var(--minWidth);
        min-height: var(--minHeight);
    }
</style>