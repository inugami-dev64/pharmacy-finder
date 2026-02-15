<script lang="ts">
    import { _ } from "svelte-i18n";
    import { LocalizedBackendError } from "$lib/service/data/error";
    import { Moderation, ReviewModerationResult } from "$lib/service/data/moderation";
    import type { ModeratorPharmacyReview } from "$lib/service/data/moderator-pharmacy-review";
    import { moderatorModalZIndex } from "$lib/utils/z-indices";
    import ModalWindow from "../common/ModalWindow.svelte";
    import LimitedTextarea from "../common/widgets/LimitedTextarea.svelte";
    import { authenticationSession } from "$lib/service/auth-session";
    import PrimaryButton from "../common/widgets/buttons/PrimaryButton.svelte";
    import TrunctablePost from "../common/TrunctablePost.svelte";
    import CenteredLoader from "../common/widgets/CenteredLoader.svelte";
    import CheckCircleIcon from "../common/icons/CheckCircleIcon.svelte";
    import TimedDeletionIcon from "../common/icons/TimedDeletionIcon.svelte";
    import BellNotificationIcon from "../common/icons/BellNotificationIcon.svelte";

    let {
        onClose,
        review
    }: {
        onClose: () => void,
        review: ModeratorPharmacyReview
    } = $props();

    const userId: string|undefined = authenticationSession.getUserId(authenticationSession.getSessionToken() ?? "");

    let isLoading: boolean = $state(true);
    let pendingSubmission: boolean = $state(false);
    let moderations: Moderation[] = $state([]);
    let errorMsg: string = $state("");
    let userModeration: Moderation | undefined = $state(undefined);

    function fetchModerations() {
        isLoading = true;
        Moderation.getModerationsForReview(review.id || 0)
            .then(res => {
                moderations = res;
                const filtered = moderations.filter(v => v.moderatorId == userId)
                if (filtered.length > 0)
                    userModeration = filtered[0];
            })
            .catch(e => {
                if (e instanceof LocalizedBackendError)
                    errorMsg = e.msg;
            })
            .finally(() => isLoading = false);
    }

    $effect(() => fetchModerations());

    /**
     * Form submission callback function
     *
     * @param e specifies the SubmitEvent object
     */
    async function submitForm(e: SubmitEvent) {
        e.preventDefault();
        pendingSubmission = true;
        const form = e.target as HTMLFormElement;
        const data = new FormData(form);

        const moderationComment = data.get("moderation-comment");
        const markedForDeletion = data.get("mark-for-deletion");
        const moderationDecision = data.get("moderation-decision");

        const moderation: Moderation = new Moderation;
        moderation.moderatorComment = moderationComment?.toString();
        moderation.markedForDeletion = markedForDeletion?.toString() === "on";
        console.log(markedForDeletion?.toString())
        moderation.result = moderationDecision?.toString() as ReviewModerationResult|undefined;
        moderation.commentId = review.id;
        moderation.id = userModeration?.id;

        try {
            if (userModeration == null)
                await moderation.createModeration(review.id || 0);
            else await moderation.updateModeration();
            moderations = [];
            fetchModerations();
        } catch (e) {
            if (e instanceof LocalizedBackendError)
                errorMsg = $_(e.msg);
        } finally {
            pendingSubmission = false;
        }
    }
</script>

