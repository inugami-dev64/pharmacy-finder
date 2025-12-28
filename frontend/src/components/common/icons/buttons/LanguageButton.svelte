<script lang="ts">
    import LanguageIcon from "../LanguageIcon.svelte";
    import IconButton from "./IconButton.svelte";

    export let size: number;
    export let title: string = "Select language";
</script>

<IconButton on:click title="{title}">
    <LanguageIcon size={size} globalClass="svg-btn"/>
</IconButton>