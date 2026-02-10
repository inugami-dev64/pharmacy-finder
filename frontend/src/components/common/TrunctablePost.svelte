<script lang="ts">
    import { _ } from "svelte-i18n";
    import { onMount, type Snippet } from "svelte";

    let {
        postText,
        children = undefined
    }: {
        postText: string,
        children?: Snippet
    } = $props();

    let truncateText: boolean = $state(true);
    let textElement: HTMLElement;
    let truncatable: boolean = $state(false);

    onMount(() =>
        truncatable = textElement.offsetHeight < textElement.scrollHeight || textElement.offsetWidth < textElement.scrollWidth
    );
</script>


<div class="truncatable-post">
    {@render children?.()}
    <div id="content">
        <p class="{truncateText ? "hide-content" : ""}" bind:this={textElement}>{postText}</p>
    </div>

    {#if truncatable}
        <button
            onclick={() => truncateText = !truncateText}
            aria-label="Toggle Read More"
            id="toggleBtn"
        >
            {truncateText ? $_("map.sidebar.review.showMore") : $_("map.sidebar.review.showLess")}
        </button>
    {/if}
</div>

<style>
    .truncatable-post {
        margin-right: 15px;
        padding-bottom: 0.25em;
        border-bottom: 1px solid #cfcfcf;
        text-wrap: wrap;

        & > #content {
            white-space: pre;
            text-wrap: wrap;

            & > .hide-content {
                overflow: hidden;
                -webkit-line-clamp: 5;
                line-clamp: 5;
                -webkit-box-orient: vertical;
                box-orient: vertical;
                display: -webkit-box;
            }
        }
    }
</style>