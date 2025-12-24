<script lang="ts">
    import Stars from "../../common/widgets/stars/Stars.svelte";
    import Estrogen from "$lib/assets/regimen/estrogen.webp";
    import Testosterone from "$lib/assets/regimen/testosterone.webp";
    import { onMount } from "svelte";
    import EditButton from "../../common/icons/EditButton.svelte";
    import Countries from "$lib/assets/countries.json"
    import DeleteButton from "../../common/icons/DeleteButton.svelte";
    import type { PharmacyReview } from "$lib/service/pharmacy-review";
    import ImagoLogo from "../../common/icons/logos/ImagoLogo.svelte";
    import GenderGPLogo from "../../common/icons/logos/GenderGPLogo.svelte";

    export let review: PharmacyReview;
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
            <Stars value={review.stars ?? 5} scale={0.75}/>
            <EditButton size={24} on:click={_ => onEdit()}/>
            <DeleteButton size={22} on:click={_ => onDelete()}/>
        </span>
        <span style="display: flex; align-items: center">
            {#if review.nationality != null && review.nationality in Countries}
            <span style="font-size: 20px" title={Countries[review.nationality as keyof typeof Countries].name}>{Countries[review.nationality as keyof typeof Countries].emoji}</span>
            {/if}
            {#if review.hrtKind === 'e'}
            <img src="{Estrogen}" alt="e" title="Estrogen based prescription">
            {:else if review.hrtKind === 't'}
            <img src="{Testosterone}" alt="t" title="Testosterone based prescription">
            {/if}
        </span>
    </div>
    <time>{new Date(review.updatedAt ?? 0).toLocaleDateString()} {new Date(review.updatedAt ?? 0).toLocaleTimeString()}</time>
    <!-- Prescription type -->
    <span style="text-align: center">
        {#if review.prescriptionType === "Imago"}
            <ImagoLogo size={24}/>
        {:else if review.prescriptionType === "GenderGP"}
            <GenderGPLogo size={22}/>
        {:else if review.prescriptionType === "National"}
            🇪🇪
        {/if}
        <b>{review.prescriptionType}</b>
    </span>
    <div id="content">
        <p class="{truncateText ? "hide-content" : ""}" bind:this={textElement}>{review.review}</p>
    </div>

    {#if truncatable}
    <button on:click|preventDefault={toggleReadMore} aria-label="Toggle Read More" id="toggleBtn">{buttonName}</button>
    {/if}
</div>

<style>
    time {
        display: block;
        margin-bottom: 0.25em;
    }

    .phr-review {
        margin-right: 15px;
        padding-bottom: 0.25em;
        border-bottom: 1px solid #cfcfcf;
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