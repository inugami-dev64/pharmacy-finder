<script lang="ts">
    import type { PharmacyInfo } from "$lib/service/pharmacy-info";
    import type { PageProps } from "../$types";
    import LeafletMap from "../../components/map/LeafletMap.svelte";
    import PharmacyView from "../../components/map/PharmacyView.svelte";

    let { data }: PageProps = $props();

    let activePharmacy: PharmacyInfo | undefined = $state();
    let visible: boolean = $state(false);

    async function showPharmacyView(pharmacy: PharmacyInfo) {
        activePharmacy = pharmacy;
        visible = true;
    }
</script>

<main>
    <PharmacyView pharmacy={<PharmacyInfo>activePharmacy} bind:visible={visible}/>
    <LeafletMap pharmacies={data.pharmacies} callback={showPharmacyView}/>
</main>

<style>
</style>