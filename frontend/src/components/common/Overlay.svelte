<script lang="ts">
    export let zIndex: number;
</script>

<div class="overlay" style="--zIndex: {zIndex}">
    <slot></slot>
</div>

<style>
    .overlay {
        display: flex;
        justify-content: center;
        align-items: center;
        position: absolute;
        left: 0;
        top: 0;
        z-index: var(--zIndex);
        width: 100%;
        height: 100%;
        background-color: rgb(0 0 0 / 50%);
    }
</style>