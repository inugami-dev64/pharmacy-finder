<script lang="ts">
    import { _ } from "svelte-i18n";
    import ModeratorHeader from "../../../components/mod/ModeratorHeader.svelte";
    import { authenticationSession } from "$lib/service/auth-session";
    import ModerationListContainer from "../../../components/mod/ModerationListContainer.svelte";
    import TitleBar from "../../../components/common/TitleBar.svelte";
    import type { HealthCheckResult } from "$lib/service/data/health";
    import type { UserProfile } from "$lib/service/data/users";
    import { onMount } from "svelte";
    import AccountIcon from "../../../components/common/icons/AccountIcon.svelte";
    import { localeSwitcher } from "$lib/service/locale";

    let {
        data
    }: {
        data: { health: HealthCheckResult, accounts: UserProfile[] }
    } = $props();

    let orderByValue: string | undefined = $state("username");
    let desc: boolean = $state(false);
    let accounts: UserProfile[] = $state([...data.accounts]);

    function sortAccounts() {
        if (!desc) {
            switch (orderByValue) {
                case "username":
                    accounts.sort((a, b) => a.username?.localeCompare(b.username || "") || 0);
                    break;
                case "firstName":
                    accounts.sort((a, b) => a.firstName?.localeCompare(b.firstName || "") || 0);
                    break;
                case "lastName":
                    accounts.sort((a, b) => a.lastName?.localeCompare(b.lastName || "") || 0);
                    break;
                case "email":
                    accounts.sort((a, b) => a.email?.localeCompare(b.email || "") || 0);
                    break;
                case "registrationTime":
                    accounts.sort((a, b) => (a.registrationTs || 0) < (b.registrationTs || 0) ? -1 : ((a.registrationTs || 0) == (b.registrationTs || 0) ? 0 : 1));
                    break;
                case "lastLoginTime":
                    accounts.sort((a, b) => (a.lastLoginTs || 0) < (b.lastLoginTs || 0) ? -1 : ((a.lastLoginTs || 0) == (b.lastLoginTs || 0) ? 0 : 1));
                    break;
            }
        } else {
            switch (orderByValue) {
                case "username":
                    accounts.sort((b, a) => a.username?.localeCompare(b.username || "") || 0);
                    break;
                case "firstName":
                    accounts.sort((b, a) => a.firstName?.localeCompare(b.firstName || "") || 0);
                    break;
                case "lastName":
                    accounts.sort((b, a) => a.lastName?.localeCompare(b.lastName || "") || 0);
                    break;
                case "email":
                    accounts.sort((b, a) => a.email?.localeCompare(b.email || "") || 0);
                    break;
                case "registrationTime":
                    accounts.sort((b, a) => (a.registrationTs || 0) < (b.registrationTs || 0) ? -1 : ((a.registrationTs || 0) == (b.registrationTs || 0) ? 0 : 1));
                    break;
                case "lastLoginTime":
                    accounts.sort((b, a) => (a.lastLoginTs || 0) < (b.lastLoginTs || 0) ? -1 : ((a.lastLoginTs || 0) == (b.lastLoginTs || 0) ? 0 : 1));
                    break;
            }
        }
    }

    onMount(() => {
        sortAccounts();
        localeSwitcher.setDefault();
    });
</script>

<svelte:head>
    <title>Moderator accounts - Pharmacy Finder</title>
</svelte:head>

<ModeratorHeader
    bannerHref="/mod"
    bannerTitle={$_("mod.header.gotoMod")}
    isAdmin={authenticationSession.isAdmin(authenticationSession.getSessionToken() || "")}
/>

<ModerationListContainer>
    <TitleBar>
        <div>
            <label for="orderBy">{$_("mod.filters.orderBy.title")}:</label>
            <select
                name="orderBy"
                bind:value={orderByValue}
                onchange={() => sortAccounts()}
            >
                <option selected value="username">{$_("mod.accounts.sorting.username")}</option>
                <option value="firstName">{$_("mod.accounts.sorting.firstName")}</option>
                <option value="lastName">{$_("mod.accounts.sorting.lastName")}</option>
                <option value="email">{$_("mod.accounts.sorting.email")}</option>
                <option value="registrationTime">{$_("mod.accounts.sorting.registrationTime")}</option>
                <option value="lastLoginTime">{$_("mod.accounts.sorting.lastLoginTime")}</option>
            </select>
            <label for="desc">
                {$_("mod.accounts.sorting.desc")}:
                <input
                    type="checkbox"
                    name="desc"
                    bind:checked={desc}
                    onchange={() => sortAccounts()}
                >
            </label>
        </div>
    </TitleBar>

    <h1>{$_("mod.accounts.listTitle")}</h1>
    {#each accounts as account}
        <user-listing>
            <a href="/mod/accounts/{account.id}" data-sveltekit-preload-data="tap">
                <AccountIcon size={22} color="#000"/>
                <p>{account.username}</p>
            </a>
            <p>- {account.firstName} {account.lastName} -</p>
            <a href="mailto:{account.email}">{account.email}</a>
        </user-listing>
    {/each}
</ModerationListContainer>


<style>
    :global(body) {
        background-color: #fdcece;
    }

    h1 {
        text-align: center;
    }

    user-listing {
        display: inline-block;
        width: calc(100% - 8px);
        box-sizing: border-box;
        padding: 0.25em;
        border-radius: 0.5em;
        background-color: #eee;
        border: 1px solid black;
        margin: 4px;

        :global(*) {
            display: inline-block;
            vertical-align: middle;
        }

        p {
            margin: 0;
        }
    }
</style>