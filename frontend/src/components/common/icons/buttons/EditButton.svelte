<script lang="ts">
    import EditIcon from "../EditIcon.svelte";
    import IconButton from "./IconButton.svelte";

    export let size;
</script>

<IconButton on:click title="Edit the review">
    <EditIcon size={size} globalClass="svg-btn"/>
</IconButton>