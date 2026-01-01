<script lang="ts">
    import type { PharmacyTierRating } from "$lib/service/pharmacy-rating";
    import { tierRatingData } from "$lib/stores";
    import { pharmacyViewZIndex } from "$lib/utils/z-indices";
    import { onDestroy } from "svelte";
    import TitleBar from "../common/TitleBar.svelte";
    import Sidepanel from "./Sidepanel.svelte";
    import Loader from "../common/widgets/Loader.svelte";

    let {
        onSelectPharmacy,
        onClose
    }: {
        onSelectPharmacy: (id: number) => void;
        onClose: () => void
    } = $props();

    let pharmacies: PharmacyTierRating[] | undefined = $state(undefined);

    const unsubscribeTierRatingData = tierRatingData.subscribe((initial) => {
        if (initial != null) {
            pharmacies = initial;
        }
    });

    onDestroy(() => {
        unsubscribeTierRatingData();
    });
</script>

<Sidepanel zIndex={pharmacyViewZIndex} rightAlign={true}>
    <TitleBar onClose={onClose}/>
    {#if pharmacies == null}
        <div class="loader-container">
            <Loader/>
        </div>
    {:else}
        <div class="pharmacy-rating-container">
        {#each pharmacies as pharmacy}
            {#if pharmacy.avgRating != null}
                <button
                    class="pharmacy-rating"
                    title="{pharmacy.name}" style="--red: {pharmacy.avgRating != 0 ? 255*(1-(pharmacy.avgRating-1)/4) : 170}; --green: {pharmacy.avgRating != 0 ? 255*(pharmacy.avgRating-1)/4 : 170}; --blue: {pharmacy.avgRating == 0 ? 170 : 0}"
                    onclick={(_) => onSelectPharmacy(pharmacy?.id ?? 0)}
                >
                    <span>{pharmacy.avgRating?.toPrecision(2)}/5</span>
                    <span class="name">{pharmacy.name}</span>
                </button>
            {/if}
        {/each}
        </div>
    {/if}
</Sidepanel>

<style>
    .loader-container {
        display: flex;
        justify-content: center;
        width: 100%;
        margin-top: 1em;
        height: fit-content;
    }

    .pharmacy-rating-container {
        flex: 1;
        width: 100%;
        overflow: auto;
    }

    .pharmacy-rating {
        all: unset;
        display: flex;
        justify-content: flex-start;
        align-items: center;
        height: 32px;
        width: 100%;
        background-color: rgba(var(--red), var(--green), var(--blue), 0.75);
        user-select: none;
        box-sizing: border-box;
        transition: 0.1s all ease-in-out;

        &:hover {
            cursor: pointer;
            background-color: rgba(var(--red), var(--green), 0, 0.5);
        }

        & > .name {
            flex: 1;
            margin-left: 1em;
            text-wrap: nowrap;
            overflow: hidden;
        }
    }
</style>