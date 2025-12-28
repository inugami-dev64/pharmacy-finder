<script lang="ts">
    import { number } from "svelte-i18n";
    import IconButton from "./IconButton.svelte";

    export let size: number;
    export let title: string = "View source code";
    export let href: string = "https://github.com/inugami-dev64/pharmacy-finder"
</script>

<IconButton on:click title="{title}" href={href}>
    <svg class="svg-btn" xmlns="http://www.w3.org/2000/svg" height="{size}" viewBox="0 -960 960 960" width="{size}" fill="#999999">
        <path d="M320-240 80-480l240-240 57 57-184 184 183 183-56 56Zm320 0-57-57 184-184-183-183 56-56 240 240-240 240Z"/>
    </svg>
</IconButton>