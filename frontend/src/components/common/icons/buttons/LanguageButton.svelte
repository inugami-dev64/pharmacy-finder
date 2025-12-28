<script lang="ts">
    import LanguageIcon from "../LanguageIcon.svelte";
    import IconButton from "./IconButton.svelte";

    export let size: number;
</script>

<IconButton on:click title="Select language">
    <LanguageIcon size={size} globalClass="svg-btn"/>
</IconButton>