<script lang="ts">
    export let size;
</script>

<button type="button" on:click aria-label="Close">
    <svg xmlns="http://www.w3.org/2000/svg" height="{size}" width="{size}" viewBox="0 -960 960 960" fill="#999999">
        <path d="m256-200-56-56 224-224-224-224 56-56 224 224 224-224 56 56-224 224 224 224-56 56-224-224-224 224Z"/>
    </svg>
</button>

<style>
    button {
        border: none;
        background-color: rgba(0, 0, 0, 0);
        padding: 0.5em;
        display: inline-block;
        vertical-align: middle;
    }

    svg {
        display: inline-block;
        vertical-align: middle;
        transition: fill 0.5s ease;
    }

    svg:hover {
        fill: #10a7ff;
        cursor: pointer;
    }
</style>