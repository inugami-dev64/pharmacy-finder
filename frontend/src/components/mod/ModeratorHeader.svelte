<script lang="ts">
    import { _ } from "svelte-i18n";
    import LanguageDropdownMenu from "../common/icons/dropdown/LanguageDropdownMenu.svelte";
    import DropdownMenuContainer from "../common/icons/dropdown/DropdownMenuContainer.svelte";
    import AccountButton from "../common/icons/buttons/AccountButton.svelte";
    import AccountIcon from "../common/icons/AccountIcon.svelte";
    import { authenticationSession } from "$lib/service/auth-session";
    import { goto } from "$app/navigation";
    import LogoutIcon from "../common/icons/LogoutIcon.svelte";
    import ShieldButton from "../common/icons/buttons/ShieldButton.svelte";
    import BANNER from "$lib/assets/banner.svg"

    let {
        bannerHref = "/",
        bannerTitle = "Go to home page",
        isAdmin = false
    }: {
        bannerHref?: string,
        bannerTitle?: string,
        isAdmin?: boolean
    } = $props();
</script>

<header>
    <a href="{bannerHref}" title="{bannerTitle}">
        <img src={BANNER} alt="Logo">
        <h3>{$_("mod.header.name")}</h3>
    </a>
    <div class="buttons">
        {#if isAdmin === true}
            <DropdownMenuContainer
                size={48}
                ButtonComponent={ShieldButton}
                title={$_("mod.header.admin")}
                color="#fff"
            >
                <a href="/mod/accounts">
                    <AccountIcon size={32}/>
                    <p>{$_("mod.header.administerAccounts")}</p>
                </a>
            </DropdownMenuContainer>
        {/if}
        <LanguageDropdownMenu
            size={48}
            color="#fff"
        />
        <DropdownMenuContainer
            size={48}
            ButtonComponent={AccountButton}
            title={$_("mod.header.myAccount")}
            color="#fff"
        >
            <a href="/mod/accounts/me">
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
        </DropdownMenuContainer>
    </div>
</header>

<style>
    header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        width: 100%;
        height: 64px;
        background-color: rgba(100, 0, 0, 0.2);
        border-bottom: 1px solid black;
        box-sizing: border-box;

        & > .buttons {
            display: flex;
            flex-direction: row;
            align-items: center;
        }

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
</style>