<script lang="ts">
    import { _ } from "svelte-i18n";
    import UserForm from "../../../components/common/widgets/UserForm.svelte";
    import LanguageSwitcher from "../../../components/common/LanguageSwitcher.svelte";

    let clientSideValidationErrors = $state(["Please complete the recaptcha challenge"])
</script>

<svelte:head>
    <title>Login - Pharmacy Finder</title>
</svelte:head>

<div class="container">
    <div class="lang-container">
        <LanguageSwitcher/>
    </div>
    <UserForm
        title={$_("mod.login.formTitle")}
        submitBtnText={$_("mod.login.submitTitle")}
    >
        <input type="text" name="username" placeholder={$_("mod.login.usernamePlaceholder")} required maxlength="32">
        <input type="password" name="password" placeholder={$_("mod.login.passwordPlaceholder")} required>
    </UserForm>
</div>

<style>
    :global(body) {
        background-color: #fdcece;
    }

    .lang-container {
        position: absolute;
        top: 0;
        left: 0;
    }

    .container {
        width: 100%;
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
    }
</style>