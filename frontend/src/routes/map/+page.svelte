<script lang="ts">
    import type { PharmacyInfo } from "$lib/service/pharmacy-info";
    import type { PageProps } from "../$types";
    import LeafletMap from "../../components/LeafletMap.svelte";

    let { data }: PageProps = $props();

    async function showPharmacyView(pharmacy: PharmacyInfo) {
        console.log(pharmacy)
    }
</script>

<main>
    <LeafletMap pharmacies={data.pharmacies} callback={showPharmacyView}/>
</main>