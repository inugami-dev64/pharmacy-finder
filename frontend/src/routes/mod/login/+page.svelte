<script lang="ts">
    import { _ } from "svelte-i18n";
    import UserForm from "../../../components/common/widgets/UserForm.svelte";
    import LanguageSwitcher from "../../../components/common/LanguageSwitcher.svelte";
    import { LoginForm } from "$lib/service/data/auth";
    import { authenticationSession } from "$lib/service/auth-session";
    import { goto } from "$app/navigation";
    import { LocalizedBackendError } from "$lib/service/data/error";

    let clientSideValidationErrors: Array<string> = $state([])
    let isSubmitted: boolean = $state(false);

    async function submitForm(e: SubmitEvent) {
        e.preventDefault();
        clientSideValidationErrors = [];
        isSubmitted = true;
        const form = e.target as HTMLFormElement;
        const data = new FormData(form);

        const loginForm: LoginForm = new LoginForm();
        loginForm.username = data.get("username")?.toString();
        loginForm.password = data.get("password")?.toString();

        try {
            const authUser = await loginForm.login(fetch);
            authenticationSession.setSessionToken(authUser.session?.token || "", authUser.session?.validFor || 0);
            goto("/mod", { replaceState: true });
        } catch (e) {
            if (e instanceof LocalizedBackendError)
                clientSideValidationErrors.push(e.msg);
        }

        isSubmitted = false;
    }
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
        submitForm={submitForm}
        validationErrors={clientSideValidationErrors}
        isSubmitted={isSubmitted}
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