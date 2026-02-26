<script lang="ts">
    export let title: string | null = null;
    export let hoverColor: string = "#10a7ff";
    export let href: string | undefined = undefined;
    export let blankHref: boolean | undefined = false;
    export let padding: number | undefined = 0.5;
    export let dummy: boolean | undefined = false;
</script>
{#if dummy}
<div style="--hover-color: {hoverColor}; --padding: {padding}em">
    <slot></slot>
</div>
{:else if href == null}
<button style="--hover-color: {hoverColor}; --padding: {padding}em" type="button" on:click aria-label="icon-btn" title={title}>
    <slot></slot>
</button>
{:else}
<a href="{href}" target={blankHref ? "_blank" : ""} style="--hover-color: {hoverColor}; --padding: {padding}em" type="button" on:click aria-label="icon-btn" title={title}>
    <slot></slot>
</a>
{/if}

<style>
    button, a {
        border: none;
        background-color: rgba(0, 0, 0, 0);
        padding: var(--padding);
        display: inline-block;
        vertical-align: middle;
    }

    :global(.svg-btn) {
        display: inline-block;
        vertical-align: middle;
        transition: fill 0.5s ease;

        :global(&:hover) {
            fill: var(--hover-color);
            cursor: pointer;
        }
    }
</style>