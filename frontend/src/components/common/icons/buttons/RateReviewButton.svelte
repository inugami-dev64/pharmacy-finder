<script lang="ts">
    import RateReviewIcon from "../RateReviewIcon.svelte";
    import IconButton from "./IconButton.svelte";

    export let size: number;
    export let title: string = "Rate this review";
</script>

<IconButton on:click title={title}>
    <RateReviewIcon size={size} globalClass="svg-btn"/>
</IconButton>