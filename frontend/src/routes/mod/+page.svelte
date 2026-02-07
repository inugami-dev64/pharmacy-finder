<script lang="ts">
    import { _ } from "svelte-i18n";
    import AccountIcon from "../../components/common/icons/AccountIcon.svelte";
    import AccountButton from "../../components/common/icons/buttons/AccountButton.svelte";
    import LogoutIcon from "../../components/common/icons/LogoutIcon.svelte";
    import { authenticationSession } from "$lib/service/auth-session";
    import { goto } from "$app/navigation";
    import TitleBar from "../../components/common/TitleBar.svelte";
    import { ModeratorPharmacyReview, ReviewModerationResult } from "$lib/service/data/moderator-pharmacy-review";
    import CenteredLoader from "../../components/common/widgets/CenteredLoader.svelte";
    import Review from "../../components/map/PharmacyView/Review.svelte";
    import { PAGER_LIMIT } from "$lib/service/data/pager";
    import IntersectionObserver from "svelte-intersection-observer/IntersectionObserver.svelte";
    import RateReviewButton from "../../components/common/icons/buttons/RateReviewButton.svelte";
    import CheckCircleIcon from "../../components/common/icons/CheckCircleIcon.svelte";
    import TimedDeletionIcon from "../../components/common/icons/TimedDeletionIcon.svelte";
    import IntermediateCheckIcon from "../../components/common/icons/IntermediateCheckIcon.svelte";
    import BellNotificationIcon from "../../components/common/icons/BellNotificationIcon.svelte";

    let accountDialog: HTMLDialogElement;
    let isAccountDialogVisible: boolean = $state(false);

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

<header>
    <a href="/" title="Go to home page">
        <img src="banner.svg" alt="Logo">
        <h3>{$_("mod.bannerName")}</h3>
    </a>
    <div class="buttons">
        <AccountButton
            size={48}
            title="My Account"
            on:click={(_) => {
                if (isAccountDialogVisible)
                    accountDialog.close();
                else
                    accountDialog.show();
                isAccountDialogVisible = !isAccountDialogVisible;
            }}
        />

        <dialog bind:this={accountDialog}>
            <a href="/mod/account/me">
                <AccountIcon size={32}/>
                <p>{$_("mod.header.myAccount")}</p>
            </a>
            <button onclick={() => {
                authenticationSession.logout();
                goto("/mod/login", { replaceState: true })
            }}>
                <LogoutIcon size={32}/>
                <p>{$_("mod.header.logOut")}</p>
            </button>
        </dialog>
    </div>
</header>

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
                <label for="moderated">{$_("mod.filters.showUnmoderated")}:</label>
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
                {:else if review.commentReviewResult === ReviewModerationResult.PersonalAttack}
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
                <RateReviewButton size={24} title="Rate the review"/>
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

<style>
    :global(body) {
        background-color: #fdcece;
    }

    .icon {
        display: inline-block;
        vertical-align: middle;
    }

    .buttons {
        display: flex;
        flex-direction: row;
        align-items: center;
    }

    header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        width: 100%;
        height: 64px;
        background-color: rgba(100, 0, 0, 0.2);
        border-bottom: 1px solid black;
        box-sizing: border-box;

        & > a {
            display: flex;
            align-items: center;
            user-select: none;
            & > img {
                margin-left: 0.5em;
                display: inline-block;
                height: 48px;
            }
            h3 {
                margin-left: 5px;
                color: white;
                text-shadow: 1px 1px 1px #eeeeee;
            }

            @media(max-width: 600px) {
                & > h3 {
                    display: none;
                }
            }
        }
    }

    dialog {
        text-align: center;
        background-color: white;
        position: absolute;
        border: 1px solid black;
        margin: 0;
        top: 64px;
        right: 0;
        left: auto;
        width: 180px;
        border-radius: 1em;
        padding: 0;
        overflow: hidden;
        height: fit-content;
        box-sizing: border-box;

        & > button, & > a {
            all: unset;
            display: flex;
            flex-direction: row;
            align-items: center;
            user-select: none;
            padding: 0;
            height: 48px;
            width: 100%;
            transition: all 0.3s ease-in-out;

            & > p {
                flex: 1;
            }

            &:hover {
                cursor: pointer;
                background-color: rgba(0, 0, 0, 0.15)
            }
        }
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