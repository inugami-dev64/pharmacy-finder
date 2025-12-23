<script lang="ts">
    import Stars from "../../common/widgets/Stars.svelte";
    import Estrogen from "$lib/assets/regimen/estrogen.webp";
    import Testosterone from "$lib/assets/regimen/testosterone.webp";
    import { onMount } from "svelte";
    import EditButton from "../../common/icons/EditButton.svelte";
    import Countries from "$lib/assets/countries.json"
    import DeleteButton from "../../common/icons/DeleteButton.svelte";

    export let rating: number;
    export let prescriptionType: string;
    export let review: string | null;
    export let regimen: string | null;
    export let countryCode: string | null;
    export let onEdit: () => void;
    export let onDelete: () => void;

    let buttonName: string = "Show more...";
    let truncateText: boolean = true;
    let textElement: HTMLElement;
    let truncatable: boolean = false;

    onMount(async () => {
        truncatable = textElement.offsetHeight < textElement.scrollHeight || textElement.offsetWidth < textElement.scrollWidth;
    })

    function toggleReadMore(_: Event) {
        truncateText = !truncateText;
        if (truncateText) {
            buttonName = "Show more..."
        } else {
            buttonName = "Show less"
        }
    }
</script>

<div class="phr-review">
    <div class="review-header">
        <span>
            <Stars value={rating} scale={0.75}/>
            <EditButton size={24} on:click={_ => onEdit()}/>
            <DeleteButton size={24} on:close={_ => onDelete()}/>
        </span>
        <span style="display: flex; align-items: center">
            {#if countryCode !== null && countryCode in Countries}
            <span style="font-size: 20px" title={Countries[countryCode].name}>{Countries[countryCode].emoji}</span>
            {/if}
            {#if regimen == 'e'}
            <img src="{Estrogen}" alt="e" title="Estrogen based prescription">
            {:else if regimen == 't'}
            <img src="{Testosterone}" alt="t" title="Testosterone based prescription">
            {/if}
        </span>
    </div>
    <b>{prescriptionType}</b>
    <div id="content">
        <p class="{truncateText ? "hide-content" : ""}" bind:this={textElement}>{review}</p>
    </div>

    {#if truncatable}
    <button on:click|preventDefault={toggleReadMore} aria-label="Toggle Read More" id="toggleBtn">{buttonName}</button>
    {/if}
</div>

<style>
    .phr-review {
        margin-right: 15px;
        padding-bottom: 0.25em;
        border-bottom: 1px solid #cfcfcf;
        white-space: pre;
        text-wrap: wrap;
    }

    .phr-review > div  {
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .review-header {
        user-select: none;
    }

    img {
        width: 24px;
        height: 24px;
        user-select: none;
    }

    .hide-content {
        overflow: hidden;
        -webkit-line-clamp: 5;
        line-clamp: 5;
        -webkit-box-orient: vertical;
        box-orient: vertical;
        display: -webkit-box;
    }
</style>