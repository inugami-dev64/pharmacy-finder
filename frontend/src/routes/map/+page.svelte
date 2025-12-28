<script lang="ts">
    import type { PharmacyInfo } from "$lib/service/pharmacy-info";
    import { PharmacyRating } from "$lib/service/pharmacy-rating";
    import { PharmacyReview } from "$lib/service/pharmacy-review";
    import { ratingData, reviewData } from "$lib/stores";
    import { navBarZIndex } from "$lib/utils/z-indices";
    import { locale } from "svelte-i18n";
    import type { PageProps } from "../$types";
    import LanguageButton from "../../components/common/icons/buttons/LanguageButton.svelte";
    import SearchButton from "../../components/common/icons/buttons/SearchButton.svelte";
    import ShinyStarButton from "../../components/common/icons/buttons/ShinyStarButton.svelte";
    import SourceCodeButton from "../../components/common/icons/buttons/SourceCodeButton.svelte";
    import NavBar from "../../components/common/widgets/navbar/NavBar.svelte";
    import LeafletMap from "../../components/map/LeafletMap.svelte";
    import SearchModal from "../../components/map/NavBar/SearchModal.svelte";
    import PharmacyView from "../../components/map/PharmacyView.svelte";
    import { languages } from "$lib/utils/languages";
    import { _ } from "svelte-i18n";

    let { data }: PageProps = $props();

    let activePharmacy: PharmacyInfo | undefined = $state();
    let pharmacyViewVisible: boolean = $state(false);
    let searchVisible: boolean = $state(false);

    async function showPharmacyView(pharmacy: PharmacyInfo) {
        // empty the reviewData store
        reviewData.set(undefined);
        ratingData.set(undefined)

        activePharmacy = pharmacy;
        pharmacyViewVisible = true;

        if (pharmacy.id != null) {
            reviewData.set(await PharmacyReview.readReviews(pharmacy.id, undefined, undefined));
            ratingData.set(await PharmacyRating.readPharmacyRatings(pharmacy.id));
        }
    }

    let langSelector: HTMLSelectElement;
</script>

<svelte:head>
    <title>Pharmacy finder | Map</title>
</svelte:head>

<main>
    <div class="navbar-container" style="--zIndex: {navBarZIndex}">
        <NavBar size={48}>
            <SearchButton size={32} on:click={() => searchVisible = true} title={$_("map.navbar.search")}/>
            <!--<ShinyStarButton size={32}/>-->
            <div>
                <LanguageButton size={32} on:click={_ => langSelector.showPicker()} title={$_("map.navbar.language")}/>
                <select
                    bind:this={langSelector}
                    onchange={e => {
                        e.preventDefault();
                        const value = (e.target as HTMLSelectElement).value;
                        locale.set(value);
                    }}
                >
                    {#each languages as selection}
                    <option value={selection.code}>{selection.language}</option>
                    {/each}
                </select>
            </div>
            <SourceCodeButton size={32} title={$_("map.navbar.source")}/>
        </NavBar>
    </div>
    {#if activePharmacy != null && pharmacyViewVisible}
    <PharmacyView pharmacy={<PharmacyInfo>activePharmacy} onClose={() => pharmacyViewVisible = false}/>
    {/if}
    {#if searchVisible}
    <SearchModal
        pharmacies={(<{pharmacies: PharmacyInfo[]}>data).pharmacies}
        onSelect={(v) => {
            activePharmacy = v;
            searchVisible = false;
        }}
        onClose={() => searchVisible = false}
        />
    {/if}
    <LeafletMap selectedPharmacy={activePharmacy} pharmacies={(<{pharmacies: PharmacyInfo[]}>data).pharmacies} callback={showPharmacyView}/>
</main>

<style>
    select {
        all: unset;
        display: block;
        width: 0;
        height: 0;
        padding: 0;
        margin: 0;
    }

    select option {
        font-size: 14px;
        color: black;
    }

    main {
        width: 100%;
        height: 100%;
    }

    .navbar-container {
        display: flex;
        flex-direction: row;
        position: absolute;
        right: 10px;
        top: 10px;
        z-index: var(--zIndex);
    }
</style>