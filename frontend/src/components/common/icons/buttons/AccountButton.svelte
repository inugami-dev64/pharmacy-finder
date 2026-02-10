<script lang="ts">
    import AccountIcon from "../AccountIcon.svelte";
    import IconButton from "./IconButton.svelte";

    export let size: number;
    export let title: string;
    export let color: string | undefined = "#999999"
    export let padding: number | undefined = 0.5;
    export let onclick: (() => void) | undefined = undefined;
</script>

<IconButton on:click={() => onclick?.()} title="{title}" padding={padding}>
    <AccountIcon size={size} globalClass="svg-btn" color={color}/>
</IconButton>