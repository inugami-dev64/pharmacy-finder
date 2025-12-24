<script lang="ts">
    export let onclick: ((e: MouseEvent) => void) | undefined = undefined;
</script>

<button class="danger-btn" onclick={onclick}>
    <slot/>
</button>

<style>
    .danger-btn {
        border: none;
        border-radius: 0.5em;
        min-width: 200px;
        box-shadow: none;
        background-color: #dc3545;
        color: white;
        box-sizing: border-box;
        padding: 1em;
        transition: all 0.5s ease-in-out;
    }

    .danger-btn:hover {
        cursor: pointer;
        background-color: #bb2d3b;
    }
</style>