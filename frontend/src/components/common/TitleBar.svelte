<script lang="ts">
    import CloseButton from "./icons/CloseButton.svelte";

    export let onClose: () => void;
</script>

<div class="titlebar">
    <CloseButton size=32 on:click={(e) => onClose()}/>
</div>

<style>
    .titlebar {
        display: flex;
        flex-direction: row;
        width: 100%;
        height: fit-content;
        justify-content: right;
    }
</style>