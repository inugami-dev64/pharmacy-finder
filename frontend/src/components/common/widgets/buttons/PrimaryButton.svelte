<script lang="ts">
    export let onclick: ((e: MouseEvent) => void) | undefined = undefined;
</script>

<button class="primary-btn" onclick={onclick}>
    <slot></slot>
</button>

<style>
    .primary-btn {
        border: none;
        border-radius: 0.5em;
        min-width: 200px;
        box-shadow: none;
        background-color: #30bfff;
        color: white;
        box-sizing: border-box;
        padding: 1em;
        transition: all 0.5s ease-in-out;
    }

    .primary-btn:hover {
        cursor: pointer;
        background-color: #10a7ff;
    }
</style>