<script lang="ts">
    import type { Snippet } from "svelte";
    import AccountButton from "../buttons/AccountButton.svelte";

    let {
        size,
        title = undefined,
        color = "#999999",
        children = undefined
    }: {
        size: number,
        title?: string,
        color?: string,
        children?: Snippet
    } = $props();

    let isDialogVisible: boolean = $state(false);
    let dialogElement: HTMLDialogElement;
</script>

<span class="dropdown-container" style="--size: {size}px">
    <AccountButton
        size={size}
        title={title || ""}
        padding={0}
        color={color}
        on:click={(_) => {
            if (isDialogVisible)
                dialogElement.close();
            else
                dialogElement.show();
            isDialogVisible = !isDialogVisible;
        }}
    />

    <dialog bind:this={dialogElement}>
        {@render children?.()}
    </dialog>
</span>

<style>
    .dropdown-container {
        display: inline-block;
        vertical-align: middle;
        padding: 0.5em;
        width: var(--size);
        height: var(--size);

        & > dialog {
            text-align: center;
            background-color: white;
            position: relative;
            border: 1px solid black;
            margin: 0;
            left: calc(var(--size) - 180px);
            width: 180px;
            border-radius: 1em;
            padding: 0;
            overflow: hidden;
            height: fit-content;
            box-sizing: border-box;

            :global(& > button), :global(& > a) {
                all: unset;
                display: flex;
                flex-direction: row;
                align-items: center;
                user-select: none;
                padding: 0;
                height: 48px;
                width: 100%;
                transition: all 0.3s ease-in-out;

                & > p {
                    flex: 1;
                }

                &:hover {
                    cursor: pointer;
                    background-color: rgba(0, 0, 0, 0.15)
                }
            }
        }
    }
</style>