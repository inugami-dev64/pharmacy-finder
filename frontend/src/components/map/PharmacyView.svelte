<script lang="ts">
    import type { PharmacyInfo } from "$lib/service/pharmacy-info";
    import CloseButton from "../common/icons/CloseButton.svelte";
    
    // Import logos
    import BenuLogo from "$lib/assets/benu-logo.svg"
    import ApothekaLogo from "$lib/assets/apotheka-logo.svg"
    import SudameapteekLogo from "$lib/assets/sudameapteek-logo.svg"
    import EuroapteekLogo from "$lib/assets/euroapteek-logo.svg"

    export let pharmacy: PharmacyInfo;
    export let visible: boolean;
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
            <p>{pharmacy.aadress}, {pharmacy.postalCode}, {pharmacy.city}, {pharmacy.county}</p>
        </div>
    </div>
</div>
{/if}

<style>
    h3, p {
        padding: 0;
        margin: 0;
    }

    .phr-view {
        position: fixed;
        left: 1em;
        top: 1em;
        width: calc(25% - 2em);
        min-width: 360px;
        max-width: 400px;
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

    .close {
        display: flex;
        flex-direction: row;
        width: 100%;
        height: fit-content;
        justify-content: right;
        margin-bottom: 2em;
    }
</style>