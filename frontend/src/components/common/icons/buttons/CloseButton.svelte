<script lang="ts">
    import CloseIcon from "../CloseIcon.svelte";
    import IconButton from "./IconButton.svelte";

    export let size;
</script>

<IconButton on:click>
    <CloseIcon size={size} globalClass="svg-btn"/>
</IconButton>