<ModalWindow zIndex={moderatorModalZIndex} onClose={onClose}>
    <div class="container">
        {#if errorMsg != ""}
            <p style="color: red">{errorMsg}</p>
        {/if}

        {#if isLoading}
            <CenteredLoader/>
        {:else}
            <form onsubmit={submitForm}>
                <LimitedTextarea
                    name="moderation-comment"
                    placeholder={$_("mod.moderations.form.moderationPlaceholder")}
                    text={userModeration?.moderatorComment || undefined}
                    maxLength={512}
                />
                <div class="markers">
                    <label for="mark-for-deletion">
                        {$_("mod.moderations.form.markForDeletion")}
                        <input type="checkbox" id="mark-for-deletion" name="mark-for-deletion" checked={userModeration != null && userModeration.markedForDeletion}>
                    </label>
                    <br>
                    <label for="moderation-decision">
                        {$_("mod.moderations.form.decisionLabel")}
                        <select name="moderation-decision">
                            <option value={ReviewModerationResult.Approved} selected={userModeration != null && userModeration.result == ReviewModerationResult.Approved}>{$_("mod.moderations.form.decision.approve")}</option>
                            <option value={ReviewModerationResult.PersonalAttack} selected={userModeration != null && userModeration.result == ReviewModerationResult.PersonalAttack}>{$_("mod.moderations.form.decision.personalAttack")}</option>
                            <option value={ReviewModerationResult.Offensive} selected={userModeration != null && userModeration.result == ReviewModerationResult.Offensive}>{$_("mod.moderations.form.decision.offensive")}</option>
                            <option value={ReviewModerationResult.Other} selected={userModeration != null && userModeration.result == ReviewModerationResult.Other}>{$_("mod.moderations.form.decision.other")}</option>
                        </select>
                    </label>
                </div>

                {#if pendingSubmission}
                    <CenteredLoader/>
                {:else}
                    {#if userModeration == null}
                        <PrimaryButton>{$_("mod.moderations.form.createNew")}</PrimaryButton>
                    {:else}
                        <PrimaryButton>{$_("mod.moderations.form.update")}</PrimaryButton>
                    {/if}
                {/if}
            </form>

            <hr>
            {#if moderations.length > 0}
                <p>{$_("mod.moderations.otherModerations")}:</p>
                {#each moderations as moderation}
                <TrunctablePost postText={moderation.moderatorComment || ""}>
                    <div class="mod-header">
                        <a href="/mod/accounts/{moderation.moderatorId}" target="_blank">{moderation.moderatorUsername}</a>
                        {#if !moderation.markedForDeletion && moderation.result === ReviewModerationResult.Approved}
                            <span title="{$_("mod.comments.approved")}" class="icon">
                                <CheckCircleIcon size={24}/>
                            </span>
                        {:else if moderation.markedForDeletion}
                            <span title="{$_("mod.comments.markedForDeletion")}" class="icon">
                                <TimedDeletionIcon size={24}/>
                            </span>
                        {/if}
                        {#if moderation.result === ReviewModerationResult.PersonalAttack}
                            <span title="{$_("mod.comments.personalAttack")}" class="icon">
                                <BellNotificationIcon size={24}/>
                            </span>
                        {:else if moderation.result === ReviewModerationResult.Offensive}
                            <span title="{$_("mod.comments.offensive")}" class="icon">
                                <BellNotificationIcon size={24}/>
                            </span>
                        {:else if moderation.result === ReviewModerationResult.Other}
                            <span title="{$_("mod.comments.other")}" class="icon">
                                <BellNotificationIcon size={24}/>
                            </span>
                        {/if}
                    </div>
                    <time>{new Date(moderation.reviewedAt ?? 0).toLocaleDateString()} {new Date(moderation.reviewedAt ?? 0).toLocaleTimeString()}</time>
                </TrunctablePost>
                {/each}
            {:else}
                <i>{$_("mod.moderations.noModerations")}</i>
            {/if}
        {/if}
    </div>
</ModalWindow>

<style>

    .mod-header {
        display: flex;

        & > a {
            font-weight: 400;
            transition: all 0.3s ease-in-out;

            &:hover {
                color: #999;
            }
        }

        & ~ time {
            display: block;
            margin-bottom: 0.25em;
        }
    }

    .container {
        width: 100%;
        height: 100%;
        overflow: auto;

        & > form {
            display: flex;
            flex-direction: column;

            & > .markers {
                margin: 10px 0;
                & > label {
                    user-select: none;
                }
            }
        }
    }
</style>