<script lang="ts">
    import type { PharmacyInfo } from "$lib/service/pharmacy-info";
    import { PharmacyReview } from "$lib/service/pharmacy-review";
    import { modifyReviewZIndex } from "$lib/utils/z-indices";
    import ModalWindow from "../ModalWindow.svelte";
    import Recaptcha from "../../common/Recaptcha.svelte";
    import DangerButton from "../../common/widgets/buttons/DangerButton.svelte";
    import Loader from "../../common/widgets/Loader.svelte";
    import { _ } from "svelte-i18n";

    let {
        pharmacy,
        onClose,
        review
    }: {pharmacy: PharmacyInfo, onClose: () => void, review: PharmacyReview} = $props();

    // stateful values
    let pendingSubmission: boolean = $state(false);
    let missingCaptcha: boolean = $state(false);
    let invalidModCode: boolean = $state(false);

    async function submitForm(e: SubmitEvent) {
        e.preventDefault();
        pendingSubmission = true;
        const form = e.target as HTMLFormElement;
        const data = new FormData(form);

        const modCode = data.get("mod-code");
        const recaptchaResponse = data.get("g-recaptcha-response");

        if (recaptchaResponse == null || recaptchaResponse.toString() === "") {
            pendingSubmission = false;
            missingCaptcha = true;
            return;
        }

        let newReview: PharmacyReview = new PharmacyReview;
        Object.assign(newReview, review);
        newReview.modCode = modCode?.toString();

        newReview.__gRecaptchaResponse = recaptchaResponse.toString();

        try {
            await newReview.deleteReview(pharmacy.id ?? 0);
            onClose();
        } catch (e) {
            pendingSubmission = false;
            invalidModCode = true;
            setTimeout(() => onClose(), 2000);
        }

    }
</script>

<ModalWindow zIndex={modifyReviewZIndex} onClose={onClose} minHeight={400}>
    <h3>{pharmacy.name}</h3>
    <form onsubmit={submitForm}>
        <div>
            <label for="mod-code">{$_("map.reviewForm.modCodeTitle")}*:</label><br>
            <input type="text" name="mod-code" required/>
        </div>

        <div class="btn-captcha">
            {#if pendingSubmission}
                <div style="display flex; justify-content: center; width: 100%">
                    <Loader/>
                </div>
            {:else if invalidModCode}
                <p style="color: red">{$_("map.reviewForm.responses.invalidModCode")}</p>
            {:else}
                {#if missingCaptcha}
                    <p style="color: red">{$_("map.reviewForm.responses.missingCaptcha")}</p>
                {/if}
                <Recaptcha/>
                <DangerButton>{$_("map.reviewForm.actions.deleteReview")}</DangerButton>
            {/if}
        </div>
    </form>
</ModalWindow>

<style>
    h3 {
        color: red;
        padding: 0;
        margin-bottom: 0.5em;
    }

    form {
        display: flex;
        width: 100%;
        flex: 1;
        flex-direction: column;
        justify-content: space-between;
    }

    input {
        min-width: 200px;
        width: 50%;
        padding: 0.5em;
        border-radius: 0.5em;
        border: 1px solid black;
    }

    .btn-captcha {
        display: flex;
        flex-direction: column;
        width: 100%;
    }
</style>