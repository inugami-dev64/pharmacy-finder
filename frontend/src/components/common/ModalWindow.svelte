<script lang="ts">
    import Overlay from "./Overlay.svelte";
import TitleBar from "./TitleBar.svelte";

    export let zIndex: number;
    export let onClose: () => void;
    export let minWidth: number = 360;
    export let minHeight: number = 500;
</script>

<Overlay zIndex={zIndex}>
    <div class="modal-window" style="--minWidth: {minWidth}px; --minHeight: {minHeight}px">
        <TitleBar onClose={onClose}/>
        <slot/>
    </div>
</Overlay>

<style>
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