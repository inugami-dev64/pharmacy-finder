<script lang="ts">
    import IconButton from "./IconButton.svelte";

    export let size;
</script>

<IconButton on:click>
    <svg class="svg-btn" xmlns="http://www.w3.org/2000/svg" height="{size}" width="{size}" viewBox="0 -960 960 960" fill="#999999">
        <path d="m256-200-56-56 224-224-224-224 56-56 224 224 224-224 56 56-224 224 224 224-56 56-224-224-224 224Z"/>
    </svg>
</IconButton>