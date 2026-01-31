<script lang="ts">
    import { _ } from "svelte-i18n";
    import LoginForm from "../../../components/common/widgets/LoginForm.svelte";
    import Recaptcha from "../../../components/common/Recaptcha.svelte";

    let clientSideValidationErrors = $state(["Please complete the recaptcha challenge"])
</script>

<svelte:head>
    <title>Login - Pharmacy Finder</title>
</svelte:head>

<main>
    <LoginForm
        title="Login to moderation panel"
        submitBtnText="Login"
    >
        <input type="text" name="username" placeholder="Username" required maxlength="32">
        <input type="password" name="password" placeholder="Password" required>
    </LoginForm>
</main>

<style>
    :global(body) {
        background-color: #fdcece;
    }

    main {
        width: 100%;
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
    }
</style>