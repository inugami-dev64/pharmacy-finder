<script lang="ts">
    export let size: number = 64;
</script>

<div class="navbar" style="--size: {size}px">
    <slot/>
</div>

<style>
    .navbar {
        display: flex;
        flex-direction: row;
        align-items: center;
        min-width: var(--size);
        min-height: var(--size);
        background-color: #ffffff;
        border: 1px solid #aaaaaa;
        border-radius: 0.5em;
        box-sizing: border-box;
    }
</style>