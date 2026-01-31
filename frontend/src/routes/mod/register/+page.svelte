<script lang="ts">
    import { _ } from "svelte-i18n";
    import UserForm from "../../../components/common/widgets/UserForm.svelte";
    import LanguageSwitcher from "../../../components/common/LanguageSwitcher.svelte";

    let clientSideValidationErrors = $state(["Please complete the recaptcha challenge"])
</script>

<svelte:head>
    <title>Login - Pharmacy Finder</title>
</svelte:head>

<main>
    <div class="lang-container">
        <LanguageSwitcher/>
    </div>
    <UserForm
        title={$_("mod.register.formTitle")}
        submitBtnText={$_("mod.register.submitTitle")}
    >
        <input type="text" name="username" placeholder={$_("mod.register.usernamePlaceholder")} required maxlength="32">
        <input type="email" name="email" placeholder={$_("mod.register.emailPlaceholder")} required maxlength="64">
        <input type="password" name="password" placeholder={$_("mod.register.passwordPlaceholder")} required>
        <input type="password" name="confirmPassword" placeholder={$_("mod.register.confirmPasswordPlaceholder")} required>
        <input type="text" name="firstName" placeholder={$_("mod.register.firstName")} required maxlength="64">
        <input type="text" name="lastName" placeholder={$_("mod.register.lastName")} required maxlength="64">
    </UserForm>
</main>

<style>
    :global(body) {
        background-color: #fdcece;
    }

    .lang-container {
        position: absolute;
        top: 0;
        left: 0;
    }

    main {
        width: 100%;
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
    }
</style>