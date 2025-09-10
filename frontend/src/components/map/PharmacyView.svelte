<script lang="ts">
    import type { PharmacyInfo } from "$lib/service/pharmacy-info";
    import CloseButton from "../common/icons/CloseButton.svelte";
    import StarRating from "../common/widgets/StarRating.svelte";

    // Import logos
    import BenuLogo from "$lib/assets/benu-logo.svg"
    import ApothekaLogo from "$lib/assets/apotheka-logo.svg"
    import SudameapteekLogo from "$lib/assets/sudameapteek-logo.svg"
    import EuroapteekLogo from "$lib/assets/euroapteek-logo.svg"
    import Loader from "../common/widgets/Loader.svelte";
    import Review from "./PharmacyView/Review.svelte";
    import { mount, onMount } from "svelte";
    import { getPharmacyReviews, type PharmacyReview } from "$lib/service/pharmacy-review";

    export let pharmacy: PharmacyInfo;
    export let visible: boolean;

    let showMoreAverageScores: boolean = false;
    let fetchDone: boolean = false;
    let reviews: PharmacyReview[] = []

    let key: number | null = null;
    let uniqueKey: number | null = null;

    let reviewContainer: HTMLElement;

    async function checkAndUpdateReviews(_: Event) {
        if (reviewContainer.scrollHeight <= reviewContainer.scrollTop + reviewContainer.clientHeight && pharmacy.id !== undefined) {
            reviews = await getPharmacyReviews(pharmacy.id, key, uniqueKey)
            if (reviews.length == 0)
                fetchDone = true
            else if (reviews[reviews.length - 1].id != uniqueKey && reviews[reviews.length - 1].updatedAt != key) {
                uniqueKey = reviews[reviews.length-1].updatedAt || 0;
                key = reviews[reviews.length - 1].id || 0;
            } else
                fetchDone = true;
        }
    }

</script>

{#if visible}
<div class="phr-view">
    <div class="close">
        <CloseButton size=32 on:click={(e) => visible = false}/>
    </div>

    {#if pharmacy.chain?.toLowerCase() == "benu"}
        <img alt="Benu logo" src="{BenuLogo}">
    {:else if pharmacy.chain?.toLowerCase() == "apotheka" }
        <img alt="Apotheka logo" src={ApothekaLogo}>
    {:else if pharmacy.chain?.toLowerCase() == "südameapteek"}
        <img alt="Südameapteek logo" src="{SudameapteekLogo}">
    {:else if pharmacy.chain?.toLowerCase() == "euroapteek"}
        <img alt="Euroapteek logo" src="{EuroapteekLogo}">
    {/if}

    <div class="phr-info">
        <div class="phr-primary">
            <p>{pharmacy.chain}</p>
            <h3>{pharmacy.name}</h3>
            <p>{pharmacy.address}, {pharmacy.city}, {pharmacy.county}, {pharmacy.postalCode}</p>
        </div>
        <!--<span><StarRating value={pharmacy.overallRating || 0} title="Overall rating"/></span>-->
        {#if !showMoreAverageScores}
        <button on:click|preventDefault={_ => showMoreAverageScores = !showMoreAverageScores}>View more...</button>
        {:else}
        <!--<span><StarRating value={pharmacy.acceptanceRating || 0} title="Acceptance rating"/></span>
        <span><StarRating value={pharmacy.eRating || 0} title="E rating"/></span>
        <span><StarRating value={pharmacy.tRating || 0} title="T rating"/></span>-->
        <button on:click|preventDefault={_ => showMoreAverageScores = !showMoreAverageScores}>View less</button>
        {/if}
    </div>

    <!-- Container for pharmacy reviews -->
    <div class="phr-reviews" on:scroll={checkAndUpdateReviews} bind:this={reviewContainer}>
        {#each reviews as review}
        <Review rating={review.stars || 0} countryCode={review.nationality || "EE"} prescriptionType={review.prescriptionType || ""} review={review.review || ""} regimen={review.hrtKind || ""}/>
        {/each}
        {#if !fetchDone}
        <div class="loader-container">
            <Loader/>
        </div>
        {/if}
    </div>
</div>
{/if}

<style>
    h3, p {
        padding: 0;
        margin: 0;
    }

    .loader-container {
        display: flex;
        justify-content: center;
        width: 100%;
        margin-top: 1em;
        height: fit-content;
    }

    .phr-view {
        display: flex;
        flex-direction: column;
        position: fixed;
        left: 1em;
        top: 1em;
        width: calc(25% - 2em);
        min-width: 380px;
        max-width: 420px;
        max-height: calc(100vh - 3em);
        padding: 1em;
        z-index: 1000;
        background-color: #ffffff;
        border: 1px solid #aaaaaa;
        border-radius: 1.25em;
        box-sizing: border-box;
    }

    .phr-view > img {
        margin-bottom: 2em;
        width: 100%;
    }

    .phr-info {
        margin-bottom: 0.5em;
        border-bottom: 1px solid #cfcfcf;
    }

    .phr-primary p, .phr-primary h3 {
        margin: 0.25em 0;
    }

    .phr-reviews {
        overflow: auto;
    }

    .close {
        display: flex;
        flex-direction: row;
        width: 100%;
        height: fit-content;
        justify-content: right;
        margin-bottom: 2em;
    }
</style>