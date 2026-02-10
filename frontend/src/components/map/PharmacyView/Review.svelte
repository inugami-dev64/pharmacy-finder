<script lang="ts">
    import Stars from "../../common/widgets/stars/Stars.svelte";
    import Estrogen from "$lib/assets/regimen/estrogen.webp";
    import Testosterone from "$lib/assets/regimen/testosterone.webp";
    import Countries from "$lib/assets/countries.json"
    import type { PharmacyReview } from "$lib/service/data/pharmacy-review";
    import ImagoLogo from "../../common/icons/logos/ImagoLogo.svelte";
    import GenderGPLogo from "../../common/icons/logos/GenderGPLogo.svelte";
    import { _ } from "svelte-i18n";
    import { ModeratorPharmacyReview } from "$lib/service/data/moderator-pharmacy-review";
    import TrunctablePost from "../../common/TrunctablePost.svelte";
    import type { Snippet } from "svelte";

    let {
        review,
        children = undefined
    }: {
        review: ModeratorPharmacyReview | PharmacyReview,
        children?: Snippet
    } = $props();
</script>

<TrunctablePost postText={review.review || ""}>
    <div class="review-header">
        <span>
            <Stars value={review.stars ?? 5} scale={0.75}/>
            {@render children?.()}
        </span>
        <span style="display: flex; align-items: center">
            {#if review.nationality != null && review.nationality in Countries}
            <span style="font-size: 20px" title={Countries[review.nationality as keyof typeof Countries].name}>{Countries[review.nationality as keyof typeof Countries].emoji}</span>
            {/if}
            {#if review.hrtKind === 'e'}
            <img src="{Estrogen}" alt="e" title={$_("map.sidebar.review.eBased")}>
            {:else if review.hrtKind === 't'}
            <img src="{Testosterone}" alt="t" title={$_("map.sidebar.review.tBased")}>
            {/if}
        </span>
    </div>
    <time>{new Date(review.updatedAt ?? 0).toLocaleDateString()} {new Date(review.updatedAt ?? 0).toLocaleTimeString()}</time>
    <!-- Prescription type -->
    <span style="text-align: center">
        {#if review.prescriptionType === "Imago"}
            <ImagoLogo size={24}/>
        {:else if review.prescriptionType === "GenderGP"}
            <GenderGPLogo size={22}/>
        {:else if review.prescriptionType === "National"}
            🇪🇪
        {/if}
        <b>{review.prescriptionType}</b>
    </span>
</TrunctablePost>

<style>
    time {
        display: block;
        margin-bottom: 0.25em;
    }

    .review-header  {
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .review-header {
        user-select: none;
    }

    img {
        width: 24px;
        height: 24px;
        user-select: none;
    }
</style>