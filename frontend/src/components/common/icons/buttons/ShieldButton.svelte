<script lang="ts">
    import ShieldIcon from "../ShieldIcon.svelte";
    import IconButton from "./IconButton.svelte";

    export let size: number;
    export let title: string;
    export let color: string | undefined = "#999";
    export let padding: number | undefined = 0.5;
    export let onclick: (() => void) | undefined = undefined;
</script>

<IconButton on:click={() => onclick?.()} title="{title}" padding={padding}>
    <ShieldIcon size={size} globalClass="svg-btn" color={color}/>
</IconButton>