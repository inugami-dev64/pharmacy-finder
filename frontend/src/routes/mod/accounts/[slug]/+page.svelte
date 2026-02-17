<script lang="ts">
    import { authenticationSession } from "$lib/service/auth-session";
    import { _ } from "svelte-i18n";
    import { LocalizedBackendError } from "$lib/service/data/error";
    import { UserProfile } from "$lib/service/data/users";
    import ModeratorHeader from "../../../../components/mod/ModeratorHeader.svelte";
    import PrimaryButton from "../../../../components/common/widgets/buttons/PrimaryButton.svelte";
    import CenteredLoader from "../../../../components/common/widgets/CenteredLoader.svelte";
    import { goto } from "$app/navigation";
    import { onMount } from "svelte";
    import { localeSwitcher } from "$lib/service/locale";

    const {
        data
    } : {
        data: { profile?: UserProfile, error?: LocalizedBackendError, isMe: boolean }
    } = $props();

    onMount(async () => {
        localeSwitcher.setDefault();
    });

    let allowEditing: boolean = $state(false);
    let pendingSubmission: boolean = $state(false);
    let errorMessage: string = $state("");
    let successMessage: string = $state("");

    // Make sure that when profile prop changes
    // we set allowEditing state value to false
    $effect(() => {
        data.profile;
        errorMessage = "";
        successMessage = "";
        allowEditing = false;
    });

    const isAdmin: boolean = authenticationSession.isAdmin(authenticationSession.getSessionToken() || "");

    /**
     * Form submission callback function
     *
     * @param e specifies the SubmitEvent object
     */
    async function submitForm(e: SubmitEvent) {
        e.preventDefault();
        pendingSubmission = true;
        const form: HTMLFormElement = e.target as HTMLFormElement;
        const formData: FormData = new FormData(form);

        const email = formData.get("email")?.toString();
        const firstName = formData.get("firstName")?.toString();
        const lastName = formData.get("lastName")?.toString();
        const isAdmin = formData.get("isAdmin")?.toString() === "on";
        const currentPassword = formData.get("currentPassword")?.toString();
        const newPassword = formData.get("newPassword")?.toString();
        const verifyPassword = formData.get("verifyPassword")?.toString();

        // If newPassword or verifyPassword values were provided, check if they match
        if ((newPassword != null || verifyPassword != null) && newPassword !== verifyPassword) {
            errorMessage = $_("mod.accounts.status.passwordsDoNotMatch");
            pendingSubmission = false;
            return;
        }

        const user = new UserProfile();
        user.email = email;
        user.firstName = firstName;
        user.lastName = lastName;
        user.administrator = isAdmin;

        const previousIsAdmin = data.profile?.administrator || false;

        try {
            if (data.isMe) {
                await user.updateCurrentlyAuthenticatedUser(currentPassword || "", newPassword);
                data.profile = user;
            } else {
                user.id = data.profile?.id;
                await user.updateUser();
                data.profile = user;
            }
        } catch (e) {
            if (e instanceof LocalizedBackendError)
                errorMessage = $_(e.msg);
            pendingSubmission = false;
            return;
        }

        successMessage = $_("mod.accounts.status.success");
        errorMessage = "";
        pendingSubmission = false;
        allowEditing = false;

        // If password or administrator state has been updated, the user shall be logged out
        if (data.isMe && (newPassword != "" || isAdmin != previousIsAdmin)) {
            setTimeout(() => {
                authenticationSession.logout();
                goto("/mod/login", { replaceState: true });
            }, 1000);
        }
    }
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
            <h1>{$_("mod.accounts.greeting")} {data.profile?.firstName}! 👋</h1>
        {:else}
            <h1>{data.profile?.firstName} {data.profile?.lastName}</h1>
        {/if}
        <hr>

        {#if errorMessage !== ""}
            <p style="color: red; text-align: center">{errorMessage}</p>
        {:else if successMessage !== ""}
            <p style="color: green; text-align: center">{successMessage}</p>
        {/if}

        <form class="profile-info" onsubmit={submitForm}>
            <span>
                🔒 {$_("mod.accounts.username")}:
                {data.profile?.username}
            </span>
            <span>
                ✉️ {$_("mod.accounts.email")}:
                {#if allowEditing}
                    <input type="email" name="email" value={data.profile?.email} placeholder="{$_("mod.accounts.email")}" required maxlength="64">
                {:else}
                    <a href="mailto:{data.profile?.email}">{data.profile?.email}</a>
                {/if}
            </span>
            <span>
                🪪 {$_("mod.accounts.firstName")}:
                {#if allowEditing}
                    <input type="text" name="firstName" value={data.profile?.firstName} placeholder="{$_("mod.accounts.firstName")}" required maxlength="64">
                {:else}
                    {data.profile?.firstName}
                {/if}
            </span>
            <span>
                🪪 {$_("mod.accounts.lastName")}:
                {#if allowEditing}
                    <input type="text" name="lastName" value={data.profile?.lastName} placeholder="{$_("mod.accounts.lastName")}" required maxlength="64">
                {:else}
                    {data.profile?.lastName}
                {/if}
            </span>
            <span>⏰ {$_("mod.accounts.registrationTimestamp")}: <time>{new Date(data.profile?.registrationTs || 0).toLocaleString()}</time></span>
            <span>⏰ {$_("mod.accounts.lastLoginTimestamp")}: <time>{new Date(data.profile?.lastLoginTs || 0).toLocaleString()}</time></span>
            <span>
                🛡️ {$_("mod.accounts.administrator")}:
                {#if allowEditing && isAdmin}
                    <input type="checkbox" name="isAdmin" checked={data.profile?.administrator}>
                {:else}
                    {data.profile?.administrator ? "yes" : "no"}
                {/if}
            </span>
            {#if allowEditing && data.isMe}
                <span>
                    🔐 {$_("mod.accounts.currentPassword")}:
                    <input type="password" name="currentPassword" placeholder="{$_("mod.accounts.currentPassword")}" required maxlength="72">
                </span>
                <h2>{$_("mod.accounts.credentials.title")}</h2>
                <p>{$_("mod.accounts.credentials.description")}</p>
                <span>
                    🔐 {$_("mod.accounts.newPassword")}:
                    <input type="password" name="newPassword" placeholder="{$_("mod.accounts.newPassword")}" maxlength="72">
                </span>
                <span>
                    🔐 {$_("mod.accounts.verifyPassword")}:
                    <input type="password" name="verifyPassword" placeholder="{$_("mod.accounts.verifyPassword")}" maxlength="72">
                </span>
            {/if}

            {#if !allowEditing && (data.isMe || isAdmin)}
                <button onclick={() => {
                    allowEditing = !allowEditing;
                    successMessage = "";
                }}>
                    {$_("mod.accounts.allowEditing")}
                </button>
            {:else if allowEditing && !pendingSubmission}
                <div class="submit-container">
                    <PrimaryButton>{$_("mod.accounts.update")}</PrimaryButton>
                </div>
            {:else if allowEditing}
                <CenteredLoader/>
            {/if}
        </form>

        <span class="user-id"><i>{$_("mod.accounts.userId")}: <a href="/mod/accounts/{data.profile?.id}">{data.profile?.id}</a></i></span>
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