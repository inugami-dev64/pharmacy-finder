<script lang="ts">
    import { PharmacyReview } from "$lib/service/pharmacy-review";
    import TitleBar from "../../common/TitleBar.svelte";
    import Countries from "$lib/assets/countries.json";
    import type { PharmacyInfo } from "$lib/service/pharmacy-info";
    import Loader from "../../common/widgets/Loader.svelte";
    import PrimaryButton from "../../common/widgets/buttons/PrimaryButton.svelte";
    import StarPicker from "../../common/widgets/stars/StarPicker.svelte";

    let {
        pharmacy,
        onClose,
        review,
    }: {pharmacy: PharmacyInfo, onClose: () => void, review?: PharmacyReview} = $props();

    // stateful values
    let pendingSubmission: boolean = $state(false);
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

        newReview.stars = Number.parseInt(stars?.toString() ?? "5");
        newReview.nationality = nationality?.toString() ?? "EE";
        newReview.review = comment?.toString();
        newReview.hrtKind = hrtKind?.toString();
        newReview.prescriptionType = prescriptionType?.toString();
        newReview.modCode = modCode?.toString();

        if (review?.id) {
            newReview.id = review.id;

            // TODO: Better error handling
            try {
                newReview = await newReview.updateReview(pharmacy.id ?? 0);
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

<div class="form-container">
    <div class="review-form">
        <TitleBar onClose={onClose}/>
        <h3>{pharmacy.name}</h3>
        <form onsubmit={submitForm}>
            <div class="form-contents">
                <label for="stars">Your rating*:</label>
                <StarPicker name="stars" defaultChecked={5}/>
                <label for="review-comment">Comment:</label>
                <textarea
                    name="review-comment"
                    placeholder="Enter your comment here"
                >{review?.review}</textarea>
                <div>
                    <label for="nationality">Nationality*:</label>
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
                    <label for="hrt-kind">HRT kind*:</label>
                    <select name="hrt-kind" required>
                        <option value="e" selected={review != null && review.hrtKind === 'e'}>Estrogen</option>
                        <option value="t" selected={review != null && review.hrtKind === 't'}>Testosterone</option>
                    </select>
                </div>
                <div>
                    <label for="prescription-type">Prescription type*:</label>
                    <select name="prescription-type" required>
                        <option value="Imago" selected={review != null && review.prescriptionType === "Imago"}>Imago</option>
                        <option value="GenderGP" selected={review != null && review.prescriptionType === "GenderGP"}>GenderGP</option>
                        <option value="National" selected={review != null && review.prescriptionType === "National"}>National</option>
                    </select>
                </div>

                {#if review != null}
                    <div>
                        <label for="mod-code">Modification code*:</label>
                        <input type="text" name="mod-code" required>
                    </div>
                {/if}
            </div>

            {#if pendingSubmission}
            <div style="display: flex; justify-content: center; width: 100%">
                <Loader/>
            </div>
            {:else if invalidModCode}
            <p style="color: red">Invalid modification code</p>
            {:else if newReview.modCode}
            <p>Your modification code is: <span style="color: green">{newReview.modCode}</span></p>
            {:else}
            <PrimaryButton>
                {review == null ? "Create a review" : "Modify review"}
            </PrimaryButton>
            {/if}
        </form>
    </div>
</div>


<style>
    h3 {
        padding: 0;
        margin-bottom: 0.5em;
    }

    .form-container {
        display: flex;
        justify-content: center;
        align-items: center;
        position: absolute;
        left: 0;
        top: 0;
        z-index: 10000;
        width: 100vw;
        height: 100vh;
        background-color: rgb(0 0 0 / 50%);
    }

    .review-form {
        display: flex;
        flex-direction: column;
        justify-content: flex-start;
        background-color: #ffffff;
        width: 50%;
        height: 50%;
        border-radius: 1.25em;
        padding: 1em;
        border: 1px solid #aaaaaa;
        min-width: 600px;
        min-height: 600px;
    }

    form {
        display: flex;
        width: 100%;
        flex: 1;
        flex-direction: column;
        justify-content: space-between;
    }

    .form-contents {
        overflow: auto;
    }

    .form-contents > label {
        display: block;
        padding: 0.5em 0;
    }

    .form-contents > textarea {
        resize: none;
        width: calc(100%);
        height: 200px;
        box-sizing: border-box;
    }
</style>