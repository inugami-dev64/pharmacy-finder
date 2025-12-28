<script lang="ts">
    import type { PharmacyInfo } from "$lib/service/pharmacy-info";
    import { PharmacyRating } from "$lib/service/pharmacy-rating";
    import { PharmacyReview } from "$lib/service/pharmacy-review";
    import { ratingData, reviewData } from "$lib/stores";
    import type { PageProps } from "../$types";
    import LanguageSwitcher from "../../components/common/LanguageSwitcher.svelte";
    import LeafletMap from "../../components/map/LeafletMap.svelte";
    import PharmacyView from "../../components/map/PharmacyView.svelte";

    let { data }: PageProps = $props();

    let activePharmacy: PharmacyInfo | undefined = $state();
    let visible: boolean = $state(false);

    async function showPharmacyView(pharmacy: PharmacyInfo) {
        // empty the reviewData store
        reviewData.set(undefined);
        ratingData.set(undefined)

        activePharmacy = pharmacy;
        visible = true;

        if (pharmacy.id != null) {
            reviewData.set(await PharmacyReview.readReviews(pharmacy.id, undefined, undefined));
            ratingData.set(await PharmacyRating.readPharmacyRatings(pharmacy.id));
        }
    }
</script>

<svelte:head>
    <title>Pharmacy finder | Map</title>
</svelte:head>

<main>
    <div class="language-switcher-container">
        <LanguageSwitcher/>
    </div>
    {#if activePharmacy != null && visible}
    <PharmacyView pharmacy={<PharmacyInfo>activePharmacy} onClose={() => visible = false}/>
    {/if}
    <LeafletMap pharmacies={(<{pharmacies: PharmacyInfo[]}>data).pharmacies} callback={showPharmacyView}/>
</main>

<style>
    main {
        width: 100%;
        height: 100%;
    }

    .language-switcher-container {
        display: block;
        position: absolute;
        right: 0;
        top: 0;
        z-index: 10000;
    }
</style>