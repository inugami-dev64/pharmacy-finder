<script lang="ts">
    import AccountIcon from "../AccountIcon.svelte";
    import IconButton from "./IconButton.svelte";

    export let size: number;
    export let title: string;
</script>

<IconButton on:click title="{title}">
    <AccountIcon size={size} globalClass="svg-btn"/>
</IconButton>