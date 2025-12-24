<script lang="ts">
    import { PharmacyReview } from "$lib/service/pharmacy-review";
    import Countries from "$lib/assets/countries.json";
    import type { PharmacyInfo } from "$lib/service/pharmacy-info";
    import Loader from "../../common/widgets/Loader.svelte";
    import PrimaryButton from "../../common/widgets/buttons/PrimaryButton.svelte";
    import StarPicker from "../../common/widgets/stars/StarPicker.svelte";
    import ModalWindow from "../../common/ModalWindow.svelte";
    import Recaptcha from "../../common/Recaptcha.svelte";
    import { _ } from "svelte-i18n";

    let {
        pharmacy,
        onClose,
        review,
    }: {pharmacy: PharmacyInfo, onClose: () => void, review?: PharmacyReview} = $props();

    // stateful values
    let missingCaptcha: boolean = $state(false);
    let pendingSubmission: boolean = $state(false);
    let successfullyModified: boolean = $state(false);
    let invalidModCode: boolean = $state(false);
    let newReview: PharmacyReview = $state(new PharmacyReview);

    /**
     * Form submission callback function
     *
     * @param e specifies the SubmitEvent object
     */
    async function submitForm(e: SubmitEvent) {
        e.preventDefault();
        pendingSubmission = true;
        const form = e.target as HTMLFormElement;
        const data = new FormData(form);

        const stars = data.get("stars");
        const nationality = data.get("nationality");
        const comment = data.get("review-comment");
        const hrtKind = data.get("hrt-kind");
        const prescriptionType = data.get("prescription-type");
        const modCode = data.get("mod-code");
        const recaptchaResponse = data.get("g-recaptcha-response");

        if (recaptchaResponse == null || recaptchaResponse.toString() === "") {
            pendingSubmission = false;
            missingCaptcha = true;
            return;
        }

        newReview.stars = Number.parseInt(stars?.toString() ?? "5");
        newReview.nationality = nationality?.toString() ?? "EE";
        newReview.review = comment?.toString();
        newReview.hrtKind = hrtKind?.toString();
        newReview.prescriptionType = prescriptionType?.toString();
        newReview.modCode = modCode?.toString();
        newReview.__gRecaptchaResponse = recaptchaResponse as string ?? "";

        if (review?.id) {
            newReview.id = review.id;

            // TODO: Better error handling
            try {
                newReview = await newReview.updateReview(pharmacy.id ?? 0);
                successfullyModified = true;
            } catch (e) {
                invalidModCode = true;
                setTimeout(() => onClose(), 2000);
            }
        } else {
            newReview = await newReview.createReview(pharmacy.id ?? 0);
        }

        pendingSubmission = false;
    }
</script>

<ModalWindow zIndex={10000} onClose={onClose}>
    <h3>{pharmacy.name}</h3>
    <form onsubmit={submitForm}>
        <div class="form-contents">
            <label for="stars">{$_("map.reviewForm.ratingTitle")}*:</label>
            <StarPicker name="stars" defaultChecked={review?.stars ?? 5}/>
            <label for="review-comment">{$_("map.reviewForm.commentTitle")}:</label>
            <textarea
                name="review-comment"
                placeholder="{$_("map.reviewForm.commentPlaceholder")}"
            >{review?.review}</textarea>
            <div>
                <label for="nationality">{$_("map.reviewForm.nationalityTitle")}*:</label>
                <select name="nationality" required>
                    {#each Object.keys(Countries).sort() as key}
                        {#if review == null && key === "EE" || review != null && key === review.nationality}
                        <option value={key} selected>{`${Countries[key as keyof typeof Countries].emoji} ${Countries[key as keyof typeof Countries].name}`}</option>
                        {:else}
                        <option value={key}>{`${Countries[key as keyof typeof Countries].emoji} ${Countries[key as keyof typeof Countries].name}`}</option>
                        {/if}
                    {/each}
                </select>
            </div>
            <div>
                <label for="hrt-kind">{$_("map.reviewForm.hrtKindTitle")}*:</label>
                <select name="hrt-kind" required>
                    <option value="e" selected={review != null && review.hrtKind === 'e'}>{$_("map.reviewForm.hrtKind.estrogen")}</option>
                    <option value="t" selected={review != null && review.hrtKind === 't'}>{$_("map.reviewForm.hrtKind.testosterone")}</option>
                </select>
            </div>
            <div>
                <label for="prescription-type">{$_("map.reviewForm.prescriptionIssuerTitle")}*:</label>
                <select name="prescription-type" required>
                    <option value="Imago" selected={review != null && review.prescriptionType === "Imago"}>Imago</option>
                    <option value="GenderGP" selected={review != null && review.prescriptionType === "GenderGP"}>GenderGP</option>
                    <option value="National" selected={review != null && review.prescriptionType === "National"}>National</option>
                </select>
            </div>

            {#if review != null}
                <div>
                    <label for="mod-code">{$_("map.reviewForm.modCodeTitle")}*:</label>
                    <input type="text" name="mod-code" required>
                </div>
            {/if}
        </div>

        {#if pendingSubmission}
            <div style="display: flex; justify-content: center; width: 100%">
                <Loader/>
            </div>
        {:else if successfullyModified}
            <p style="color: green">{$_("map.reviewForm.responses.modSuccess")}</p>
        {:else if invalidModCode}
            <p style="color: red">{$_("map.reviewForm.responses.invalidModCode")}</p>
        {:else if newReview.modCode}
            <p>{$_("map.reviewForm.responses.newModCode")}: <span style="color: green">{newReview.modCode}</span></p>
        {:else}
            {#if missingCaptcha}
                <p style="color: red">{$_("map.reviewForm.responses.missingCaptcha")}</p>
            {/if}
            <Recaptcha/>
            <PrimaryButton>
                {review == null ? $_("map.reviewForm.actions.createReview") : $_("map.reviewForm.actions.modifyReview")}
            </PrimaryButton>
        {/if}
    </form>
</ModalWindow>


<style>
    h3 {
        padding: 0;
        margin-bottom: 0.5em;
    }

    form {
        display: flex;
        overflow: auto;
        width: 100%;
        flex: 1;
        flex-direction: column;
        justify-content: space-between;
    }

    .form-contents > label {
        display: block;
        padding: 0.5em 0;
    }

    .form-contents > textarea {
        resize: none;
        width: calc(100%);
        height: 150px;
        box-sizing: border-box;
    }
</style>