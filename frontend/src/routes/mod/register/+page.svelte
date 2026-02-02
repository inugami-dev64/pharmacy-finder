<script lang="ts">
    import { _ } from "svelte-i18n";
    import UserForm from "../../../components/common/widgets/UserForm.svelte";
    import LanguageSwitcher from "../../../components/common/LanguageSwitcher.svelte";
    import { RegistrationForm } from "$lib/service/data/auth";
    import { LocalizedBackendError } from "$lib/service/data/error";
    import { authenticationSession } from "$lib/service/auth-session";
    import { goto } from "$app/navigation";

    let clientSideValidationErrors: Array<string> = $state([]);
    let isSubmitted: boolean = $state(false);

    async function submitForm(e: SubmitEvent) {
        e.preventDefault();
        clientSideValidationErrors = [];
        isSubmitted = true
        const form = e.target as HTMLFormElement;
        const data = new FormData(form);

        const username = data.get("username");
        const email = data.get("email");
        const password = data.get("password");
        const confirmPassword = data.get("confirmPassword");
        const firstName = data.get("firstName");
        const lastName = data.get("lastName");

        if (password != confirmPassword) {
            clientSideValidationErrors.push("mod.register.errors.passwordsDoNotMatch");
            isSubmitted = false;
            return
        }

        let regForm = new RegistrationForm();
        regForm.username = username?.toString(),
        regForm.email = email?.toString(),
        regForm.password = password?.toString(),
        regForm.firstName = firstName?.toString(),
        regForm.lastName = lastName?.toString()

        try {
            const adminUser = await regForm.registerAdministrator(fetch);
            authenticationSession.setSessionToken(adminUser.session?.token || "", adminUser.session?.validFor || 0);
            goto("/mod", { replaceState: true });
        } catch (e) {
            if (e instanceof LocalizedBackendError)
                clientSideValidationErrors.push(e.msg);
        }

        isSubmitted = false;
    }
</script>

<svelte:head>
    <title>Register - Pharmacy Finder</title>
</svelte:head>

<div class="container">
    <div class="lang-container">
        <LanguageSwitcher/>
    </div>
    <UserForm
        title={$_("mod.register.formTitle")}
        submitBtnText={$_("mod.register.submitTitle")}
        isSubmitted={isSubmitted}
        validationErrors={clientSideValidationErrors}
        submitForm={submitForm}
    >
        <input type="text" name="username" placeholder={$_("mod.register.usernamePlaceholder")} required maxlength="32">
        <input type="email" name="email" placeholder={$_("mod.register.emailPlaceholder")} required maxlength="64">
        <input type="password" name="password" placeholder={$_("mod.register.passwordPlaceholder")} required>
        <input type="password" name="confirmPassword" placeholder={$_("mod.register.confirmPasswordPlaceholder")} required>
        <input type="text" name="firstName" placeholder={$_("mod.register.firstName")} required maxlength="64">
        <input type="text" name="lastName" placeholder={$_("mod.register.lastName")} required maxlength="64">
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