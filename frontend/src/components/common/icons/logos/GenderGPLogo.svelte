<script lang="ts">
    import GenderGPLogo from "$lib/assets/gendergp-logo.webp";
    export let size: number = 24;
</script>

<span title="GenderGP">
    <img alt="" src="{GenderGPLogo}" style="--size: {size}px"/>
</span>

<style>
    img {
        display: inline;
        width: var(--size);
        height: var(---size);
    }
</style>
