<script lang="ts">
    import type { PharmacyInfo } from "$lib/service/pharmacy-info";
    import CloseButton from "../common/icons/CloseButton.svelte";
    import StarRating from "../common/widgets/StarRating.svelte";

    // Import logos
    import BenuLogo from "$lib/assets/benu-logo.svg"
    import ApothekaLogo from "$lib/assets/apotheka-logo.svg"
    import SudameapteekLogo from "$lib/assets/sudameapteek-logo.svg"
    import EuroapteekLogo from "$lib/assets/euroapteek-logo.svg"
    import Loader from "../common/widgets/Loader.svelte";

    export let pharmacy: PharmacyInfo;
    export let visible: boolean;

    let showMoreAverageScores: boolean = false;

</script>

{#if visible}
<div class="phr-view">
    <div class="close">
        <CloseButton size=32 on:click={(e) => visible = !visible}/>
    </div>

    <div>
        {#if pharmacy.chain?.toLowerCase() == "benu"}
            <img alt="Benu logo" src="{BenuLogo}">
        {:else if pharmacy.chain?.toLowerCase() == "apotheka" }
            <img alt="Apotheka logo" src={ApothekaLogo}>
        {:else if pharmacy.chain?.toLowerCase() == "südameapteek"}
            <img alt="Südameapteek logo" src="{SudameapteekLogo}">
        {:else if pharmacy.chain?.toLowerCase() == "euroapteek"}
            <img alt="Euroapteek logo" src="{EuroapteekLogo}">
        {/if}

        <div class="phr-info">
            <p>{pharmacy.chain}</p>
            <h3>{pharmacy.name}</h3>
            <p>{pharmacy.address},</p>
            <p>{pharmacy.city},</p>
            <p>{pharmacy.county},</p>
            <p>{pharmacy.postalCode}</p>
            <!--<span><StarRating value={pharmacy.overallRating || 0} title="Overall rating"/></span>-->
            {#if !showMoreAverageScores}
            <button on:click|preventDefault={_ => showMoreAverageScores = !showMoreAverageScores}>View more...</button>
            {:else}
            <!--<span><StarRating value={pharmacy.acceptanceRating || 0} title="Acceptance rating"/></span>
            <span><StarRating value={pharmacy.eRating || 0} title="E rating"/></span>
            <span><StarRating value={pharmacy.tRating || 0} title="T rating"/></span>-->
            <button on:click|preventDefault={_ => showMoreAverageScores = !showMoreAverageScores}>View less</button>
            {/if}
        </div>

        <div class="phr-ratings">
            <div class="loader-container">
                <Loader/>
            </div>
        </div>
    </div>
</div>
{/if}

<style>
    h3, p {
        padding: 0;
        margin: 0;
    }

    .loader-container {
        display: flex;
        justify-content: center;
        width: 100%;
        height: fit-content;
    }

    .phr-view {
        position: fixed;
        left: 1em;
        top: 1em;
        width: calc(25% - 2em);
        min-width: 380px;
        max-width: 420px;
        height: fit-content;
        max-height: calc(100% - 2em);
        padding: 1em;
        z-index: 1000;
        background-color: #ffffff;
        border: 1px solid #aaaaaa;
        border-radius: 1.25em;
        box-sizing: border-box;
    }

    .phr-view > div > img {
        margin-bottom: 2em;
        width: 100%;
    }

    .phr-info {
        margin-bottom: 0.5em;
        border-bottom: 1px solid #cfcfcf;
    }

    .close {
        display: flex;
        flex-direction: row;
        width: 100%;
        height: fit-content;
        justify-content: right;
        margin-bottom: 2em;
    }
</style>