<script lang="ts">
    export let title: string | null = null;
</script>

<button type="button" on:click aria-label="icon-btn" title={title}>
    <slot></slot>
</button>

<style>
    button {
        border: none;
        background-color: rgba(0, 0, 0, 0);
        padding: 0.5em;
        display: inline-block;
        vertical-align: middle;
    }

    :global(.svg-btn) {
        display: inline-block;
        vertical-align: middle;
        transition: fill 0.5s ease;
    }

    :global(.svg-btn:hover) {
        fill: #10a7ff;
        cursor: pointer;
    }
</style>