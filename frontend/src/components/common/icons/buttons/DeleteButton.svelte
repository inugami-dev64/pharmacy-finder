<script lang="ts">
    import DeleteIcon from "../DeleteIcon.svelte";
    import IconButton from "./IconButton.svelte";

    export let size;
</script>

<IconButton hoverColor="#fb3c4f" on:click title="Delete the review">
    <DeleteIcon size={size} globalClass="svg-btn"/>
</IconButton>