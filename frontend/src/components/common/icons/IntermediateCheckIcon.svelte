<script lang="ts">
    export let globalClass: string | undefined = undefined;
    export let color: string | undefined = "#BEA804";
    export let size: number;
</script>

<svg class="{globalClass}" xmlns="http://www.w3.org/2000/svg" height="{size}" viewBox="0 -960 960 960" width="{size}" fill="{color}">
    <path d="M240-440v-80h480v80H240Z"/>
</svg>