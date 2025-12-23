<script lang="ts">
    import { PharmacyReview } from "$lib/service/pharmacy-review";
    import TitleBar from "../../common/TitleBar.svelte";
    import Countries from "$lib/assets/countries.json";
    import type { PharmacyInfo } from "$lib/service/pharmacy-info";
    import Loader from "../../common/widgets/Loader.svelte";

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
                <div class="star-container">
                    <input type="radio" name="stars" class="star" value="1" required>
                    <input type="radio" name="stars" class="star" value="2" required>
                    <input type="radio" name="stars" class="star" value="3" required>
                    <input type="radio" name="stars" class="star" value="4" required>
                    <input type="radio" name="stars" class="star" value="5" required checked>
                </div>
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
                            <option value={key} selected>{`${Countries[key].emoji} ${Countries[key].name}`}</option>
                            {:else}
                            <option value={key}>{`${Countries[key].emoji} ${Countries[key].name}`}</option>
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
            <p>Your modification code is: {newReview.modCode}</p>
            {:else}
            <button class="primary-btn">
                {review == null ? "Create a review" : "Modify review"}
            </button>
            {/if}
        </form>
    </div>
</div>


<style>
    .primary-btn {
        border: none;
        border-radius: 0.5em;
        min-width: 200px;
        box-shadow: none;
        background-color: #30bfff;
        color: white;
        box-sizing: border-box;
        padding: 1em;
        transition: all 0.5s ease-in-out;
    }

    .primary-btn:hover {
        cursor: pointer;
        background-color: #10a7ff;
    }

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

    .star-container {
        display: inline-block;
        user-select: none;
        padding: 0;
        margin: 0;
        outline: none;
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

    .star {
        margin: 0;
        appearance: none;
        width: 24px;
        height: 24px;
        background: url("data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjQiIGhlaWdodD0iMjQiIHZlcnNpb249IjEuMSIgdmlld0JveD0iMCAwIDguNDY2NjY2NiA4LjQ2NjY2NjYiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CiAgICA8ZyBmaWxsPSIjZmZiMDAwIiBzdHJva2Utd2lkdGg9Ii4yMzk3NDYiPgogICAgICAgIDxwYXRoIGQ9Im00LjIzMzI3MDMgMC40MjMzMzMwMWMtMC4yNTgwNjc4IDAtMC41MTYwODU1IDAuMTQ5MDUyNi0wLjYxMDE1MyAwLjQ0NjgyNjZsLTAuNTAzMzA2MSAxLjU5MjY3MDZjLTAuMDczNDkxMyAwLjIzMjYyNTctMC4yODk0MzQ3IDAuMzkxMTQ5LTAuNTMzNzY2MyAwLjM5MTE0OWgtMS41MjE2MzQxYy0wLjYxOTcwNDA0IDAtMC44Nzc1MjI1NCAwLjc5MTMzNTgtMC4zNzYzMDgwNCAxLjE1NTIwMDRsMS4zNTUyNzE1IDAuOTgzOTU1M2MwLjIwNDk4MzQgMC4xNDg4MTMyIDAuMjg0MDU3OSAwLjQxNjUyMDQgMC4xOTMwNzM3IDAuNjUyNjk1N2wtMC41ODc2NTg1IDEuNTI1NzYyOGMtMC4yMjg4MjU0IDAuNTk0MDAwOSAwLjQ1NzgzMjEgMS4xMjA2MDI2IDAuOTczMzM5OSAwLjc0NjI3MWwxLjI4MTY5NTktMC45MzEwODUxYzAuMTk2MzE0OS0wLjE0MjU1NTUgMC40NjI1NzU5LTAuMTQyNTU1NSAwLjY1ODg5MDggMGwxLjI4MTY5NTkgMC45MzEwODUxYzAuNTE1NTA2OSAwLjM3NDMzMTYgMS4yMDIxNzM1LTAuMTUyMjcwMSAwLjk3MzMzOS0wLjc0NjI3MWwtMC41ODc2NTg2LTEuNTI1NzYyOGMtMC4wOTEwMTE0LTAuMjM2MTc1My0wLjAxMTQ0NDQtMC41MDM4ODI1IDAuMTkzNTQzNS0wLjY1MjY5NTdsMS4zNTQ4MDE3LTAuOTgzOTU1M2MwLjUwMTIyNDktMC4zNjM4NjQ2IDAuMjQzMzk2Ni0xLjE1NTIwMDQtMC4zNzYzMDcyLTEuMTU1MjAwNGgtMS41MjE2MzRjLTAuMjQ0MzA1MiAwLTAuNDYwMzEyMy0wLjE1ODUyMzMtMC41MzM3NjcxLTAuMzkxMTQ5bC0wLjUwMzMwNjEtMS41OTI2NzA2Yy0wLjA5NDA4MTItMC4yOTc3NzQtMC4zNTIwODUzLTAuNDQ2ODI2Ni0wLjYxMDE1MjktMC40NDY4MjY2eiIvPgogICAgPC9nPgo8L3N2Zz4=");
    }

    .star:hover {
        cursor: pointer;
    }

    .star:checked ~ .star {
        background: url("data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjQiIGhlaWdodD0iMjQiIHZlcnNpb249IjEuMSIgdmlld0JveD0iMCAwIDguNDY2NjY2NiA4LjQ2NjY2NjYiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CiAgICA8ZyBmaWxsPSIjY2ZjZmNmIiBzdHJva2Utd2lkdGg9Ii4yMzk3NDYiPgogICAgICAgIDxwYXRoIGQ9Im00LjIzMzI3MDMgMC40MjMzMzMwMWMtMC4yNTgwNjc4IDAtMC41MTYwODU1IDAuMTQ5MDUyNi0wLjYxMDE1MyAwLjQ0NjgyNjZsLTAuNTAzMzA2MSAxLjU5MjY3MDZjLTAuMDczNDkxMyAwLjIzMjYyNTctMC4yODk0MzQ3IDAuMzkxMTQ5LTAuNTMzNzY2MyAwLjM5MTE0OWgtMS41MjE2MzQxYy0wLjYxOTcwNDA0IDAtMC44Nzc1MjI1NCAwLjc5MTMzNTgtMC4zNzYzMDgwNCAxLjE1NTIwMDRsMS4zNTUyNzE1IDAuOTgzOTU1M2MwLjIwNDk4MzQgMC4xNDg4MTMyIDAuMjg0MDU3OSAwLjQxNjUyMDQgMC4xOTMwNzM3IDAuNjUyNjk1N2wtMC41ODc2NTg1IDEuNTI1NzYyOGMtMC4yMjg4MjU0IDAuNTk0MDAwOSAwLjQ1NzgzMjEgMS4xMjA2MDI2IDAuOTczMzM5OSAwLjc0NjI3MWwxLjI4MTY5NTktMC45MzEwODUxYzAuMTk2MzE0OS0wLjE0MjU1NTUgMC40NjI1NzU5LTAuMTQyNTU1NSAwLjY1ODg5MDggMGwxLjI4MTY5NTkgMC45MzEwODUxYzAuNTE1NTA2OSAwLjM3NDMzMTYgMS4yMDIxNzM1LTAuMTUyMjcwMSAwLjk3MzMzOS0wLjc0NjI3MWwtMC41ODc2NTg2LTEuNTI1NzYyOGMtMC4wOTEwMTE0LTAuMjM2MTc1My0wLjAxMTQ0NDQtMC41MDM4ODI1IDAuMTkzNTQzNS0wLjY1MjY5NTdsMS4zNTQ4MDE3LTAuOTgzOTU1M2MwLjUwMTIyNDktMC4zNjM4NjQ2IDAuMjQzMzk2Ni0xLjE1NTIwMDQtMC4zNzYzMDcyLTEuMTU1MjAwNGgtMS41MjE2MzRjLTAuMjQ0MzA1MiAwLTAuNDYwMzEyMy0wLjE1ODUyMzMtMC41MzM3NjcxLTAuMzkxMTQ5bC0wLjUwMzMwNjEtMS41OTI2NzA2Yy0wLjA5NDA4MTItMC4yOTc3NzQtMC4zNTIwODUzLTAuNDQ2ODI2Ni0wLjYxMDE1MjktMC40NDY4MjY2eiIvPgogICAgPC9nPgo8L3N2Zz4=");
    }
</style>