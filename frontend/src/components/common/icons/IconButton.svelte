<script lang="ts">
    export let title: string | null = null;
    export let hoverColor: string = "#10a7ff";
    export let href: string | undefined = undefined;
</script>

{#if href == null}
<button style="--hover-color: {hoverColor}" type="button" on:click aria-label="icon-btn" title={title}>
    <slot></slot>
</button>
{:else}
<a href="{href}" target="_blank" style="--hover-color: {hoverColor}" type="button" on:click aria-label="icon-btn" title={title}>
    <slot></slot>
</a>
{/if}

<style>
    button, a {
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
        fill: var(--hover-color);
        cursor: pointer;
    }
</style>