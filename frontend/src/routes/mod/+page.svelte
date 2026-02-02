<script lang="ts">
    import { _ } from "svelte-i18n";
    import AccountIcon from "../../components/common/icons/AccountIcon.svelte";
    import AccountButton from "../../components/common/icons/buttons/AccountButton.svelte";

    let accountDialog: HTMLDialogElement;
    let isAccountDialogVisible: boolean = $state(false);
</script>

<svelte:head>
    <title>Moderation panel - Pharmacy Finder</title>
</svelte:head>

<header>
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
                <p>My account</p>
            </a>
            <button>
                <LogoutIcon size={32}/>
                <p>Log out</p>
            </button>
        </dialog>
    </div>
</header>



<style>
    :global(body) {
        background-color: #fdcece;
    }

    .buttons {
        display: flex;
        flex-direction: row;
        align-items: center;
    }

    header {
        display: flex;
        justify-content: end;
        width: 100%;
        height: 64px;
        background-color: rgba(100, 0, 0, 0.2);
        border-bottom: 1px solid black;
    }

    dialog {
        text-align: center;
        background-color: white;
        position: absolute;
        border: none;
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
</style>