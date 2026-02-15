<script lang="ts">
    import { _ } from "svelte-i18n";
    import { authenticationSession } from "$lib/service/auth-session";
    import TitleBar from "../../components/common/TitleBar.svelte";
    import { ModeratorPharmacyReview } from "$lib/service/data/moderator-pharmacy-review";
    import CenteredLoader from "../../components/common/widgets/CenteredLoader.svelte";
    import Review from "../../components/map/PharmacyView/Review.svelte";
    import { PAGER_LIMIT } from "$lib/service/data/pager";
    import IntersectionObserver from "svelte-intersection-observer/IntersectionObserver.svelte";
    import RateReviewButton from "../../components/common/icons/buttons/RateReviewButton.svelte";
    import CheckCircleIcon from "../../components/common/icons/CheckCircleIcon.svelte";
    import TimedDeletionIcon from "../../components/common/icons/TimedDeletionIcon.svelte";
    import IntermediateCheckIcon from "../../components/common/icons/IntermediateCheckIcon.svelte";
    import BellNotificationIcon from "../../components/common/icons/BellNotificationIcon.svelte";
    import ModerationModal from "../../components/mod/ModerationModal.svelte";
    import { ReviewModerationResult } from "$lib/service/data/moderation";
    import { onMount } from "svelte";
    import { localeSwitcher } from "$lib/service/locale";
    import ModeratorHeader from "../../components/mod/ModeratorHeader.svelte";

    class ReviewFilter {
        key?: number;
        uniqueKey?: number;
        desc: boolean = true;
        showUnmoderated: boolean = true;
        showModerated: boolean = false;
    }

    let filter: ReviewFilter = $state(new ReviewFilter());
    let fetchDone: boolean = $state(false);
    let reviews: ModeratorPharmacyReview[] = $state([]);

    let element: HTMLElement | undefined = $state();
    let orderByValue: string | undefined = $state();
    let showModerationModal: boolean = $state(false);
    let pendingReview: ModeratorPharmacyReview | undefined = $state(undefined);

    onMount(() => localeSwitcher.setDefault())

    const clearPager = () => {
        reviews = [];
        filter.key = undefined;
        fetchDone = false;
        filter.uniqueKey = undefined;
    }

    async function updateReviewList() {
        let newReviews = await ModeratorPharmacyReview.getModeratorPharmacyReviews(filter.key, filter.uniqueKey, filter.desc, filter.showUnmoderated, filter.showModerated);
        if (newReviews.length != 0) {
            filter.key = newReviews[newReviews.length-1].updatedAt;
            filter.uniqueKey = newReviews[newReviews.length-1].id;
        }

        if (newReviews.length < PAGER_LIMIT)
            fetchDone = true;

        reviews.push(...newReviews)
    }
</script>

<svelte:head>
    <title>Moderation panel - Pharmacy Finder</title>
</svelte:head>

<ModeratorHeader
    bannerHref="/"
    bannerTitle={$_("mod.header.gotoHome")}
    isAdmin={authenticationSession.isAdmin(authenticationSession.getSessionToken() || "")}
/>

<div class="mod-container">
    <div class="comment-container">
        <TitleBar>
            <div class="filters">
                <label for="orderBy">{$_("mod.filters.orderBy.title")}:</label>
                <select
                    name="orderBy"
                    bind:value={orderByValue}
                    onchange={() => {
                        console.log(orderByValue)
                        if (orderByValue == "newest")
                            filter.desc = true;
                        else filter.desc = false;
                        clearPager();
                    }}
                >
                    <option selected={filter.desc} value="newest">{$_("mod.filters.orderBy.newestFirst")}</option>
                    <option selected={!filter.desc} value="oldest">{$_("mod.filters.orderBy.oldestFirst")}</option>
                </select>
                <label for="unmoderated">{$_("mod.filters.showUnmoderated")}:</label>
                <input
                    type="checkbox"
                    name="unmoderated"
                    bind:checked={filter.showUnmoderated}
                    onchange={(_) => clearPager()}
                >
                <label for="moderated">{$_("mod.filters.showModerated")}:</label>
                <input
                    type="checkbox"
                    name="moderated"
                    bind:checked={filter.showModerated}
                    onchange={(_) => clearPager()}
                >
            </div>
        </TitleBar>
        {#each reviews as review}
            <Review review={review}>
                {#if !review.markedForDeletion && review.commentReviewResult !== ReviewModerationResult.None}
                    <span title="{$_("mod.comments.approved")}" class="icon">
                        <CheckCircleIcon size={24}/>
                    </span>
                {:else if review.markedForDeletion}
                    <span title="{$_("mod.comments.markedForDeletion")}" class="icon">
                        <TimedDeletionIcon size={24}/>
                    </span>
                {:else if review.commentReviewResult === ReviewModerationResult.None}
                    <span title="{$_("mod.comments.notModerated")}" class="icon">
                        <IntermediateCheckIcon size={24}/>
                    </span>
                {/if}
                {#if review.commentReviewResult === ReviewModerationResult.PersonalAttack}
                    <span title="{$_("mod.comments.personalAttack")}" class="icon">
                        <BellNotificationIcon size={24}/>
                    </span>
                {:else if review.commentReviewResult === ReviewModerationResult.Offensive}
                    <span title="{$_("mod.comments.offensive")}" class="icon">
                        <BellNotificationIcon size={24}/>
                    </span>
                {:else if review.commentReviewResult === ReviewModerationResult.Other}
                    <span title="{$_("mod.comments.other")}" class="icon">
                        <BellNotificationIcon size={24}/>
                    </span>
                {/if}
                <RateReviewButton
                    size={24}
                    title="Rate the review"
                    on:click={() => {
                        pendingReview = review;
                        showModerationModal = true;
                    }}
                />
            </Review>
        {/each}
        {#if reviews.length === 0 && fetchDone}
            <div style="width: 100%; text-align: center; margin-top: 1em">
                <i>{$_("mod.noReviews")}</i>
            </div>
        {/if}
        {#if !fetchDone}
            <IntersectionObserver
                {element}
                on:intersect={(e: CustomEvent<IntersectionObserverEntry>) => {
                    updateReviewList().then();
                }}
            >
                <CenteredLoader bind:node={element}/>
            </IntersectionObserver>
        {/if}
    </div>
</div>

{#if showModerationModal && pendingReview != null}
    <ModerationModal
        onClose={() => {
            showModerationModal = false;
            clearPager();
        }}
        review={pendingReview}
    />
{/if}

<style>
    :global(body) {
        background-color: #fdcece;
    }

    .icon {
        display: inline-block;
        vertical-align: middle;
    }

    .mod-container {
        width: 100%;
        height: calc(100% - 65px);
        overflow: hidden;
        box-sizing: border-box;
    }

    .comment-container {
        margin: 30px auto;
        overflow: auto;
        max-height: calc(100% - 60px);
        max-width: 50%;
        min-width: 400px;
        box-sizing: border-box;
        padding: 0.5em;
        border-radius: 0.5em;
        background-color: #ffffff;
    }
</style>