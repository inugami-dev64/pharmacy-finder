<script lang="ts">
    import type { PharmacyInfo } from "$lib/service/pharmacy-info";
    import CloseButton from "../common/icons/CloseButton.svelte";
    import StarRating from "../common/widgets/StarRating.svelte";
    import {_} from "svelte-i18n";
    import Loader from "../common/widgets/Loader.svelte";
    import Review from "./PharmacyView/Review.svelte";
    import { onDestroy } from "svelte";
    import { PAGER_LIMIT, PharmacyReview } from "$lib/service/pharmacy-review";
    import { ratingData, reviewData } from "$lib/stores";
    import IntersectionObserver from "svelte-intersection-observer";

    // Import logos
    import BenuLogo from "$lib/assets/benu-logo.svg"
    import ApothekaLogo from "$lib/assets/apotheka-logo.svg"
    import SudameapteekLogo from "$lib/assets/sudameapteek-logo.svg"
    import EuroapteekLogo from "$lib/assets/euroapteek-logo.svg"
    import ModifyReviewForm from "./PharmacyView/ModifyReviewForm.svelte";
    import TitleBar from "../common/TitleBar.svelte";

    // props
    let {
        pharmacy,
        onClose
    }: {pharmacy: PharmacyInfo, onClose: () => void} = $props();

    let showMoreAverageScores: boolean = $state(false);
    let showModifyReview: boolean = $state(false);

    let key: number | undefined = $state(undefined);
    let uniqueKey: number | undefined = $state(undefined);
    let fetchDone: boolean = $state(false);

    let reviews: PharmacyReview[] | undefined = $state(undefined);
    let overAllRating: number | undefined = $state(undefined);
    let eRating: number | undefined = $state(undefined);
    let tRating: number | undefined = $state(undefined);

    let element: HTMLElement | undefined = $state();

    const unsubscribeReviewData = reviewData.subscribe((initialReviews) => {
        reviews = initialReviews;
        fetchDone = false;
        if (reviews != null && reviews.length != 0) {
            key = reviews[reviews.length-1].updatedAt;
            uniqueKey = reviews[reviews.length-1].id;
        }

        if (reviews != null && reviews.length < PAGER_LIMIT)
            fetchDone = true;
    });

    const unsubscribeRatingData = ratingData.subscribe((ratings) => {
        if (ratings != null) {
            overAllRating = ratings.filter(v => v.hrtKind == null).at(0)?.stars || 0;
            eRating = ratings.filter(v => v.hrtKind == 'e').at(0)?.stars;
            tRating = ratings.filter(v => v.hrtKind == 't').at(0)?.stars;
        }
    });

    onDestroy(() => {
        unsubscribeReviewData();
        unsubscribeRatingData();
    });

    /**
     * Fetches next page in the review list and appends it to reviews list
     */
    async function updateReviewList() {
        if (pharmacy.id && reviews) {
            let newReviews = await PharmacyReview.readReviews(pharmacy.id, key, uniqueKey);
            if (newReviews.length != 0) {
                key = newReviews[newReviews.length-1].updatedAt;
                uniqueKey = newReviews[newReviews.length-1].id;
            }

            if (newReviews.length < PAGER_LIMIT)
                fetchDone = true;

            reviews.push(...newReviews);
        }
    }
</script>

<div class="phr-view">
    <TitleBar onClose={onClose}/>

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
        {#if overAllRating == null}
            <div class="loader-container">
                <Loader/>
            </div>
        {:else}
            <span><StarRating value={overAllRating} title="Overall rating"/></span>
            {#if (eRating || tRating) && !showMoreAverageScores}
                <button onclick={_ => showMoreAverageScores = !showMoreAverageScores}>{$_("map.sidebar.viewMoreRatings")}</button>
            {:else if (eRating || tRating)}
                {#if eRating}
                    <span><StarRating value={eRating || 0} title="E rating"/></span>
                {/if}
                {#if tRating}
                    <span><StarRating value={tRating || 0} title="T rating"/></span>
                {/if}
            <button onclick={_ => showMoreAverageScores = !showMoreAverageScores}>{$_("map.sidebar.viewLessRatings")}</button>
            {/if}
        {/if}
        <button onclick={_ => showModifyReview = true}>Write a review</button>
    </div>

    <!-- Container for pharmacy reviews -->
    <div class="phr-reviews">
    {#if reviews == null}
            <div class="loader-container">
                <Loader/>
            </div>
        {:else}
            {#if reviews.length == 0}
                <i>{$_("map.sidebar.noReviews")}</i>
            {:else}
                {#each reviews as review}
                <Review rating={review.stars || 0} countryCode={review.nationality || "EE"} prescriptionType={review.prescriptionType || ""} review={review.review || ""} regimen={review.hrtKind || ""}/>
                {/each}
                {#if !fetchDone}
                    <IntersectionObserver
                        {element}
                        on:observe={(e) => {
                            if (e.detail.isIntersecting)
                                updateReviewList().then();
                        }}
                    >
                        <div bind:this={element}>
                            Loading...
                        </div>
                    </IntersectionObserver>
                {/if}
            {/if}
        {/if}
    </div>
</div>

{#if showModifyReview}
    <ModifyReviewForm pharmacy={pharmacy} onClose={async () => {
        showModifyReview = false;
        await updateReviewList();
    }}/>
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
</style>