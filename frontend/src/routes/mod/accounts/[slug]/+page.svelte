<script lang="ts">
    import { authenticationSession } from "$lib/service/auth-session";
    import { _ } from "svelte-i18n";
    import type { LocalizedBackendError } from "$lib/service/data/error";
    import type { UserProfile } from "$lib/service/data/users";
    import ModeratorHeader from "../../../../components/mod/ModeratorHeader.svelte";
    import PrimaryButton from "../../../../components/common/widgets/buttons/PrimaryButton.svelte";
    import CenteredLoader from "../../../../components/common/widgets/CenteredLoader.svelte";

    const {
        data
    } : {
        data: { profile?: UserProfile, error?: LocalizedBackendError, isMe: boolean }
    } = $props();

    let allowEditing: boolean = $state(false);
    let isSubmitted: boolean = $state(false);

    // Make sure that when profile prop changes
    // we set allowEditing state value to false
    $effect(() => {
        data.profile;
        allowEditing = false;
    });

    const isAdmin: boolean = authenticationSession.isAdmin(authenticationSession.getSessionToken() || "");
</script>

<svelte:head>
    <title>Moderator accounts - Pharmacy Finder</title>
</svelte:head>

<ModeratorHeader
    bannerHref="/mod"
    bannerTitle={$_("mod.header.gotoMod")}
    isAdmin={isAdmin}
/>

<div class="container">
    <div class="profile-container">
        {#if data.isMe}
            <h1>Hi {data.profile?.firstName}! 👋</h1>
        {:else}
            <h1>{data.profile?.firstName} {data.profile?.lastName}</h1>
        {/if}
        <hr>

        <form class="profile-info">
            <span>
                🔒 Username:
                {data.profile?.username}
            </span>
            <span>
                ✉️ Email:
                {#if allowEditing}
                    <input type="email" name="email" value={data.profile?.email} placeholder="Email" required>
                {:else}
                    <a href="mailto:{data.profile?.email}">{data.profile?.email}</a>
                {/if}
            </span>
            <span>
                🪪 First name:
                {#if allowEditing}
                    <input type="text" name="firstName" value={data.profile?.firstName} placeholder="First name" required>
                {:else}
                    {data.profile?.firstName}
                {/if}
            </span>
            <span>
                🪪 Last name:
                {#if allowEditing}
                    <input type="text" name="lastName" value={data.profile?.lastName} placeholder="Last name" required>
                {:else}
                    {data.profile?.lastName}
                {/if}
            </span>
            <span>⏰ Registration timestamp: <time>{new Date(data.profile?.registrationTs || 0).toLocaleString()}</time></span>
            <span>⏰ Last login timestamp: <time>{new Date(data.profile?.lastLoginTs || 0).toLocaleString()}</time></span>
            <span>
                🛡️ Administrator:
                {#if allowEditing && isAdmin}
                    <input type="checkbox" name="isAdmin" checked={data.profile?.administrator} required>
                {:else}
                    {data.profile?.administrator ? "yes" : "no"}
                {/if}
            </span>
            {#if allowEditing && data.isMe}
                <span>
                    🔐 Current password:
                    <input type="password" name="currentPassword" placeholder="Current password">
                </span>
                <h2>Credentials</h2>
                <p>Leave empty if you do not want to change the password</p>
                <span>
                    🔐 New password:
                    <input type="password" name="newPassword" placeholder="New password">
                </span>
                <span>
                    🔐 Verify password:
                    <input type="password" name="verifyPassword" placeholder="Verify password">
                </span>
            {/if}

            {#if !allowEditing}
                <button onclick={() => allowEditing = !allowEditing}>
                    Allow editing
                </button>
            {:else if !isSubmitted}
                <div class="submit-container">
                    <PrimaryButton>Update user profile</PrimaryButton>
                </div>
            {:else}
                <CenteredLoader/>
            {/if}
        </form>

        <span class="user-id"><i>User ID: <a href="/mod/accounts/{data.profile?.id}">{data.profile?.id}</a></i></span>
    </div>
</div>

<style>
    :global(body) {
        background-color: #fdcece;
    }

    h1, h2, p {
        text-align: center;
    }

    .container {
        display: flex;
        justify-content: center;
        align-items: center;
        width: 100%;
        height: 100%;

        & > .profile-container {
            width: 30%;
            min-width: 400px;
            background-color: #fff;
            padding: 0.5em;
            border-radius: 0.5em;

            & > .user-id {
                color: #aaa;
                display: block;
                text-align: center;
                margin-top: 20px;

                & > i > a {
                    color: #aaa;
                    transition: 0.3s all ease-in-out;

                    &:hover {
                        color: #777;
                    }
                }
            }

            & > .profile-info {
                & > span {
                    margin: 5px 0;
                    display: block;
                }

                & > .submit-container {
                    margin-top: 20px;
                    display: flex;
                    flex-direction: column;
                    width: 100%;
                }
            }
        }
    }
</style>