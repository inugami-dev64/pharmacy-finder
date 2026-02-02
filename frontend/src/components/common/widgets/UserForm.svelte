<script lang="ts">
    import type { Snippet } from "svelte";
    import PrimaryButton from "./buttons/PrimaryButton.svelte";
    import Loader from "./Loader.svelte";

    let {
        title,
        submitBtnText,
        isSubmitted = false,
        validationErrors = [],
        submitForm = undefined,
        children = undefined
    }: {
        title: string,
        submitBtnText: string,
        isSubmitted?: boolean,
        validationErrors?: Array<string>,
        submitForm?: (e: SubmitEvent) => void,
        children?: Snippet
    } = $props();
</script>

<form
    onsubmit={
        (e: SubmitEvent) => {
            submitForm?.(e);
        }
    }
>
    <h3>{title}</h3>
    {#if validationErrors.length > 0}
        <ul>
        {#each validationErrors as err}
            <li class="err">{err}</li>
        {/each}
        </ul>
    {/if}
    {@render children?.()}

    {#if !isSubmitted}
        <PrimaryButton>{submitBtnText}</PrimaryButton>
    {:else}
        <div class="loader-container">
            <Loader/>
        </div>
    {/if}
</form>

<style>
    .err {
        color: red;
    }

    form {
        text-align: center;
        width: 400px;
        display: flex;
        flex-direction: column;
        padding: 10px;
        border: 1px solid black;
        border-radius: 1em;

        :global(& > input) {
            all: unset;
            padding: 0.5em;
            margin: 5px 0;
            border-radius: 0.5em;
            background-color: white;
            border: 1px solid #30bfff;
            transition: all 0.5s ease-in-out;

            :global(&:focus) {
                border: 1px solid #0000ff;
            }
        }
    }
</style>
